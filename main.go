package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {

	// ---------------- BACKENDS ---------------- //
	backends := []string{
		"localhost:9001",
		"localhost:9002",
		"localhost:9003",
	}

	for _, addr := range backends {
		b := Backend{Address: addr}
		b.Start()
	}

	// ---------------- BALANCER ---------------- //
	balancer := NewLoadBalancer(backends)

	// Canal de entrada de peticiones
	requests := make(chan Request, 100)

	// Contadores por worker
	workerCount := make([]int64, 4)

	// Crear worker pool
	for i := 0; i < 4; i++ {
		workerID := i
		go func() {
			for req := range requests {
				// Incrementar contador
				atomic.AddInt64(&workerCount[workerID], 1)

				// Procesar petición
				balancer.Forward(context.Background(), req)

				// Notificar finalización
				if req.Done != nil {
					close(req.Done)
				}
			}
		}()
	}

	// ---------------- HANDLER PRINCIPAL ---------------- //
	handler := func(w http.ResponseWriter, r *http.Request) {
		done := make(chan struct{})
		requests <- Request{Writer: w, Request: r, Done: done}
		<-done
	}

	// ---------------- ENDPOINT DE ESTADÍSTICAS ---------------- //
	statsHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		stats := map[string]int64{
			"worker_0": atomic.LoadInt64(&workerCount[0]),
			"worker_1": atomic.LoadInt64(&workerCount[1]),
			"worker_2": atomic.LoadInt64(&workerCount[2]),
			"worker_3": atomic.LoadInt64(&workerCount[3]),
		}

		json.NewEncoder(w).Encode(stats)
	}

	// ---------------- REGISTRO DE ENDPOINTS ---------------- //
	http.HandleFunc("/", handler)
	http.HandleFunc("/stats", statsHandler)

	log.Println("Load balancer escuchando en :8080")

	// Arrancar servidor
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
