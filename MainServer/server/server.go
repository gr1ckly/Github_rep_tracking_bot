package server

import (
	"Crypto_Bot/MainServer/storage"
	mux2 "github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	serverUrl           string
	router              *mux2.Router
	chatStore           *storage.ChatStore
	repoStore           *storage.RepoStore
	chatRepoRecordStore *storage.ChatRepoRecordStore
}

func BuildServer(serverUrl string, chatStore *storage.ChatStore, repoStore *storage.RepoStore, chatRepoRecordStore *storage.ChatRepoRecordStore) *Server {
	router := mux2.NewRouter()
	router.HandleFunc("/user", handleAddUser).Methods("POST")
	router.HandleFunc("/user", handleDeleteUser).Methods("DELETE")
	router.HandleFunc("/repo", handleGetLinks).Methods("GET")
	router.HandleFunc("/repo", handleAddLink).Methods("POST")
	router.HandleFunc("/repo", handleDeleteLink).Methods("DELETE")
	return &Server{serverUrl: serverUrl, router: router, chatStore: chatStore, repoStore: repoStore, chatRepoRecordStore: chatRepoRecordStore}
}

func (serv *Server) Start() error {
	err := http.ListenAndServe(serv.serverUrl, serv.router)
	return err
}

func (serv *Server) Stop() {

}
