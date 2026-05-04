import testing from '@taskcluster/lib-testing';
import SchemaSet from '@taskcluster/lib-validate';
import { MonitorManager } from '@taskcluster/lib-monitor';
import assert from 'assert';
import net from 'net';
import path from 'path';
import { App } from '@taskcluster/lib-app';

const __dirname = new URL('.', import.meta.url).pathname;
let runningServer = null;

const pickFreePort = () => new Promise((resolve, reject) => {
  const srv = net.createServer();
  srv.unref();
  srv.on('error', reject);
  srv.listen(0, () => {
    const { port } = srv.address();
    srv.close(() => resolve(port));
  });
});

const port = await pickFreePort();
export const rootUrl = `http://localhost:${port}`;

export let monitor = null;
export let monitorManager = null;

export const setupMonitor = () => {
  monitor = MonitorManager.setup({
    serviceName: 'lib-api',
    fake: true,
    debug: true,
    verify: true,
    level: 'debug',
  });
  monitorManager = monitor.manager;
};

export const resetMonitorManager = () => {
  monitorManager.reset();
};

export const setupServer = async ({ builder, context }) => {
  testing.fakeauth.start({
    'client-with-aa-bb-dd': ['aa', 'bb', 'dd'],
  }, { rootUrl });
  assert(runningServer === null);

  const schemaset = new SchemaSet({
    serviceName: 'test',
    folder: path.join(__dirname, 'schemas'),
  });

  const api = await builder.build({
    rootUrl,
    schemaset,
    monitor,
    context,
  });

  runningServer = await App({
    port,
    env: 'development',
    forceSSL: false,
    trustProxy: false,
    apis: [api],
  });
};

export const teardownServer = async () => {
  if (runningServer) {
    await new Promise(function(accept) {
      runningServer.once('close', function() {
        runningServer = null;
        accept();
      });
      runningServer.close();
    });
  }
  testing.fakeauth.stop();
};

export default { rootUrl, setupServer, teardownServer };
