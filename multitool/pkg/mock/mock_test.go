package mock

import "testing"

func TestMockBasic(t *testing.T) {
	   // Just check that mock stores are instantiable
	   if NewHetznerObjectStorage() == nil {
			   t.Error("NewHetznerObjectStorage returned nil")
	   }
	   if NewIonosObjectStorage() == nil {
			   t.Error("NewIonosObjectStorage returned nil")
	   }
	   if NewStackITObjectStorage() == nil {
			   t.Error("NewStackITObjectStorage returned nil")
	   }
}
