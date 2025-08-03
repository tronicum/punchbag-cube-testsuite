package cmd

import "testing"

func TestRootBasic(t *testing.T) {
	// Test that Execute() does not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute panicked: %v", r)
		}
	}()
	Execute()
}
