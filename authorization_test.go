package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/taskcluster/httpbackoff"
	"github.com/taskcluster/slugid-go/slugid"
	"github.com/taskcluster/taskcluster-client-go/tcclient"
)

var (
	permCredentials = &tcclient.Credentials{
		ClientId:    os.Getenv("TASKCLUSTER_CLIENT_ID"),
		AccessToken: os.Getenv("TASKCLUSTER_ACCESS_TOKEN"),
	}
)

// Requires scope "auth:azure-table-access:fakeaccount/DuMmYtAbLe"
func sharedAccessSignature() string {
	return fmt.Sprintf(
		"http://localhost:60024/auth/v1/azure/%s/table/%s/read-write",
		"fakeaccount",
		"DuMmYtAbLe",
	)
}

type IntegrationTest func(t *testing.T, creds *tcclient.Credentials)

func skipIfNoPermCreds(t *testing.T) {
	if permCredentials.ClientId == "" {
		t.Skip("TASKCLUSTER_CLIENT_ID not set - skipping test")
	}
	if permCredentials.AccessToken == "" {
		t.Skip("TASKCLUSTER_ACCESS_TOKEN not set - skipping test")
	}
}

func testWithPermCreds(t *testing.T, test IntegrationTest) {
	skipIfNoPermCreds(t)
	test(t, permCredentials)
}

func testWithTempCreds(t *testing.T, test IntegrationTest) {
	skipIfNoPermCreds(t)
	tempCredentials, err := permCredentials.CreateTemporaryCredentials(1*time.Hour,
		"auth:azure-table-access:fakeaccount/DuMmYtAbLe",
		"queue:define-task:win-provisioner/win2008-worker",
		"queue:get-artifact:private/build/sources.xml",
		"queue:route:tc-treeherder.mozilla-inbound.*",
		"queue:route:tc-treeherder-stage.mozilla-inbound.*",
		"queue:task-priority:high",
		"test-worker:image:toastposter/pumpkin:0.5.6",
	)
	if err != nil {
		t.Fatalf("Could not generate temp credentials")
	}
	test(t, tempCredentials)
}

func TestBewit(t *testing.T) {
	test := func(t *testing.T, creds *tcclient.Credentials) {

		// Test setup
		routes := Routes(tcclient.ConnectionData{
			Credentials: creds,
		})
		req, err := http.NewRequest(
			"POST",
			"http://localhost:60024/bewit",
			bytes.NewBufferString("https://queue.taskcluster.net/v1/task/DD1kmgFiRMWTjyiNoEJIMA/runs/0/artifacts/private%2Fbuild%2Fsources.xml"),
		)
		if err != nil {
			log.Fatal(err)
		}
		res := httptest.NewRecorder()

		// Function to test
		routes.ServeHTTP(res, req)

		// Validate results
		if res.Code != 303 {
			t.Fatalf("Expected status code 303 but got %v", res.Code)
		}
		bewitUrl := res.Header().Get("Location")
		resp, _, err := httpbackoff.Get(bewitUrl)
		if err != nil {
			t.Fatalf("Exception thrown:\n%s", err)
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Exception thrown:\n%s", err)
		}
		if len(respBody) != 18170 {
			t.Logf("Response received:\n%s", string(respBody))
			t.Fatalf("Expected response body to be 18170 bytes, but was %v bytes", len(respBody))
		}
	}
	testWithPermCreds(t, test)
	testWithTempCreds(t, test)
}

func TestAuthorizationDelegate(t *testing.T) {
	test := func(name string, statusCode int, scopes []string) IntegrationTest {
		return func(t *testing.T, creds *tcclient.Credentials) {
			// Test setup
			routes := Routes(tcclient.ConnectionData{
				Authenticate: true,
				Credentials: &tcclient.Credentials{
					ClientId:         creds.ClientId,
					AccessToken:      creds.AccessToken,
					Certificate:      creds.Certificate,
					AuthorizedScopes: scopes,
				},
			})

			req, err := http.NewRequest(
				"GET",
				sharedAccessSignature(),
				// Note: we don't set body to nil as a server http request
				// cannot have a nil body. See:
				// https://golang.org/pkg/net/http/#Request
				new(bytes.Buffer),
			)
			if err != nil {
				log.Fatal(err)
			}
			res := httptest.NewRecorder()

			// Function to test
			routes.ServeHTTP(res, req)

			// Validate results

			if res.Code != statusCode {
				t.Logf("Part %s) Expected delgated request to fail with HTTP %v - but got HTTP %v", name, statusCode, res.Code)
				respBody, err := ioutil.ReadAll(res.Body)
				t.Logf("Headers: %s", res.Header())
				if err == nil {
					t.Logf("Response received:\n%s", string(respBody))
				}
				t.FailNow()
			}
		}
	}
	testWithPermCreds(t, test("A", 404, []string{"auth:azure-table-access:fakeaccount/DuMmYtAbLe"}))
	testWithTempCreds(t, test("B", 404, []string{"auth:azure-table-access:fakeaccount/DuMmYtAbLe"}))
	testWithPermCreds(t, test("C", 401, []string{"queue:get-artifact:private/build/sources.xml"}))
	testWithTempCreds(t, test("D", 401, []string{"queue:get-artifact:private/build/sources.xml"}))
}

func TestAPICallWithPayload(t *testing.T) {
	test := func(t *testing.T, creds *tcclient.Credentials) {

		// Test setup
		routes := Routes(tcclient.ConnectionData{
			Authenticate: true,
			Credentials:  creds,
		})
		taskId := slugid.Nice()
		taskGroupId := slugid.Nice()
		created := time.Now()
		deadline := created.AddDate(0, 0, 1)
		expires := deadline

		req, err := http.NewRequest(
			"POST",
			"http://localhost:60024/queue/v1/task/"+taskId+"/define",
			bytes.NewBufferString(
				`
{
  "provisionerId": "win-provisioner",
  "workerType": "win2008-worker",
  "schedulerId": "go-test-test-scheduler",
  "taskGroupId": "`+taskGroupId+`",
  "routes": [
    "tc-treeherder.mozilla-inbound.bcf29c305519d6e120b2e4d3b8aa33baaf5f0163",
    "tc-treeherder-stage.mozilla-inbound.bcf29c305519d6e120b2e4d3b8aa33baaf5f0163"
  ],
  "priority": "high",
  "retries": 5,
  "created": "`+tcclient.Time(created).String()+`",
  "deadline": "`+tcclient.Time(deadline).String()+`",
  "expires": "`+tcclient.Time(expires).String()+`",
  "scopes": [
    "test-worker:image:toastposter/pumpkin:0.5.6"
  ],
  "payload": {
    "features": {
      "relengApiProxy": true
    }
  },
  "metadata": {
    "description": "Stuff",
    "name": "[TC] Pete",
    "owner": "pmoore@mozilla.com",
    "source": "http://everywhere.com/"
  },
  "tags": {
    "createdForUser": "cbook@mozilla.com"
  },
  "extra": {
    "index": {
      "rank": 12345
    }
  }
}
`,
			),
		)
		if err != nil {
			log.Fatal(err)
		}
		res := httptest.NewRecorder()

		// Function to test
		routes.ServeHTTP(res, req)

		// Validate results
		if res.Code != 200 {
			t.Logf("Expected status code 200 but got %v", res.Code)
			respBody, err := ioutil.ReadAll(res.Body)
			t.Logf("Headers: %s", res.Header())
			if err == nil {
				t.Logf("Response received:\n%s", string(respBody))
			}
			t.FailNow()
		}
		t.Logf("Created task https://queue.taskcluster.net/v1/task/%v", taskId)
	}
	testWithPermCreds(t, test)
	testWithTempCreds(t, test)
}

func TestNon200HasErrorBody(t *testing.T) {
	test := func(t *testing.T, creds *tcclient.Credentials) {

		// Test setup
		routes := Routes(tcclient.ConnectionData{
			Authenticate: true,
			Credentials:  creds,
		})
		taskId := slugid.Nice()

		req, err := http.NewRequest(
			"POST",
			"http://localhost:60024/queue/v1/task/"+taskId+"/define",
			bytes.NewBufferString(
				`{"comment": "Valid json so that we hit endpoint, but should not result in http 200"}`,
			),
		)
		if err != nil {
			log.Fatal(err)
		}
		res := httptest.NewRecorder()

		// Function to test
		routes.ServeHTTP(res, req)

		respBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Could not read response body: %v", err)
		}
		// Validate results
		if res.Code != 400 {
			t.Logf("Expected status code 400 but got %v", res.Code)
			t.Logf("Headers: %s", res.Header())
			t.Logf("Response received:\n%s", string(respBody))
			t.FailNow()
		}
		if len(respBody) < 20 {
			t.Logf("Headers: %s", res.Header())
			t.Logf("Response received:\n%s", string(respBody))
			t.Log("Expected a response body (at least 20 bytes) with HTTP 400 error, but get less (or none).")
			t.FailNow()
		}

	}
	testWithPermCreds(t, test)
	testWithTempCreds(t, test)
}

func TestOversteppedScopes(t *testing.T) {
	test := func(t *testing.T, creds *tcclient.Credentials) {

		// Test setup
		routes := Routes(tcclient.ConnectionData{
			Authenticate: true,
			Credentials:  creds,
		})

		// This scope is not in the scopes of the temp credentials, which would
		// happen if a task declares a scope that the provisioner does not
		// grant.
		routes.Credentials.AuthorizedScopes = []string{"secrets:get:garbage/pmoore/foo"}

		req, err := http.NewRequest(
			"GET",
			"http://localhost:60024/secrets/v1/secret/garbage/pmoore/foo",
			new(bytes.Buffer),
		)
		if err != nil {
			log.Fatal(err)
		}
		res := httptest.NewRecorder()

		// Function to test
		routes.ServeHTTP(res, req)

		respBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Could not read response body: %v", err)
		}
		// Validate results
		if res.Code != 401 {
			t.Logf("Expected status code 401 but got %v", res.Code)
			t.Logf("Headers: %s", res.Header())
			t.Logf("Response received:\n%s", string(respBody))
			t.FailNow()
		}
		for headerKey, expectedHeaderValue := range map[string]string{
			"X-Taskcluster-Endpoint":          "https://secrets.taskcluster.net/v1/secret/garbage/pmoore/foo",
			"X-Taskcluster-Proxy-Version":     version,
			"X-Taskcluster-Authorized-Scopes": "[secrets:get:garbage/pmoore/foo]",
			"X-Taskcluster-Proxy-Temp-Scopes": "[auth:azure-table-access:fakeaccount/DuMmYtAbLe queue:define-task:win-provisioner/win2008-worker queue:get-artifact:private/build/sources.xml queue:route:tc-treeherder.mozilla-inbound.* queue:route:tc-treeherder-stage.mozilla-inbound.* queue:task-priority:high test-worker:image:toastposter/pumpkin:0.5.6]",
		} {
			actualHeaderValue := res.Header().Get(headerKey)
			if actualHeaderValue != expectedHeaderValue {
				t.Fatalf("Expected header %q to be %q but it was %q", headerKey, expectedHeaderValue, actualHeaderValue)
			}
		}
	}
	testWithTempCreds(t, test)
}
