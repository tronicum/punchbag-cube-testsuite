package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Cube Server is running...")
	})
	fmt.Println("Cube Server binary is running...")
	http.ListenAndServe(":8080", nil)
}
