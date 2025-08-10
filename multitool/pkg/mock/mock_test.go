package mock

import "testing"

func TestMockBasic(t *testing.T) {
	// Just check that mock stores are instantiable
	if NewIonosObjectStorage() == nil {
		t.Error("NewIonosObjectStorage returned nil")
	}
	if NewStackITObjectStorage() == nil {
		t.Error("NewStackITObjectStorage returned nil")
	}
}
