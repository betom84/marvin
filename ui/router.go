package ui

import (
	"fmt"
	"marvin/metrics"
	"net/http"
)

type router struct {
	handlers map[string]http.HandlerFunc
	fallback http.HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]http.HandlerFunc),
	}
}

func (r *router) get(uri string, f http.HandlerFunc) {
	r.handlers[fmt.Sprintf("%s:%s", http.MethodGet, uri)] = f
}

func (r *router) post(uri string, f http.HandlerFunc) {
	r.handlers[fmt.Sprintf("%s:%s", http.MethodPost, uri)] = f
}

func (r *router) put(uri string, f http.HandlerFunc) {
	r.handlers[fmt.Sprintf("%s:%s", http.MethodPut, uri)] = f
}

func (r *router) serveHTTP(w http.ResponseWriter, req *http.Request) {
	if f, ok := r.handlers[fmt.Sprintf("%s:%s", req.Method, req.URL.Path)]; ok {
		metrics.Middleware(f).ServeHTTP(w, req)
		return
	}

	metrics.Middleware(r.fallback).ServeHTTP(w, req)
}
