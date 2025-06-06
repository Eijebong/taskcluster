import { Octokit } from '@octokit/rest';
import { createAppAuth } from '@octokit/auth-app';
import { throttling } from '@octokit/plugin-throttling';
import { retry } from '@octokit/plugin-retry';
import Bottleneck from "bottleneck";

const PluggedOctokit = Octokit.plugin(retry, throttling);

const tokenCache = new Map();
/**
 * Github app access token have an expiration date of 1h
 * https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#authenticating-as-an-installation
 *
 * As we run multiple handlers for the same repository, we can save a lot of resources by using cached token
 * On average it takes 0.5-1s to get one.
 */
export const getCachedInstallationToken = async (gh, inst_id) => {
  let tokenData = tokenCache.get(inst_id);
  const timeMargin = 10 * 60 * 1000; // 10min before expiry
  if (tokenData) {
    if (new Date(tokenData.expires_at).getTime() > new Date().getTime() + timeMargin) {
      return tokenData;
    }
  }

  tokenData = (await gh.apps.createInstallationAccessToken({
    installation_id: inst_id,
  })).data;

  tokenCache.set(inst_id, tokenData);
  return tokenData;
};

export const getPrivatePEM = cfg => {
  const keyRe = /-----BEGIN RSA PRIVATE KEY-----(\n|\\n).*(\n|\\n)-----END RSA PRIVATE KEY-----(\n|\\n)?/s;
  const privatePEM = cfg.github.credentials.privatePEM;
  if (!keyRe.test(privatePEM)) {
    throw new Error(`Malformed GITHUB_PRIVATE_PEM: must match ${keyRe}; ` +
      `got a value of length ${privatePEM.length}`);
  }

  // sometimes it's easier to provide this config value with embedded backslash-n characters
  // than to convince everything to correctly handle newlines.  So, we'll be friendly to that
  // arrangement, too.
  return privatePEM.replace(/\\n/g, '\n');
};

export default async ({ cfg, monitor }) => {
  const privatePEM = getPrivatePEM(cfg);

  const OctokitOptions = {
    log: {
      debug: message => monitor.debug(message),
      info: message => monitor.info(message),
      warn: message => monitor.warning(message),
      error: message => monitor.err(message),
    },
    throttle: {
      write: new Bottleneck.Group({ minTime: 50 }),
      onRateLimit: (retryAfter, options, octokit, retryCount) => {
        octokit.log.warn(
          `Request quota exhausted for request ${options.method} ${options.url}`,
        );

        if (retryCount < 3) {
          octokit.log.info(`Retrying after ${retryAfter} seconds!`);
          return true;
        }
      },
      onSecondaryRateLimit: (retryAfter, options, octokit) => {
        octokit.log.warn(
          `SecondaryRateLimit detected for request ${options.method} ${options.url}`,
        );
      },
    },
    retry: {
      // 404 and 401 are both retried because they can occur spuriously, likely due to MySQL db replication
      // delays at GitHub.
      doNotRetry: [400, 403],
    },
  };

  const getAppGithub = async () => {
    return new PluggedOctokit({
      authStrategy: createAppAuth,
      auth: {
        appId: cfg.github.credentials.appId,
        privateKey: privatePEM,
      },
      ...OctokitOptions,
    });
  };

  const getInstallationGithub = async (inst_id) => {
    try {
      const inteGithub = await getAppGithub();
      // Authenticating as installation
      const instaToken = await getCachedInstallationToken(inteGithub, inst_id);
      const instaGithub = new PluggedOctokit({
        auth: `token ${instaToken.token}`,
        ...OctokitOptions,
      });
      return instaGithub;
    } catch (err) {
      err.installationId = inst_id;
      throw err;
    }
  };

  // This object insures that the authentication is delayed until we need it.
  // Also, the authentication happens not just once in the beginning, but for each request.
  return { getAppGithub, getInstallationGithub };
};
