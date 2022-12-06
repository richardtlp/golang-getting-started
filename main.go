package main

import (
	"net/http"
	"fmt"
)

func setupRouter() {
	notesHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			GetNotesHandler(w, r)
		case "POST":
			CreateNotesHandler(w, r)
		case "PUT":
			UpdateNotesHandler(w, r)
		case "DELETE":
			DeleteNotesHandler(w, r)
		default:
			DefaultErrorHandler(w, r)
		}
	}
	router := Router{}
	router.AddHandler("/notes/", notesHandler)
}

func launchServer(address string) {
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("http server crashed: %s", err.Error())
	}
}

func main() {
	setupRouter()
	launchServer(":8080")
}
