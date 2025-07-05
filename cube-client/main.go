package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Cube Client binary is running...")
	resp, err := http.Get("http://localhost:8080/api/v1/status")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Server response:", resp.Status)
}
