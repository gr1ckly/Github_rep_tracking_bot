package server

import (
	dtos2 "Crypto_Bot/MainServer/server/dtos"
	"Crypto_Bot/MainServer/server/validators"
	"encoding/json"
	mux2 "github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

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
	chats, err := serv.storeManager.GetChats()
	if err != nil {
		errDto := dtos2.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
	}
	dtos := make([]*dtos2.ChatDTO, len(chats))
	for idx, chat := range chats {
		conv, err := dtos2.ConvertChatDTO(&chat)
		if err != nil {
			errDto := dtos2.ErrorDTO{err.Error()}
			logger.Error(err.Error())
			serv.sendAns(errDto, 500, w)
		}
		dtos[idx] = conv
	}
	result := dtos2.ResultDTO[[]*dtos2.ChatDTO]{dtos}
	serv.sendAns(result, 200, w)
}

func (serv *Server) handleAddChat(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var chatDto dtos2.ChatDTO
	err := json.NewDecoder(req.Body).Decode(&chatDto)
	if err != nil {
		errDto := dtos2.ErrorDTO{"Invalid JSON"}
		logger.Error("Invalid chat_id format " + err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	newChat, err := dtos2.ParseChat(chatDto)
	if err != nil {
		errDto := dtos2.ErrorDTO{err.Error()}
		logger.Error("Error of parsing chat " + err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	id, err := serv.storeManager.AddChat(newChat)
	if err != nil {
		errDto := dtos2.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	resultDto := dtos2.ResultDTO[int]{id}
	serv.sendAns(resultDto, 200, w)
}

func (serv *Server) handleDeleteChat(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chatId, err := strconv.Atoi(strings.TrimSpace(req.URL.Query().Get("chat_id")))
	if err != nil {
		errDto := dtos2.ErrorDTO{"Invalid chat_id format"}
		logger.Error("Invalid chat_id format " + err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	err = serv.storeManager.DeleteChat(chatId)
	if err != nil {
		errDto := dtos2.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	serv.sendAns(nil, 200, w)
}

func (serv *Server) handleGetRepos(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chatId, err := strconv.Atoi(strings.TrimSpace(req.URL.Query().Get("chat_id")))
	if err != nil {
		errDto := dtos2.ErrorDTO{"Invalid chat_id format"}
		logger.Error("Invalid chat_id format " + err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	records, err := serv.storeManager.GetReposByChat(chatId)
	if err != nil {
		errDto := dtos2.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	ans := make([]*dtos2.RepoDTO, len(records))
	for idx, record := range records {
		conv, err := dtos2.ConvertRepoDTO(&record)
		if err != nil {
			errDto := dtos2.ErrorDTO{err.Error()}
			logger.Error(err.Error())
			serv.sendAns(errDto, 500, w)
		}
		ans[idx] = conv
	}
	resultDto := dtos2.ResultDTO[[]*dtos2.RepoDTO]{ans}
	serv.sendAns(resultDto, 200, w)
}

func (serv *Server) handleAddRepo(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var repoDto dtos2.RepoDTO
	err := json.NewDecoder(req.Body).Decode(&repoDto)
	if err != nil {
		errDto := dtos2.ErrorDTO{"Invalid JSON"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	if !serv.urlValidator.Check(repoDto.Link) {
		errDto := dtos2.ErrorDTO{"Invalid link"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	repo, err := dtos2.ParseRepo(repoDto)
	if err != nil {
		errDto := dtos2.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	id, err := serv.storeManager.AddRepo(repo, &repoDto)
	if err != nil {
		errDto := dtos2.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	ans := dtos2.ResultDTO[int]{id}
	serv.sendAns(ans, 200, w)
}

func (serv *Server) handleDeleteRepo(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chatID, err := strconv.Atoi(req.URL.Query().Get("chat_id"))
	if err != nil {
		errDto := dtos2.ErrorDTO{"Invalid chat_id"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	rawLink := req.URL.Query().Get("link")
	if rawLink == "" {
		errDto := dtos2.ErrorDTO{"Invalid link"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	name, owner, err := dtos2.ParseNameAndOwner(rawLink)
	if err != nil {
		errDto := dtos2.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	err = serv.storeManager.DeleteRepo(chatID, owner, name)
	if err != nil {
		errDto := dtos2.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	serv.sendAns(nil, 200, w)
}

func (serv *Server) handleGetReposByTag(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chat_id, err := strconv.Atoi(req.URL.Query().Get("chat_id"))
	if err != nil {
		errDto := dtos2.ErrorDTO{"Invalid chat_id"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	tag := req.URL.Query().Get("tag")
	if tag == "" {
		errDto := dtos2.ErrorDTO{"Invalid tag"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	records, err := serv.storeManager.GetReposByTag(chat_id, tag)
	if err != nil {
		errDto := dtos2.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	ans := make([]*dtos2.RepoDTO, len(records))
	for idx, record := range records {
		conv, err := dtos2.ConvertRepoDTO(record)
		if err != nil {
			errDto := dtos2.ErrorDTO{err.Error()}
			logger.Error(err.Error())
			serv.sendAns(errDto, 500, w)
			return
		}
		ans[idx] = conv
	}
	result := dtos2.ResultDTO[[]*dtos2.RepoDTO]{ans}
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
