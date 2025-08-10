package output

import "testing"

func TestOutputBasic(t *testing.T) {
	// Just check that color functions do not panic
	Success("success")
	Error("error")
	Warn("warn")
	Info("info")
}
