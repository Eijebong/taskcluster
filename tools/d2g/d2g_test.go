package d2g_test

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/mcuadros/go-defaults"
	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"

	d2g "github.com/taskcluster/taskcluster/v53/tools/d2g"
	"github.com/taskcluster/taskcluster/v53/tools/d2g/dockerworker"
)

func ExampleScopes_mixture() {
	dwScopes := []string{
		"foo",
		"bar:dog",
		"cat:docker-worker:feet",
		"docker-worker",
		"docker-worker:monkey",
		"generic-worker:teapot",
		"docker-worker:docker-worker:potato",
		"docker-worker:capability:device:loopbackVideo",
		"docker-worker:capability:device:loopbackVideo:",
		"docker-worker:capability:device:loopbackVideo:x/y/z",
		"docker-worker:capability:device:kvm:x/y/z",
	}
	gwScopes := d2g.Scopes(dwScopes)
	for _, s := range gwScopes {
		fmt.Printf("\t%#v\n", s)
	}

	// Output:
	// 	"foo"
	// 	"bar:dog"
	// 	"cat:docker-worker:feet"
	// 	"docker-worker"
	// 	"generic-worker:monkey"
	// 	"generic-worker:teapot"
	// 	"generic-worker:docker-worker:potato"
	// 	"generic-worker:loopback-video:*"
	// 	"generic-worker:loopback-video:"
	// 	"generic-worker:loopback-video:x/y/z"
	// 	"generic-worker:capability:device:kvm:x/y/z"
}

// TestDataTestCases runs all the test cases found in directory testdata/testcases.
func TestDataTestCases(t *testing.T) {
	b := yamlToJSON(t, "schemas/test_suites.yml")
	schemaLoader := gojsonschema.NewBytesLoader(b)
	// Enumerate all test suites under testdata/testcases, and execute a subtest for each suite
	err := filepath.WalkDir(
		"testdata/testcases",
		func(path string, file fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if file.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".yml" {
				t.Logf("Skipping %v", path)
				return nil
			}
			t.Run(
				path,
				testSuite(schemaLoader, path),
			)
			return nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}

// testSuite returns a go test the given testSuite
func testSuite(schemaLoader gojsonschema.JSONLoader, path string) func(t *testing.T) {
	return func(t *testing.T) {
		validateTestSuite(t, schemaLoader, path)
		// Iterate through test cases in the test suite, and execute a subtest for each test case
		var d D2GTestCases
		defaults.SetDefaults(&d)
		unmarshalYAML(t, &d, path)
		// TODO: check if this works as expected (create explicit test cases for this)
		// setting defaults and unmarshaling a second time
		// because this seems to fix the issue with defaults
		// not being applied to slices
		defaults.SetDefaults(&d)
		unmarshalYAML(t, &d, path)
		for _, tc := range d.TestSuite.Tests {
			t.Run(
				tc.Name,
				tc.TestCase(),
			)
		}
	}
}

func (tc TestCase) TestCase() func(t *testing.T) {
	return func(t *testing.T) {
		// need to convert types using unsafe package because ... https://github.com/golang/go/issues/58965
		dwPayload := *(*dockerworker.DockerWorkerPayload)(unsafe.Pointer(&tc.DockerWorkerTaskPayload))
		actualGWPayload, err := d2g.Convert(&dwPayload)
		if err != nil {
			t.Fatalf("Cannot convert Docker Worker payload %#v to Generic Worker payload: %s", dwPayload, err)
		}
		formattedActualGWPayload, err := json.MarshalIndent(*actualGWPayload, "", "  ")
		if err != nil {
			t.Fatalf("Cannot convert Generic Worker payload %#v to JSON: %s", *actualGWPayload, err)
		}
		formattedExpectedGWPayload, err := json.MarshalIndent(tc.GenericWorkerTaskPayload, "", "  ")
		if err != nil {
			t.Fatalf("Cannot convert Generic Worker payload %#v to JSON: %s", tc.GenericWorkerTaskPayload, err)
		}
		if string(formattedExpectedGWPayload) != string(formattedActualGWPayload) {
			t.Fatalf("Converted task does not match expected value.\nExpected:%v\nActual:%v", string(formattedExpectedGWPayload), string(formattedActualGWPayload))
		}
	}
}

func yamlToJSON(t *testing.T, path string) []byte {
	t.Helper()
	yml, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	j, err := yaml.YAMLToJSON(yml)
	if err != nil {
		t.Fatal(err)
	}
	// Takes json []byte input, unmarshals and then marshals, in order to get a
	// canonical representation of json (i.e. formatted with objects ordered).
	// Ugly and perhaps inefficient, but effective! :p
	tmpObj := new(interface{})
	err = json.Unmarshal(j, &tmpObj)
	if err != nil {
		t.Fatal(err)
	}
	formatted, err := json.MarshalIndent(&tmpObj, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return formatted
}

func unmarshalYAML(t *testing.T, dest interface{}, path string) {
	t.Helper()
	j := yamlToJSON(t, path)
	err := json.Unmarshal(j, dest)
	if err != nil {
		t.Fatal(err)
	}
}

func validateTestSuite(t *testing.T, schemaLoader gojsonschema.JSONLoader, path string) {
	t.Helper()
	b := yamlToJSON(t, path)
	documentLoader := gojsonschema.NewBytesLoader(b)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Valid() {
		t.Logf("The document is not valid:")
		for _, desc := range result.Errors() {
			t.Logf("- %s", desc)
		}
		t.FailNow()
	}
}