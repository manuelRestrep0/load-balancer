package main

import (
	"context"
	"io"
	"net/http"
	"sync/atomic"
)

type LoadBalancer struct {
	backends []string
	counter  uint64
}

func NewLoadBalancer(backends []string) *LoadBalancer {
	return &LoadBalancer{backends: backends}
}

func (lb *LoadBalancer) getNextBackend() string {
	// Round-robin
	idx := atomic.AddUint64(&lb.counter, 1)
	return lb.backends[(idx-1)%uint64(len(lb.backends))]
}

func (lb *LoadBalancer) Forward(ctx context.Context, req Request) {
	backend := lb.getNextBackend()

	// Crear petición hacia el backend seleccionado
	proxyReq, err := http.NewRequestWithContext(
		ctx,
		req.Request.Method,
		"http://"+backend+req.Request.URL.Path,
		req.Request.Body,
	)
	if err != nil {
		http.Error(req.Writer, "No se pudo crear la petición", http.StatusInternalServerError)
		return
	}

	// Copiar headers
	proxyReq.Header = req.Request.Header.Clone()

	// Hacer la petición al backend
	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		http.Error(req.Writer, "Backend no disponible", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copiar headers y body al cliente
	for k, v := range resp.Header {
		req.Writer.Header()[k] = v
	}
	req.Writer.WriteHeader(resp.StatusCode)
	io.Copy(req.Writer, resp.Body)
}
