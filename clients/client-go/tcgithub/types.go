// This source code file is AUTO-GENERATED by github.com/taskcluster/jsonschema2go

package tcgithub

import (
	"encoding/json"

	tcclient "github.com/taskcluster/taskcluster/v87/clients/client-go"
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

		// Associated pull request number for 'pull_request' events.
		PullRequestNumber int64 `json:"pullRequestNumber,omitempty"`

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
		SHA string `json:"sha"`

		// Github status associated with the build.
		//
		// Possible values:
		//   * "pending"
		//   * "success"
		//   * "error"
		//   * "failure"
		//   * "cancelled"
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

	// Emulate one of the github events with mocked payload.
	// Some of the events have sub-actions, that can be specified.
	// Event type names follow the `tasks_for` naming convention.
	IssueCommentEvents struct {

		// Possible values:
		//   * "created"
		//   * "edited"
		Action string `json:"action"`

		// Additional data to be mixed to the mocked event object.
		// This can be used to set some specific properties of the event or override the existing ones.
		// For example:
		//   "ref": "refs/heads/main"
		//   "before": "000"
		//   "after": "111"
		// To make sure which properties are available for each event type,
		// please refer to the github [documentation](https://docs.github.com/en/webhooks-and-events/webhooks/webhook-events-and-payloads)
		//
		// Additional properties allowed
		Overrides json.RawMessage `json:"overrides,omitempty"`

		// Possible values:
		//   * "github-issue-comment"
		Type string `json:"type"`
	}

	// .taskcluster.yml supports `github-pull-request` and `github-pull-request-untrusted` events.
	// The difference is that `github-pull-request-untrusted` will use different set of scopes.
	// See [RFC 175](https://github.com/taskcluster/taskcluster-rfcs/blob/main/rfcs/0175-restricted-pull-requests.md)
	PullRequestEvents struct {

		// Possible values:
		//   * "opened"
		//   * "synchronize"
		//   * "reopened"
		//   * "assigned"
		//   * "auto_merge_disabled"
		//   * "auto_merge_enabled"
		//   * "closed"
		//   * "converted_to_draft"
		//   * "dequeued"
		//   * "edited"
		//   * "enqueued"
		//   * "labeled"
		//   * "ready_for_review"
		//   * "review_requested"
		//   * "review_request_removed"
		//   * "unassigned"
		//   * "unlabeled"
		Action string `json:"action"`

		// Additional data to be mixed to the mocked event object.
		// This can be used to set some specific properties of the event or override the existing ones.
		// For example:
		//   "ref": "refs/heads/main"
		//   "before": "000"
		//   "after": "111"
		// To make sure which properties are available for each event type,
		// please refer to the github [documentation](https://docs.github.com/en/webhooks-and-events/webhooks/webhook-events-and-payloads)
		//
		// Additional properties allowed
		Overrides json.RawMessage `json:"overrides,omitempty"`

		// Possible values:
		//   * "github-pull-request"
		//   * "github-pull-request-untrusted"
		Type string `json:"type"`
	}

	// Github sends `push` event for commits and for tags.
	// To distinguish between those two, the `ref` property is used.
	// If you want to mock a tag push, please specify `ref` property in the overrides:
	// "ref": "refs/tags/v1.0.0"
	PushEvents struct {

		// Additional data to be mixed to the mocked event object.
		// This can be used to set some specific properties of the event or override the existing ones.
		// For example:
		//   "ref": "refs/heads/main"
		//   "before": "000"
		//   "after": "111"
		// To make sure which properties are available for each event type,
		// please refer to the github [documentation](https://docs.github.com/en/webhooks-and-events/webhooks/webhook-events-and-payloads)
		//
		// Additional properties allowed
		Overrides json.RawMessage `json:"overrides,omitempty"`

		// Possible values:
		//   * "github-push"
		Type string `json:"type"`
	}

	// Emulate one of the github events with mocked payload.
	// Some of the events have sub-actions, that can be specified.
	// Event type names follow the `tasks_for` naming convention.
	ReleaseEvents struct {

		// Possible values:
		//   * "published"
		//   * "unpublished"
		//   * "created"
		//   * "edited"
		//   * "deleted"
		//   * "prereleased"
		//   * "released"
		Action string `json:"action"`

		// Additional data to be mixed to the mocked event object.
		// This can be used to set some specific properties of the event or override the existing ones.
		// For example:
		//   "ref": "refs/heads/main"
		//   "before": "000"
		//   "after": "111"
		// To make sure which properties are available for each event type,
		// please refer to the github [documentation](https://docs.github.com/en/webhooks-and-events/webhooks/webhook-events-and-payloads)
		//
		// Additional properties allowed
		Overrides json.RawMessage `json:"overrides,omitempty"`

		// Possible values:
		//   * "github-release"
		Type string `json:"type"`
	}

	// Render .taskcluster.yml for one of the supported events.
	//
	// Read more about the `.taskcluster.yml` file format in
	// [documentation](https://docs.taskcluster.net/docs/reference/integrations/github/taskcluster-yml-v1)
	RenderTaskclusterYmlInput struct {

		// The contents of the .taskcluster.yml file.
		Body string `json:"body"`

		// Emulate one of the github events with mocked payload.
		// Some of the events have sub-actions, that can be specified.
		// Event type names follow the `tasks_for` naming convention.
		//
		// One of:
		//   * PushEvents
		//   * PullRequestEvents
		//   * ReleaseEvents
		//   * IssueCommentEvents
		FakeEvent json.RawMessage `json:"fakeEvent"`

		// Syntax:     ^[-a-zA-Z0-9]{1,39}$
		Organization string `json:"organization,omitempty"`

		// Syntax:     ^[-a-zA-Z0-9_.]{1,100}$
		Repository string `json:"repository,omitempty"`
	}

	// Rendered .taskcluster.yml output.
	RenderTaskclusterYmlOutput struct {

		// Scopes that will be used by the github client to create tasks.
		// Those are different that the scopes inside the tasks itself.
		//
		// Array items:
		Scopes []string `json:"scopes"`

		// Rendered tasks objects.
		// Those objects not guaranteed to produce valid task definitions
		// that conform to the json schema.
		//
		// Array items:
		// Additional properties allowed
		Tasks []json.RawMessage `json:"tasks"`
	}

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
