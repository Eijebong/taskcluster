// This source code file is AUTO-GENERATED by github.com/taskcluster/jsonschema2go

package tcgithub

import (
	tcclient "github.com/taskcluster/taskcluster/v49/clients/client-go"
)

type (
	Build struct {

		// The initial creation time of the build. This is when it became pending.
		Created tcclient.Time `json:"created"`

		// The GitHub webhook deliveryId. Extracted from the header 'X-GitHub-Delivery'
		//
		// One of:
		//   * GithubGUID
		//   * UnknownGithubGUID
		EventID string `json:"eventId"`

		// Type of Github event that triggered the build (i.e. push, pull_request.opened).
		EventType string `json:"eventType"`

		// Github organization associated with the build.
		//
		// Syntax:     ^([a-zA-Z0-9-_%]*)$
		// Min length: 1
		// Max length: 100
		Organization string `json:"organization"`

		// Github repository associated with the build.
		//
		// Syntax:     ^([a-zA-Z0-9-_%]*)$
		// Min length: 1
		// Max length: 100
		Repository string `json:"repository"`

		// Github revision associated with the build.
		//
		// Min length: 40
		// Max length: 40
		Sha string `json:"sha"`

		// Github status associated with the build.
		//
		// Possible values:
		//   * "pending"
		//   * "success"
		//   * "error"
		//   * "failure"
		State string `json:"state"`

		// Taskcluster task-group associated with the build.
		//
		// Syntax:     ^[A-Za-z0-9_-]{8}[Q-T][A-Za-z0-9_-][CGKOSWaeimquy26-][A-Za-z0-9_-]{10}[AQgw]$
		TaskGroupID string `json:"taskGroupId"`

		// The last updated of the build. If it is done, this is when it finished.
		Updated tcclient.Time `json:"updated"`
	}

	// A paginated list of builds
	BuildsResponse struct {

		// A simple list of builds.
		Builds []Build `json:"builds"`

		// Passed back from Azure to allow us to page through long result sets.
		ContinuationToken string `json:"continuationToken,omitempty"`
	}

	// Write a new comment on a GitHub Issue or Pull Request.
	// Full specification on [GitHub docs](https://developer.github.com/v3/issues/comments/#create-a-comment)
	CreateCommentRequest struct {

		// The contents of the comment.
		Body string `json:"body"`
	}

	// Create a commit status on GitHub.
	// Full specification on [GitHub docs](https://developer.github.com/v3/repos/statuses/#create-a-status)
	CreateStatusRequest struct {

		// A string label to differentiate this status from the status of other systems.
		Context string `json:"context,omitempty"`

		// A short description of the status.
		Description string `json:"description,omitempty"`

		// The state of the status.
		//
		// Possible values:
		//   * "pending"
		//   * "success"
		//   * "error"
		//   * "failure"
		State string `json:"state"`

		// The target URL to associate with this status. This URL will be linked from the GitHub UI to allow users to easily see the 'source' of the Status.
		Target_URL string `json:"target_url,omitempty"`
	}

	// The GitHub webhook deliveryId. Extracted from the header 'X-GitHub-Delivery'
	//
	// Syntax:     ^[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}$
	GithubGUID string

	// Any Taskcluster-specific Github repository information.
	RepositoryResponse struct {

		// True if integration is installed, False otherwise.
		Installed bool `json:"installed"`
	}

	// The GitHub webhook deliveryId. Extracted from the header 'X-GitHub-Delivery'
	//
	// Possible values:
	//   * "Unknown"
	UnknownGithubGUID string
)
