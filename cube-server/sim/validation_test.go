package sim

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleValidation(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/validation?provider=azure&resource=aks", nil)
	w := httptest.NewRecorder()
	HandleValidation(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		t.Errorf("expected status OK or BadRequest, got %d", w.Code)
	}
}
