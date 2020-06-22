package storage

import (
	"testing"
)

func TestRandomGenerator(t *testing.T) {
	raw, err := generateRandomStringForKey(apiKeyLength)
	if err != nil {
		t.Fatalf("failed to generate random string: %v", err)
	}
	t.Logf("generated random string: %s", raw)
}
