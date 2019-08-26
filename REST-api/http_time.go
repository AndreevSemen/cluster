package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func HttpTimeFunc(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(time.Now().Format(time.Stamp))); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/gossip/time", HttpTimeFunc).Methods("GET")

	log.Fatal(http.ListenAndServe(":" + os.Getenv("GOSSIP_TIME_PORT"), router))
}
