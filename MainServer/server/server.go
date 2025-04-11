package server

import (
	"Common"
	dtos2 "Crypto_Bot/MainServer/server/dtos"
	"Crypto_Bot/MainServer/server/validators"
	"context"
	"encoding/json"
	"errors"
	mux2 "github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

type Server struct {
	urlValidator *validators.UrlValidator
	storeManager *StoreManager
	server       *http.Server
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
		errDto := Common.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
	}
	dtos := make([]*Common.ChatDTO, len(chats))
	for idx, chat := range chats {
		conv, err := dtos2.ConvertChatDTO(&chat)
		if err != nil {
			errDto := Common.ErrorDTO{err.Error()}
			logger.Error(err.Error())
			serv.sendAns(errDto, 500, w)
		}
		dtos[idx] = conv
	}
	result := Common.ResultDTO[[]*Common.ChatDTO]{dtos}
	serv.sendAns(result, 200, w)
}

func (serv *Server) handleAddChat(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var chatDto Common.ChatDTO
	err := json.NewDecoder(req.Body).Decode(&chatDto)
	if err != nil {
		errDto := Common.ErrorDTO{"Invalid JSON"}
		logger.Error("Invalid chat_id format " + err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	newChat, err := dtos2.ParseChat(chatDto)
	if err != nil {
		errDto := Common.ErrorDTO{err.Error()}
		logger.Error("Error of parsing chat " + err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	id, err := serv.storeManager.AddChat(newChat)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			errDto := Common.ErrorDTO{"Chat already register"}
			logger.Error(err.Error())
			serv.sendAns(errDto, 409, w)
			return
		}
	}
	if err != nil {
		errDto := Common.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	resultDto := Common.ResultDTO[int]{id}
	serv.sendAns(resultDto, 200, w)
}

func (serv *Server) handleDeleteChat(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chatId, err := strconv.ParseInt(strings.TrimSpace(req.URL.Query().Get("chat_id")), 10, 64)
	if err != nil {
		errDto := Common.ErrorDTO{"Invalid chat_id format"}
		logger.Error("Invalid chat_id format " + err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	err = serv.storeManager.DeleteChat(chatId)
	if err != nil {
		errDto := Common.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	serv.sendAns(nil, 200, w)
}

func (serv *Server) handleGetRepos(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chatId, err := strconv.ParseInt(strings.TrimSpace(req.URL.Query().Get("chat_id")), 10, 64)
	if err != nil {
		errDto := Common.ErrorDTO{"Invalid chat_id format"}
		logger.Error("Invalid chat_id format " + err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	records, err := serv.storeManager.GetReposByChat(chatId)
	if err != nil {
		errDto := Common.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	ans := make([]*Common.RepoDTO, len(records))
	for idx, record := range records {
		conv, err := dtos2.ConvertRepoDTO(&record)
		if err != nil {
			errDto := Common.ErrorDTO{err.Error()}
			logger.Error(err.Error())
			serv.sendAns(errDto, 500, w)
		}
		ans[idx] = conv
	}
	resultDto := Common.ResultDTO[[]*Common.RepoDTO]{ans}
	serv.sendAns(resultDto, 200, w)
}

func (serv *Server) handleAddRepo(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var repoDto Common.RepoDTO
	err := json.NewDecoder(req.Body).Decode(&repoDto)
	if err != nil {
		errDto := Common.ErrorDTO{"Invalid JSON"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	if !serv.urlValidator.Check(repoDto.Link) {
		errDto := Common.ErrorDTO{"Invalid link"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	repo, err := dtos2.ParseRepo(repoDto)
	if err != nil {
		errDto := Common.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	vars := mux2.Vars(req)
	chatId, err := strconv.ParseInt(vars["chat_id"], 10, 64)
	if err != nil {
		errDto := Common.ErrorDTO{"Invalid chat_id"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	id, err := serv.storeManager.AddRepo(repo, &repoDto, chatId)
	if err != nil {
		errDto := Common.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	ans := Common.ResultDTO[int]{id}
	serv.sendAns(ans, 200, w)
}

func (serv *Server) handleDeleteRepo(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	vars := mux2.Vars(req)
	chatId, err := strconv.ParseInt(vars["chat_id"], 10, 64)
	if err != nil {
		errDto := Common.ErrorDTO{"Invalid chat_id"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	rawLink := req.URL.Query().Get("link")
	if rawLink == "" {
		errDto := Common.ErrorDTO{"Invalid link"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	name, owner, err := dtos2.ParseNameAndOwner(rawLink)
	if err != nil {
		errDto := Common.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	num, err := serv.storeManager.DeleteRepo(chatId, owner, name)
	if err != nil {
		errDto := Common.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	if num == 0 {
		serv.sendAns(nil, 404, w)
		return
	}
	serv.sendAns(nil, 200, w)
}

func (serv *Server) handleGetReposByTag(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	vars := mux2.Vars(req)
	chatId, err := strconv.ParseInt(vars["chat_id"], 10, 64)
	if err != nil {
		errDto := Common.ErrorDTO{"Invalid chat_id"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	tag := vars["tag"]
	if tag == "" {
		errDto := Common.ErrorDTO{"Invalid tag"}
		logger.Error(err.Error())
		serv.sendAns(errDto, 400, w)
		return
	}
	records, err := serv.storeManager.GetReposByTag(chatId, tag)
	if err != nil {
		errDto := Common.ErrorDTO{err.Error()}
		logger.Error(err.Error())
		serv.sendAns(errDto, 500, w)
		return
	}
	ans := make([]*Common.RepoDTO, len(records))
	for idx, record := range records {
		conv, err := dtos2.ConvertRepoDTO(record)
		if err != nil {
			errDto := Common.ErrorDTO{err.Error()}
			logger.Error(err.Error())
			serv.sendAns(errDto, 500, w)
			return
		}
		ans[idx] = conv
	}
	result := Common.ResultDTO[[]*Common.RepoDTO]{ans}
	serv.sendAns(result, 200, w)
}

func BuildServer(serverUrl string, urlValidator *validators.UrlValidator, storeManager *StoreManager) *Server {
	server := &Server{urlValidator: urlValidator, storeManager: storeManager}
	router := mux2.NewRouter()
	router.HandleFunc("/chat", server.handleAddChat).Methods("POST")
	router.HandleFunc("/chat", server.handleDeleteChat).Methods("DELETE")
	router.HandleFunc("/chat", server.handleGetChats).Methods("GET")
	router.HandleFunc("/repo/{chat_id}", server.handleGetRepos).Methods("GET")
	router.HandleFunc("/repo/{chat_id}/{tag}", server.handleGetReposByTag).Methods("GET")
	router.HandleFunc("/repo/{chat_id}", server.handleAddRepo).Methods("POST")
	router.HandleFunc("/repo/{chat_id}", server.handleDeleteRepo).Methods("DELETE")
	svr := &http.Server{
		Addr:    serverUrl,
		Handler: router,
	}
	server.server = svr
	return server
}

func (serv *Server) Start() error {
	return serv.server.ListenAndServe()
}

func (serv *Server) Stop(ctx context.Context) error {
	return serv.server.Shutdown(ctx)
}
