package client

import (
	"io"
	"net/http"
)

// ProxyOrSimHandler switches between proxy and simulation mode for S3 endpoints.
type ProxyOrSimHandler struct {
	ProxyURL   string
	SimHandler http.Handler
	Mode       string // "proxy" or "simulation"
}

// ServeHTTP forwards to proxy or simulation handler based on mode.
func (h *ProxyOrSimHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Mode == "proxy" {
		// Forward request to cube-server S3 endpoint
		proxyReq, err := http.NewRequest(r.Method, h.ProxyURL+r.URL.Path, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Proxy request error"))
			return
		}
		proxyReq.Header = r.Header.Clone()
		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("Proxy error: " + err.Error()))
			return
		}
		defer resp.Body.Close()
		for k, v := range resp.Header {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, resp.Body)
		return
	}
	// Simulation mode: use in-memory S3 mock
	h.SimHandler.ServeHTTP(w, r)
}
