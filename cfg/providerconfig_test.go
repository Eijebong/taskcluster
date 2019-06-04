package cfg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"
)

func TestNoProviderType(t *testing.T) {
	var pc ProviderConfig

	err := yaml.Unmarshal([]byte(`{}`), &pc)
	if err == nil {
		t.Fatalf("did not fail")
	}

	assert.Equal(t, fmt.Errorf("provider config must have a `providerType` property"), err, "should have errored")
}

func TestNonStringProviderType(t *testing.T) {
	var pc ProviderConfig

	err := yaml.Unmarshal([]byte(`{"providerType": false}`), &pc)
	if err == nil {
		t.Fatalf("did not fail")
	}

	assert.Equal(t, fmt.Errorf("provider config's `providerType` property must be a string"), err, "should have errored")
}

func TestOK(t *testing.T) {
	var pc ProviderConfig
	err := yaml.Unmarshal([]byte(`{"providerType": "something", "value": "sure"}`), &pc)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}

	assert.Equal(t, map[string]interface{}{"value": "sure"}, pc.data, "did not get expected config")
}

func TestUnpack(t *testing.T) {
	type mypc struct {
		Value   int
		Another string `provider:"anotherValue"`
	}

	var pc ProviderConfig
	yaml.Unmarshal([]byte(`{"providerType": "x", "Value": 10, "anotherValue": "hi"}`), &pc)

	var c mypc
	err := pc.Unpack(&c)
	if err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}
	assert.Equal(t, mypc{10, "hi"}, c, "unpacked values correctly")
}

func TestUnpackMissing(t *testing.T) {
	type mypc struct {
		Value int
	}

	var pc ProviderConfig
	yaml.Unmarshal([]byte(`{}`), &pc)

	var c mypc
	err := pc.Unpack(&c)
	if err == nil {
		t.Fatalf("failed to fail")
	}
}

func TestUnpackWrongType(t *testing.T) {
	type mypc struct {
		Value int
	}

	var pc ProviderConfig
	yaml.Unmarshal([]byte(`{"Value": "yo"}`), &pc)

	var c mypc
	err := pc.Unpack(&c)
	if err == nil {
		t.Fatalf("failed to fail")
	}
}
