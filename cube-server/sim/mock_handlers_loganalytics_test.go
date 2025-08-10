package sim

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleLogAnalytics_POST(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/simulate/azure/loganalytics", nil)
	w := httptest.NewRecorder()
	HandleLogAnalytics(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected status Created, got %d", w.Code)
	}
}

func TestHandleLogAnalytics_GET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/simulate/azure/loganalytics?id=test", nil)
	w := httptest.NewRecorder()
	HandleLogAnalytics(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %d", w.Code)
	}
}

func TestHandleLogAnalytics_DELETE(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/simulate/azure/loganalytics?id=test", nil)
	w := httptest.NewRecorder()
	HandleLogAnalytics(w, req)
	if w.Code != http.StatusNoContent && w.Code != http.StatusNotFound {
		t.Errorf("expected status NoContent or NotFound, got %d", w.Code)
	}
}

func TestHandleLogAnalytics_DELETE_NoID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/simulate/azure/loganalytics", nil)
	w := httptest.NewRecorder()
	HandleLogAnalytics(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status NotFound, got %d", w.Code)
	}
}
