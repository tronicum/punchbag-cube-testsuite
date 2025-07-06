package main

import (
	"fmt"
	"net/http"
)

// func main() { ... } // moved to cmd/cube-server/main.go

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Cube Server is running...")
}

func init() {
	http.HandleFunc("/api/v1/status", statusHandler)
}
