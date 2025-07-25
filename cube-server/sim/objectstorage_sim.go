// S3/Hetzner object storage simulation server for cube-server
package sim

import (
	"fmt"
	"net/http"
	"os"
)

// HetznerS3MockHandler exposes the HetznerS3Mock as an HTTP handler for simulation
func HetznerS3MockHandler() http.Handler {
	// Import the mock from multitool/pkg/client (move if needed)
	// For now, just return NotImplemented
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Hetzner S3 simulation not yet implemented"))
	})
}

// StartObjectStorageSimulationServer starts a simulation server for S3-like APIs
func StartObjectStorageSimulationServer(provider, port string) {
	switch provider {
	case "hetzner":
		fmt.Printf("Starting Hetzner S3 simulation on http://localhost:%s\n", port)
		http.Handle("/", HetznerS3MockHandler())
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Provider not supported for simulation.")
		os.Exit(1)
	}
}
