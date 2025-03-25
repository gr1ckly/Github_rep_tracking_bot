package server

import (
	mux2 "github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	serverUrl string
	router    *mux2.Router
}

func BuildServer(serverUrl string, dbUrl string) (*Server, error) {
	router := mux2.NewRouter()
	router.HandleFunc("/user", handleAddUser).Methods("POST")
	router.HandleFunc("/user", handleDeleteUser).Methods("DELETE")
	router.HandleFunc("/repo", handleGetLinks).Methods("GET")
	router.HandleFunc("/repo", handleAddLink).Methods("POST")
	router.HandleFunc("/repo", handleDeleteLink).Methods("DELETE")
	return &Server{router: router}, nil
}

func (serv *Server) Start() error {
	err := http.ListenAndServe(serv.serverUrl, serv.router)
	return err
}

func (serv *Server) Stop() {

}
