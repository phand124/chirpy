package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serv.ListenAndServe()
}
