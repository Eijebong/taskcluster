// This source code file is AUTO-GENERATED by github.com/taskcluster/jsonschema2go

package gwconfig

import (
	"encoding/json"
)

type (
	// This schema defines the structure of the generic-worker.config file.
	ConfigFileForGenericWorker struct {

		// Taskcluster access token used by generic worker
		// to talk to taskcluster queue.
		AccessToken string `json:"accessToken"`

		// The EC2 availability zone of the worker.
		AvailabilityZone string `json:"availabilityZone,omitempty"`

		// The directory where task caches should be stored on
		// the worker. The directory will be created if it does
		// not exist. This may be a relative path to the
		// current directory, or an absolute path.
		//
		// Default:    "caches"
		CachesDir string `json:"cachesDir,omitempty"`

		// Taskcluster certificate, when using temporary
		// credentials only.
		Certificate string `json:"certificate,omitempty"`

		// The number of seconds between consecutive calls
		// to the provisioner, to check if there has been a
		// new deployment of the current worker type. If a
		// new deployment is discovered, worker will shut
		// down. See deploymentId property.
		//
		// Default:    1800
		CheckForNewDeploymentEverySecs float64 `json:"checkForNewDeploymentEverySecs,omitempty"`

		// Whether to delete the home directories of the task
		// users after the task completes. Normally you would
		// want to do this to avoid filling up disk space,
		// but for one-off troubleshooting, it can be useful
		// to (temporarily) leave home directories in place.
		// Accepted values: true or false.
		//
		// Default:    true
		CleanUpTaskDirs bool `json:"cleanUpTaskDirs,omitempty"`

		// Taskcluster client ID used by generic worker to
		// talk to taskcluster queue.
		ClientID string `json:"clientId"`

		// If running with --configure-for-aws, then between
		// tasks, at a chosen maximum frequency (see
		// checkForNewDeploymentEverySecs property), the
		// worker will query the provisioner to get the
		// updated worker type definition. If the deploymentId
		// in the config of the worker type definition is
		// different to the worker's current deploymentId, the
		// worker will shut itself down. See
		// https://bugzil.la/1298010
		DeploymentID string `json:"deploymentId,omitempty"`

		// If true, no system reboot will be initiated by
		// generic-worker program, but it will still return
		// with exit code 67 if the system needs rebooting.
		// This allows custom logic to be executed before
		// rebooting, by patching run-generic-worker.bat
		// script to check for exit code 67, perform steps
		// (such as formatting a hard drive) and then
		// rebooting in the run-generic-worker.bat script.
		//
		// Default:    false
		DisableReboots bool `json:"disableReboots,omitempty"`

		// The directory to cache downloaded files for
		// populating preloaded caches and readonly mounts. The
		// directory will be created if it does not exist. This
		// may be a relative path to the current directory, or
		// an absolute path.
		//
		// Default:    "downloads"
		DownloadsDir string `json:"downloadsDir,omitempty"`

		// The ed25519 signing key for signing artifacts with.
		Ed25519SigningKeyLocation string `json:"ed25519SigningKeyLocation"`

		// How many seconds to wait without getting a new
		// task to perform, before the worker process exits.
		// An integer, >= 0. A value of 0 means "never reach
		// the idle state" - i.e. continue running
		// indefinitely. See also shutdownMachineOnIdle.
		//
		// Default:    0
		IdleTimeoutSecs float64 `json:"idleTimeoutSecs,omitempty"`

		// The EC2 instance ID of the worker. Used by chain of trust.
		InstanceID string `json:"instanceId,omitempty"`

		// The EC2 instance Type of the worker. Used by chain of trust.
		InstanceType string `json:"instanceType,omitempty"`

		// Filepath of LiveLog executable to use; see
		// https://github.com/taskcluster/livelog
		//
		// Default:    "livelog"
		LivelogExecutable string `json:"livelogExecutable,omitempty"`

		// If zero, run tasks indefinitely. Otherwise, after
		// this many tasks, exit.
		//
		// Default:    0
		NumberOfTasksToRun int64 `json:"numberOfTasksToRun,omitempty"`

		// The private IP of the worker, used by chain of trust.
		PrivateIP string `json:"privateIP,omitempty"`

		// The taskcluster provisioner which is taking care
		// of provisioning environments with generic-worker
		// running on them.
		//
		// Default:    "test-provisioner"
		ProvisionerID string `json:"provisionerId,omitempty"`

		// The IP address for VNC access.  Also used by chain of
		// trust when present.
		PublicIP string `json:"publicIP,omitempty"`

		// The EC2 region of the worker. Used by chain of trust.
		Region string `json:"region,omitempty"`

		// The garbage collector will ensure at least this
		// number of megabytes of disk space are available
		// when each task starts. If it cannot free enough
		// disk space, the worker will shut itself down.
		//
		// Default:    10240
		RequiredDiskSpaceMegabytes float64 `json:"requiredDiskSpaceMegabytes,omitempty"`

		// The root URL of the taskcluster deployment to which
		// clientId and accessToken grant access. For example,
		// 'https://community-tc.services.mozilla.com/'.
		RootURL string `json:"rootURL"`

		// A string, that if non-empty, will be treated as a
		// command to be executed as the newly generated task
		// user, after the user has been created, the machine
		// has rebooted and the user has logged in, but before
		// a task is run as that user. This is a way to
		// provide generic user initialisation logic that
		// should apply to all generated users (and thus all
		// tasks) and be run as the task user itself. This
		// option does *not* support running a command as
		// Administrator. Furthermore, even if
		// runTasksAsCurrentUser is true, the script will still
		// be executed as the task user, rather than the
		// current user (that runs the generic-worker process).
		RunAfterUserCreation string `json:"runAfterUserCreation,omitempty"`

		// The project name used in https://sentry.io for
		// reporting worker crashes. Permission to publish
		// crash reports is granted via the scope
		// auth:sentry:<sentryProject>. If the taskcluster
		// client (see clientId property above) does not
		// posses this scope, no crash reports will be sent.
		// Similarly, if this property is not specified or
		// is the empty string, no reports will be sent.
		//
		// Default:    "generic-worker"
		SentryProject string `json:"sentryProject,omitempty"`

		// If true, when the worker is deemed to have been
		// idle for enough time (see idleTimeoutSecs) the
		// worker will issue an OS shutdown command. If false,
		// the worker process will simply terminate, but the
		// machine will not be shut down.
		//
		// Default:    false
		ShutdownMachineOnIdle bool `json:"shutdownMachineOnIdle,omitempty"`

		// If true, if the worker encounters an unrecoverable
		// error (such as not being able to write to a
		// required file) it will shutdown the host
		// computer. Note this is generally only desired
		// for machines running in production, such as on AWS
		// EC2 spot instances. Use with caution!
		//
		// Default:    false
		ShutdownMachineOnInternalError bool `json:"shutdownMachineOnInternalError,omitempty"`

		// Filepath of taskcluster-proxy executable to use; see
		// https://github.com/taskcluster/taskcluster/tree/main/tools/taskcluster-proxy
		//
		// Default:    "taskcluster-proxy"
		TaskclusterProxyExecutable string `json:"taskclusterProxyExecutable,omitempty"`

		// Port number for taskcluster-proxy HTTP requests.
		//
		// Default:    80
		TaskclusterProxyPort int64 `json:"taskclusterProxyPort,omitempty"`

		// The location where task directories should be
		// created on the worker.
		//
		// Default:    "${DEFAULT_TASKS_DIR}"
		TasksDir string `json:"tasksDir,omitempty"`

		// Typically this would be an aws region - an
		// identifier to uniquely identify which pool of
		// workers this worker logically belongs to.
		//
		// Default:    "test-worker-group"
		WorkerGroup string `json:"workerGroup,omitempty"`

		// A unique name for the worker.
		WorkerID string `json:"workerId"`

		// If a non-empty string, task commands will have environment variable
		// TASKCLUSTER_WORKER_LOCATION set to the value provided.
		// Otherwise TASKCLUSTER_WORKER_LOCATION environment
		// variable will not be implicitly set in task commands.
		//
		// Default:    ""
		WorkerLocation string `json:"workerLocation,omitempty"`

		// This should match a worker_type managed by the
		// specified provisioner.
		WorkerType string `json:"workerType"`

		// This arbitrary json blob will be included at the
		// top of each task log. Providing information here,
		// such as a URL to the code/config used to set up the
		// worker type will mean that people running tasks on
		// the worker type will have more information about how
		// it was set up (for example what has been installed on
		// the machine).
		//
		// Additional properties allowed
		WorkerTypeMetadata json.RawMessage `json:"workerTypeMetadata,omitempty"`

		// The audience value for which to request websocktunnel
		// credentials, identifying a set of WST servers this
		// worker could connect to.  Optional if not using websocktunnel
		// to expose live logs.
		WstAudience string `json:"wstAudience,omitempty"`

		// The URL of the websocktunnel server with which to expose
		// live logs.  Optional if not using websocktunnel to expose
		// live logs.
		WstServerURL string `json:"wstServerURL,omitempty"`
	}
)

// Returns json schema for the payload part of the task definition. Please
// note we use a go string and do not load an external file, since we want this
// to be *part of the compiled executable*. If this sat in another file that
// was loaded at runtime, it would not be burned into the build, which would be
// bad for the following two reasons:
//  1) we could no longer distribute a single binary file that didn't require
//     installation/extraction
//  2) the payload schema is specific to the version of the code, therefore
//     should be versioned directly with the code and *frozen on build*.
//
// Run `generic-worker show-payload-schema` to output this schema to standard
// out.
func taskPayloadSchema() string {
	return `{
  "$id": "/schemas/generic-worker/config.json#",
  "$schema": "/schemas/common/metaschema.json#",
  "additionalProperties": false,
  "description": "This schema defines the structure of the generic-worker.config file.",
  "properties": {
    "accessToken": {
      "description": "Taskcluster access token used by generic worker\nto talk to taskcluster queue.",
      "title": "Taskcluster Access Token",
      "type": "string"
    },
    "availabilityZone": {
      "description": "The EC2 availability zone of the worker.",
      "title": "Availability Zone (EC2)",
      "type": "string"
    },
    "cachesDir": {
      "default": "caches",
      "description": "The directory where task caches should be stored on\nthe worker. The directory will be created if it does\nnot exist. This may be a relative path to the\ncurrent directory, or an absolute path.",
      "title": "Caches Directory",
      "type": "string"
    },
    "certificate": {
      "description": "Taskcluster certificate, when using temporary\ncredentials only.",
      "title": "Taskcluster Certificate",
      "type": "string"
    },
    "checkForNewDeploymentEverySecs": {
      "default": 1800,
      "description": "The number of seconds between consecutive calls\nto the provisioner, to check if there has been a\nnew deployment of the current worker type. If a\nnew deployment is discovered, worker will shut\ndown. See deploymentId property.",
      "title": "Check For New Deployment (every X seconds)",
      "type": "number"
    },
    "cleanUpTaskDirs": {
      "default": true,
      "description": "Whether to delete the home directories of the task\nusers after the task completes. Normally you would\nwant to do this to avoid filling up disk space,\nbut for one-off troubleshooting, it can be useful\nto (temporarily) leave home directories in place.\nAccepted values: true or false.",
      "title": "Clean Up Task Directories",
      "type": "boolean"
    },
    "clientId": {
      "description": "Taskcluster client ID used by generic worker to\ntalk to taskcluster queue.",
      "title": "Taskcluster Client ID",
      "type": "string"
    },
    "deploymentId": {
      "description": "If running with --configure-for-aws, then between\ntasks, at a chosen maximum frequency (see\ncheckForNewDeploymentEverySecs property), the\nworker will query the provisioner to get the\nupdated worker type definition. If the deploymentId\nin the config of the worker type definition is\ndifferent to the worker's current deploymentId, the\nworker will shut itself down. See\nhttps://bugzil.la/1298010",
      "title": null,
      "type": "string"
    },
    "disableReboots": {
      "default": false,
      "description": "If true, no system reboot will be initiated by\ngeneric-worker program, but it will still return\nwith exit code 67 if the system needs rebooting.\nThis allows custom logic to be executed before\nrebooting, by patching run-generic-worker.bat\nscript to check for exit code 67, perform steps\n(such as formatting a hard drive) and then\nrebooting in the run-generic-worker.bat script.",
      "title": null,
      "type": "boolean"
    },
    "downloadsDir": {
      "default": "downloads",
      "description": "The directory to cache downloaded files for\npopulating preloaded caches and readonly mounts. The\ndirectory will be created if it does not exist. This\nmay be a relative path to the current directory, or\nan absolute path.",
      "title": null,
      "type": "string"
    },
    "ed25519SigningKeyLocation": {
      "description": "The ed25519 signing key for signing artifacts with.",
      "title": null,
      "type": "string"
    },
    "idleTimeoutSecs": {
      "default": 0,
      "description": "How many seconds to wait without getting a new\ntask to perform, before the worker process exits.\nAn integer, \u003e= 0. A value of 0 means \"never reach\nthe idle state\" - i.e. continue running\nindefinitely. See also shutdownMachineOnIdle.",
      "title": null,
      "type": "number"
    },
    "instanceId": {
      "description": "The EC2 instance ID of the worker. Used by chain of trust.",
      "title": null,
      "type": "string"
    },
    "instanceType": {
      "description": "The EC2 instance Type of the worker. Used by chain of trust.",
      "title": null,
      "type": "string"
    },
    "livelogExecutable": {
      "default": "livelog",
      "description": "Filepath of LiveLog executable to use; see\nhttps://github.com/taskcluster/livelog",
      "title": null,
      "type": "string"
    },
    "numberOfTasksToRun": {
      "default": 0,
      "description": "If zero, run tasks indefinitely. Otherwise, after\nthis many tasks, exit.",
      "title": null,
      "type": "integer"
    },
    "privateIP": {
      "description": "The private IP of the worker, used by chain of trust.",
      "format": "ipv4",
      "title": null,
      "type": "string"
    },
    "provisionerId": {
      "default": "test-provisioner",
      "description": "The taskcluster provisioner which is taking care\nof provisioning environments with generic-worker\nrunning on them.",
      "title": null,
      "type": "string"
    },
    "publicIP": {
      "description": "The IP address for VNC access.  Also used by chain of\ntrust when present.",
      "format": "ipv4",
      "title": null,
      "type": "string"
    },
    "region": {
      "description": "The EC2 region of the worker. Used by chain of trust.",
      "title": null,
      "type": "string"
    },
    "requiredDiskSpaceMegabytes": {
      "default": 10240,
      "description": "The garbage collector will ensure at least this\nnumber of megabytes of disk space are available\nwhen each task starts. If it cannot free enough\ndisk space, the worker will shut itself down.",
      "title": null,
      "type": "number"
    },
    "rootURL": {
      "description": "The root URL of the taskcluster deployment to which\nclientId and accessToken grant access. For example,\n'https://community-tc.services.mozilla.com/'.",
      "format": "uri",
      "title": null,
      "type": "string"
    },
    "runAfterUserCreation": {
      "description": "A string, that if non-empty, will be treated as a\ncommand to be executed as the newly generated task\nuser, after the user has been created, the machine\nhas rebooted and the user has logged in, but before\na task is run as that user. This is a way to\nprovide generic user initialisation logic that\nshould apply to all generated users (and thus all\ntasks) and be run as the task user itself. This\noption does *not* support running a command as\nAdministrator. Furthermore, even if\nrunTasksAsCurrentUser is true, the script will still\nbe executed as the task user, rather than the\ncurrent user (that runs the generic-worker process).",
      "title": null,
      "type": "string"
    },
    "sentryProject": {
      "default": "generic-worker",
      "description": "The project name used in https://sentry.io for\nreporting worker crashes. Permission to publish\ncrash reports is granted via the scope\nauth:sentry:\u003csentryProject\u003e. If the taskcluster\nclient (see clientId property above) does not\nposses this scope, no crash reports will be sent.\nSimilarly, if this property is not specified or\nis the empty string, no reports will be sent.",
      "title": null,
      "type": "string"
    },
    "shutdownMachineOnIdle": {
      "default": false,
      "description": "If true, when the worker is deemed to have been\nidle for enough time (see idleTimeoutSecs) the\nworker will issue an OS shutdown command. If false,\nthe worker process will simply terminate, but the\nmachine will not be shut down.",
      "title": null,
      "type": "boolean"
    },
    "shutdownMachineOnInternalError": {
      "default": false,
      "description": "If true, if the worker encounters an unrecoverable\nerror (such as not being able to write to a\nrequired file) it will shutdown the host\ncomputer. Note this is generally only desired\nfor machines running in production, such as on AWS\nEC2 spot instances. Use with caution!",
      "title": null,
      "type": "boolean"
    },
    "taskclusterProxyExecutable": {
      "default": "taskcluster-proxy",
      "description": "Filepath of taskcluster-proxy executable to use; see\nhttps://github.com/taskcluster/taskcluster/tree/main/tools/taskcluster-proxy",
      "title": null,
      "type": "string"
    },
    "taskclusterProxyPort": {
      "default": 80,
      "description": "Port number for taskcluster-proxy HTTP requests.",
      "title": null,
      "type": "integer"
    },
    "tasksDir": {
      "default": "${DEFAULT_TASKS_DIR}",
      "description": "The location where task directories should be\ncreated on the worker.",
      "title": null,
      "type": "string"
    },
    "workerGroup": {
      "default": "test-worker-group",
      "description": "Typically this would be an aws region - an\nidentifier to uniquely identify which pool of\nworkers this worker logically belongs to.",
      "title": null,
      "type": "string"
    },
    "workerId": {
      "description": "A unique name for the worker.",
      "title": null,
      "type": "string"
    },
    "workerLocation": {
      "default": "",
      "description": "If a non-empty string, task commands will have environment variable\nTASKCLUSTER_WORKER_LOCATION set to the value provided.\nOtherwise TASKCLUSTER_WORKER_LOCATION environment\nvariable will not be implicitly set in task commands.",
      "title": null,
      "type": "string"
    },
    "workerType": {
      "description": "This should match a worker_type managed by the\nspecified provisioner.",
      "title": null,
      "type": "string"
    },
    "workerTypeMetadata": {
      "description": "This arbitrary json blob will be included at the\ntop of each task log. Providing information here,\nsuch as a URL to the code/config used to set up the\nworker type will mean that people running tasks on\nthe worker type will have more information about how\nit was set up (for example what has been installed on\nthe machine).",
      "title": null,
      "type": "object"
    },
    "wstAudience": {
      "description": "The audience value for which to request websocktunnel\ncredentials, identifying a set of WST servers this\nworker could connect to.  Optional if not using websocktunnel\nto expose live logs.",
      "title": null,
      "type": "string"
    },
    "wstServerURL": {
      "description": "The URL of the websocktunnel server with which to expose\nlive logs.  Optional if not using websocktunnel to expose\nlive logs.",
      "format": "uri",
      "title": null,
      "type": "string"
    }
  },
  "required": [
    "accessToken",
    "clientId",
    "ed25519SigningKeyLocation",
    "rootURL",
    "workerId",
    "workerType"
  ],
  "title": "Config file for generic-worker",
  "type": "object"
}`
}