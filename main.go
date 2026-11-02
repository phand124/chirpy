package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serv.ListenAndServe()
}
