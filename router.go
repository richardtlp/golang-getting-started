package main

import (
	"net/http"
)

type Router struct{}

func (r Router) AddHandler(endpoint string, handler func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc("/notes/", handler)
}
