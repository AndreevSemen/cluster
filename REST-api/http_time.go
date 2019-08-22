package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func TIME(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(time.Now().Format(time.Stamp))); err != nil {
		w.WriteHeader(505)
	}
	w.WriteHeader(200)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/gossip/time", TIME).Methods("GET")

	log.Fatal(http.ListenAndServe(":82", router))
}
