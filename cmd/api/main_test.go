package main

import (
	"net/http"
	"testing"
)

func TestNewHTTPServer(t *testing.T) {

	handler := http.NewServeMux()
	server := newHTTPServer(":8080", handler)

	if server.Addr != ":8080" {
		t.Errorf("expected address to be :8080, got %q", server.Addr)
	}

	if server.Handler != handler {
		t.Error("expected handler to be the same as provided")
	}

	if server.ReadHeaderTimeout != serverReadHeaderTimeout {
		t.Errorf("expected ReadHeaderTimeout to be %v, got %v", serverReadHeaderTimeout, server.ReadHeaderTimeout)
	}

	if server.ReadTimeout != serverReadTimeout {
		t.Errorf("expected ReadTimeout to be %v, got %v", serverReadTimeout, server.ReadTimeout)
	}

	if server.WriteTimeout != serverWriteTimeout {
		t.Errorf("expected WriteTimeout to be %v, got %v", serverWriteTimeout, server.WriteTimeout)
	}

	if server.IdleTimeout != serverIdleTimeout {
		t.Errorf("expected IdleTimeout to be %v, got %v", serverIdleTimeout, server.IdleTimeout)
	}
}
