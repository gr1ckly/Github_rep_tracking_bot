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

func (serv *Server) handleAddUser(w http.ResponseWriter, req *http.Request) {

}

func (serv *Server) handleDeleteUser(w http.ResponseWriter, req *http.Request) {

}

func (serv *Server) handleGetLinks(w http.ResponseWriter, req *http.Request) {

}

func (serv *Server) handleAddLink(w http.ResponseWriter, req *http.Request) {

}

func (serv *Server) handleDeleteLink(w http.ResponseWriter, req *http.Request) {

}

func (serv *Server) handleGetLinksByTags(w http.ResponseWriter, req *http.Request) {

}

func BuildServer(serverUrl string, chatStore *storage.ChatStore, repoStore *storage.RepoStore, chatRepoRecordStore *storage.ChatRepoRecordStore) *Server {
	server := &Server{serverUrl: serverUrl, chatStore: chatStore, repoStore: repoStore, chatRepoRecordStore: chatRepoRecordStore}
	server.router = mux2.NewRouter()
	server.router.HandleFunc("/user", server.handleAddUser).Methods("POST")
	server.router.HandleFunc("/user", server.handleDeleteUser).Methods("DELETE")
	server.router.HandleFunc("/repo", server.handleGetLinks).Methods("GET")
	server.router.HandleFunc("/repo/tags", server.handleGetLinksByTags).Methods("GET")
	server.router.HandleFunc("/repo", server.handleAddLink).Methods("POST")
	server.router.HandleFunc("/repo", server.handleDeleteLink).Methods("DELETE")
	return server
}

func (serv *Server) Start() error {
	err := http.ListenAndServe(serv.serverUrl, serv.router)
	return err
}

func (serv *Server) Stop() {

}
