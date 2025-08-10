// Validation handler migrated from sim/
package sim

import (
	"net/http"
)

type ValidationResult struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

func HandleValidation(w http.ResponseWriter, r *http.Request) {
	// ...existing code from sim/validation.go...
}
