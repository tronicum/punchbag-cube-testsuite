package importpkg

import (
	"os"
	"testing"
)

func TestImportCloudState(t *testing.T) {
	f := "testdata.json"
	os.WriteFile(f, []byte(`{"provider":"aws","region":"us-east-1"}`), 0644)
	defer os.Remove(f)
	state, err := ImportCloudState(f)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if state["provider"] != "aws" || state["region"] != "us-east-1" {
		t.Errorf("unexpected state: %v", state)
	}
}
