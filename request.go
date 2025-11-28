package main

import "net/http"

type Request struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Done    chan struct{}
}
