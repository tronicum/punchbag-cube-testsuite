package client

import "testing"

func TestClientBasic(t *testing.T) {
	client := NewAPIClient("http://localhost:1234")
	if client == nil {
		t.Error("NewAPIClient returned nil")
	}
}
