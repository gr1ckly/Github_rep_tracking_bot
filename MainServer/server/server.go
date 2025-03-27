package server

import (
	"Crypto_Bot/MainServer/server/validators"
	"context"
	"encoding/json"
	mux2 "github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	serverUrl    string
	router       *mux2.Router
	urlValidator *validators.UrlValidator
	storeManager *StoreManager
}

func (serv *Server) sendAns(msg any, statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	if msg != nil {
		body, _ := json.Marshal(msg)
		w.Write(body)
	}
}

func (serv *Server) handleGetChats(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chats, err := serv.storeManager.GetChats(context.Background())
	if err != nil {
		errDto := ErrorDTO{err.Error()}
		serv.sendAns(errDto, 500, w)
	}
	dtos := make([]*ChatDTO, len(chats))
	for idx, chat := range chats {
		conv, err := ConvertChatDTO(&chat)
		if err != nil {
			errDto := ErrorDTO{err.Error()}
			serv.sendAns(errDto, 500, w)
		}
		dtos[idx] = conv
	}
	result := ResultDTO[[]*ChatDTO]{dtos}
	serv.sendAns(result, 200, w)
}

func (serv *Server) handleAddChat(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var chatDto ChatDTO
	err := json.NewDecoder(req.Body).Decode(&chatDto)
	if err != nil {
		errDto := ErrorDTO{"Invalid JSON"}
		serv.sendAns(errDto, 400, w)
		return
	}
	newChat, err := ParseChat(chatDto)
	if err != nil {
		errDto := ErrorDTO{err.Error()}
		serv.sendAns(errDto, 400, w)
		return
	}
	id, err := serv.storeManager.AddChat(context.Background(), newChat)
	if err != nil {
		errDto := ErrorDTO{err.Error()}
		serv.sendAns(errDto, 500, w)
		return
	}
	resultDto := ResultDTO[int]{id}
	serv.sendAns(resultDto, 200, w)
}

func (serv *Server) handleDeleteChat(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chatId, err := strconv.Atoi(strings.TrimSpace(req.URL.Query().Get("chat_id")))
	if err != nil {
		errDto := ErrorDTO{"Invalid chat_id format"}
		serv.sendAns(errDto, 400, w)
		return
	}
	err = serv.storeManager.DeleteChat(context.Background(), chatId)
	if err != nil {
		errDto := ErrorDTO{err.Error()}
		serv.sendAns(errDto, 500, w)
		return
	}
	serv.sendAns(nil, 200, w)
}

func (serv *Server) handleGetRepos(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chatId, err := strconv.Atoi(strings.TrimSpace(req.URL.Query().Get("chat_id")))
	if err != nil {
		errDto := ErrorDTO{"Invalid chat_id format"}
		serv.sendAns(errDto, 400, w)
		return
	}
	records, err := serv.storeManager.GetReposByChat(context.Background(), chatId)
	if err != nil {
		errDto := ErrorDTO{err.Error()}
		serv.sendAns(errDto, 500, w)
		return
	}
	ans := make([]*RepoDTO, len(records))
	for idx, record := range records {
		conv, err := ConvertRepoDTO(&record)
		if err != nil {
			errDto := ErrorDTO{err.Error()}
			serv.sendAns(errDto, 500, w)
		}
		ans[idx] = conv
	}
	resultDto := ResultDTO[[]*RepoDTO]{ans}
	serv.sendAns(resultDto, 200, w)
}

func (serv *Server) handleAddRepo(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var repoDto RepoDTO
	err := json.NewDecoder(req.Body).Decode(&repoDto)
	if err != nil {
		errDto := ErrorDTO{"Invalid JSON"}
		serv.sendAns(errDto, 400, w)
		return
	}
	if !serv.urlValidator.Check(repoDto.Link) {
		errDto := ErrorDTO{"Invalid link"}
		serv.sendAns(errDto, 400, w)
		return
	}
	repo, err := ParseRepo(repoDto)
	if err != nil {
		errDto := ErrorDTO{err.Error()}
		serv.sendAns(errDto, 400, w)
		return
	}
	id, err := serv.storeManager.AddRepo(context.Background(), repo, &repoDto)
	if err != nil {
		errDto := ErrorDTO{err.Error()}
		serv.sendAns(errDto, 500, w)
		return
	}
	ans := ResultDTO[int]{id}
	serv.sendAns(ans, 200, w)
}

func (serv *Server) handleDeleteRepo(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chatID, err := strconv.Atoi(req.URL.Query().Get("chat_id"))
	if err != nil {
		errDto := ErrorDTO{"Invalid chat_id"}
		serv.sendAns(errDto, 400, w)
		return
	}
	rawLink := req.URL.Query().Get("link")
	if rawLink == "" {
		errDto := ErrorDTO{"Invalid link"}
		serv.sendAns(errDto, 400, w)
		return
	}
	name, owner, err := ParseNameAndOwner(rawLink)
	if err != nil {
		errDto := ErrorDTO{err.Error()}
		serv.sendAns(errDto, 400, w)
		return
	}
	err = serv.storeManager.DeleteRepo(context.Background(), chatID, owner, name)
	if err != nil {
		errDto := ErrorDTO{err.Error()}
		serv.sendAns(errDto, 500, w)
		return
	}
	serv.sendAns(nil, 200, w)
}

func (serv *Server) handleGetReposByTag(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chat_id, err := strconv.Atoi(req.URL.Query().Get("chat_id"))
	if err != nil {
		errDto := ErrorDTO{"Invalid chat_id"}
		serv.sendAns(errDto, 400, w)
		return
	}
	tag := req.URL.Query().Get("tag")
	if tag == "" {
		errDto := ErrorDTO{"Invalid tag"}
		serv.sendAns(errDto, 400, w)
		return
	}
	records, err := serv.storeManager.GetReposByTag(context.Background(), chat_id, tag)
	if err != nil {
		errDto := ErrorDTO{err.Error()}
		serv.sendAns(errDto, 500, w)
		return
	}
	ans := make([]*RepoDTO, len(records))
	for idx, record := range records {
		conv, err := ConvertRepoDTO(record)
		if err != nil {
			errDto := ErrorDTO{err.Error()}
			serv.sendAns(errDto, 500, w)
			return
		}
		ans[idx] = conv
	}
	result := ResultDTO[[]*RepoDTO]{ans}
	serv.sendAns(result, 200, w)
}

func BuildServer(serverUrl string, urlValidator *validators.UrlValidator, storeManager *StoreManager) *Server {
	server := &Server{serverUrl: serverUrl, urlValidator: urlValidator, storeManager: storeManager}
	server.router = mux2.NewRouter()
	server.router.HandleFunc("/chat", server.handleAddChat).Methods("POST")
	server.router.HandleFunc("/chat", server.handleDeleteChat).Methods("DELETE")
	server.router.HandleFunc("/chat", server.handleGetChats).Methods("GET")
	server.router.HandleFunc("/repo", server.handleGetRepos).Methods("GET")
	server.router.HandleFunc("/repo/tags", server.handleGetReposByTag).Methods("GET")
	server.router.HandleFunc("/repo", server.handleAddRepo).Methods("POST")
	server.router.HandleFunc("/repo", server.handleDeleteRepo).Methods("DELETE")
	return server
}

func (serv *Server) Start() error {
	err := http.ListenAndServe(serv.serverUrl, serv.router)
	return err
}

func (serv *Server) Stop() {

}
