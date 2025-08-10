package sim

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAks_POST(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/simulate/azure/aks", nil)
	w := httptest.NewRecorder()
	HandleAks(w, req)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusOK {
		t.Errorf("expected status BadRequest or OK, got %d", w.Code)
	}
}

func TestHandleAks_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/simulate/azure/aks?id=test", nil)
	w := httptest.NewRecorder()
	HandleAks(w, req)
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("expected status OK or NotFound, got %d", w.Code)
	}
}

func TestHandleAks_DELETE(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/simulate/azure/aks?id=test", nil)
	w := httptest.NewRecorder()
	HandleAks(w, req)
	if w.Code != http.StatusNoContent && w.Code != http.StatusNotFound {
		t.Errorf("expected status NoContent or NotFound, got %d", w.Code)
	}
}
