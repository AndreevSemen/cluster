package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"sort"
)

type ServerInfo struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	State string `json:"state"`
}

var states = []string{"alive", "dead", "inactive", "onload"}

type ServerInfoArray struct {
	Servers []ServerInfo `json:"servers"`
}

func IsValidState(info ServerInfo) bool {
	result := sort.SearchStrings(states, info.State)
	return result < len(states) && info.State == states[result]
}

type ServerMap = map[string]map[int]string

var database = make(ServerMap)

func GetServersArray(database* ServerMap) ServerInfoArray {
	var arr ServerInfoArray

	for ip := range *database {
		for port := range (*database)[ip] {
			arr.Servers = append(arr.Servers, ServerInfo{ip,port, (*database)[ip][port]})
		}
	}

	return arr
}

func NeighbourServersArray(database *ServerMap, currSrv ServerInfo) ServerInfoArray {
	var arr ServerInfoArray

	for port := range (*database)[currSrv.IP] {
		if port != currSrv.Port {
			arr.Servers = append(arr.Servers, ServerInfo{currSrv.IP, port, (*database)[currSrv.IP][port]})
		}
	}

	return arr
}

func AddServerInfo(database *ServerMap, srv ServerInfo) {
	if (*database)[srv.IP] == nil {
		(*database)[srv.IP] = make(map[int]string)
	}

	(*database)[srv.IP][srv.Port] = srv.State
}

func Empty(database *ServerMap) {
	*database = make(ServerMap)
}

func HttpGetAll(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	if  err := json.NewEncoder(w).Encode(GetServersArray(&database)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HttpPost(w http.ResponseWriter, r *http.Request) {
	var requestedServer ServerInfo

	if err := json.NewDecoder(r.Body).Decode(&requestedServer); err != nil {
		w.Write([]byte("sent invalid json"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !IsValidState(requestedServer) {
		w.Write([]byte("sent server with invalid state"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	AddServerInfo(&database, requestedServer)

	if err := json.NewEncoder(w).Encode(NeighbourServersArray(&database, requestedServer)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HttpDropDatabase(w http.ResponseWriter, r *http.Request) {
	Empty(&database)
	w.WriteHeader(http.StatusOK)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/gossip/table", HttpGetAll).Methods("GET")
	router.HandleFunc("/gossip/table/update", HttpPost).Methods("POST")
	router.HandleFunc("/gossip/table/drop", HttpDropDatabase).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":" + os.Getenv("GOSSIP_TABLE_PORT"), router))
}