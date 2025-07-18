package client

import (
	"io"
	"net/http"
)

// S3Proxy abstracts S3 API proxying for different providers.
type S3Proxy interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// HetznerS3Proxy proxies S3 requests to Hetzner (AWS S3 compatible, but with endpoint/location tweaks).
type HetznerS3Proxy struct {
	ProxyURL string
}

func (h *HetznerS3Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Forward request to Hetzner S3 endpoint, adjusting as needed
	proxyReq, err := http.NewRequest(r.Method, h.ProxyURL+r.URL.Path, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Hetzner S3 proxy request error"))
		return
	}
	proxyReq.Header = r.Header.Clone()
	// Hetzner S3 may require region/location header tweaks here
	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Hetzner S3 proxy error: " + err.Error()))
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
}

// AzureBlobProxy proxies requests to Azure Blob Storage (using Azure SDK or REST API).
type AzureBlobProxy struct {
	ProxyURL string
}

func (a *AzureBlobProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Use Azure SDK or REST API to forward requests
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Azure Blob proxy not implemented"))
}

// S3SimOrProxyHandler switches between simulation and proxy for S3 endpoints, per provider.
type S3SimOrProxyHandler struct {
	Mode      string // "proxy" or "simulation"
	Sim       http.Handler
	Proxy     S3Proxy
}

func (h *S3SimOrProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.Mode == "proxy" {
		h.Proxy.ServeHTTP(w, r)
		return
	}
	h.Sim.ServeHTTP(w, r)
}
