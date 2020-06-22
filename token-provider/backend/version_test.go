package backend

import "testing"

func TestGetVersion(t *testing.T) {

	expectedResult := versionInfo{
		Version:   "devel",
		CommitSHA: "unknown",
		BuildDate: "unknown",
	}
	if getVersionResponse() != expectedResult {
		t.Error("version output has changed")
	}
}
