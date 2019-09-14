package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"sort"
)

var states = []string{"alive", "dead", "inactive", "onload"}

type ServerInfo struct {
	IP    string `json:"ip"`
	Port  int    `json:"port"`
	State string `json:"state"`
}
func (info ServerInfo) IsValidState() bool {
	result := sort.SearchStrings(states, info.State)
	return result < len(states) && info.State == states[result]
}

type ServerInfoArray struct {
	Servers []ServerInfo `json:"servers"`
}


type ServerMap struct {
	data map[string]map[int]string
}
func makeServerMap() ServerMap {
	return ServerMap{make(map[string]map[int]string)}
}
func (s ServerMap) GetServersArray() ServerInfoArray {
	var arr ServerInfoArray

	for ip := range s.data {
		for port := range s.data[ip] {
			arr.Servers = append(arr.Servers, ServerInfo{ip,port, s.data[ip][port]})
		}
	}

	return arr
}
func (s ServerMap) NeighbourServersArray(currSrv ServerInfo) ServerInfoArray {
	var arr ServerInfoArray

	for port := range s.data[currSrv.IP] {
		if port != currSrv.Port {
			arr.Servers = append(arr.Servers, ServerInfo{currSrv.IP, port, s.data[currSrv.IP][port]})
		}
	}

	return arr
}
func (s ServerMap) AddServerInfo(srv ServerInfo) {
	if s.data[srv.IP] == nil {
		s.data[srv.IP] = make(map[int]string)
	}

	s.data[srv.IP][srv.Port] = srv.State
}
func (s ServerMap) Empty() {
	s.data = make(map[string]map[int]string)
}

var database = makeServerMap()

func HttpGetAll(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	if  err := json.NewEncoder(w).Encode(database.GetServersArray()); err != nil {
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
	if !requestedServer.IsValidState() {
		w.Write([]byte("sent server with invalid state"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	database.AddServerInfo(requestedServer)

	if err := json.NewEncoder(w).Encode(database.NeighbourServersArray(requestedServer)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HttpDropDatabase(w http.ResponseWriter, r *http.Request) {
	database.Empty()
	w.WriteHeader(http.StatusOK)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/gossip/table", HttpGetAll).Methods("GET")
	router.HandleFunc("/gossip/table/update", HttpPost).Methods("POST")
	router.HandleFunc("/gossip/table/drop", HttpDropDatabase).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":" + os.Getenv("GOSSIP_TABLE_PORT"), router))
}