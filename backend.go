package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Backend struct {
	Address string
}

func (b *Backend) Start() {
	handler := func(w http.ResponseWriter, r *http.Request) {

		// Simular un tiempo de procesamiento aleatorio (200–800 ms)
		delay := time.Duration(rand.Intn(600)+200) * time.Millisecond
		time.Sleep(delay)

		fmt.Fprintf(w,
			"Respuesta desde backend %s (delay: %v)\n",
			b.Address,
			delay,
		)
	}

	srv := &http.Server{
		Addr:    b.Address,
		Handler: http.HandlerFunc(handler),
	}

	log.Printf("[Backend %s] Iniciado", b.Address)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("[Backend %s] Error: %v", b.Address, err)
		}
	}()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
