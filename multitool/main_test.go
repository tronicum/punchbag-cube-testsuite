package main

import "testing"

func TestMainBasic(t *testing.T) {
	   // Just check that main() runs without panic
	   defer func() {
			   if r := recover(); r != nil {
					   t.Errorf("main panicked: %v", r)
			   }
	   }()
	   main()
}
