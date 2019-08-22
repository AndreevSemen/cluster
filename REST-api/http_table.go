package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type ServerInfo struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	State string `json:"state"`
}

func IsValidState(info ServerInfo) bool {
	return info.State == "alive"    ||
		   info.State == "onload"   ||
		   info.State == "inactive" ||
		   info.State == "dead"
}

type ServerInfoArray struct {
	Servers []ServerInfo `json:"servers"`
}

var database ServerInfoArray;

func GETALL(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	if  err := json.NewEncoder(w).Encode(database); err != nil {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
}

func POST(w http.ResponseWriter, r *http.Request) {
	var requestedServer ServerInfo


	if err := json.NewDecoder(r.Body).Decode(&requestedServer); err != nil || !IsValidState(requestedServer) {
		w.WriteHeader(400)
		return
	}

	updated := false
	for i := 0; i < len(database.Servers); i++ {
		if database.Servers[i].IP == requestedServer.IP &&
		   database.Servers[i].Port == requestedServer.Port {
			database.Servers[i] = requestedServer
			updated = true
			break
		}
	}

	if !updated {
		database.Servers = append(database.Servers, requestedServer)
	}

	var neighbourServers ServerInfoArray
	for _, server := range database.Servers {
		if server.IP == requestedServer.IP && server.Port != requestedServer.Port {
			neighbourServers.Servers = append(neighbourServers.Servers, server)
		}
	}

	if err := json.NewEncoder(w).Encode(neighbourServers); err != nil {
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(200)
}

func DROP(w http.ResponseWriter, r *http.Request) {
	database.Servers = nil
	w.WriteHeader(200)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/gossip/table", GETALL).Methods("GET")
	router.HandleFunc("/gossip/table/update", POST).Methods("POST")
	router.HandleFunc("/gossip/table/drop", DROP).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":80", router))
}
