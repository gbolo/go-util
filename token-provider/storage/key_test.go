package storage

import (
	"testing"
)

func TestApiKeyGenerator(t *testing.T) {
	raw, err := GenerateApiKey()
	if err != nil {
		t.Fatalf("failed to generate api key: %v", err)
	}
	t.Logf("generated api key: %s", raw)
	if !raw.HasValidFormat() {
		t.Fatalf("api key does not have valid format: %s", raw)
	}
	t.Logf("key prefix: %s", raw.GetPrefix())
	t.Logf("key hash: %s", raw.GetHash())
}
