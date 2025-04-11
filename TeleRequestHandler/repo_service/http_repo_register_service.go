package repo_service

import (
	"Common"
	"TeleRequestHandler/custom_erros"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpRepoRegisterService struct {
	client         *http.Client
	serverUrl      string
	repoApiExt     string
	repoTagsApiExt string
}

func NewHttpChatRegisterService(serverUrl string, chatApiExt string, repoTagsApiExt string, timeout int) *HttpRepoRegisterService {
	return &HttpRepoRegisterService{client: &http.Client{Timeout: time.Duration(timeout) * time.Second}, serverUrl: serverUrl, repoApiExt: chatApiExt, repoTagsApiExt: repoTagsApiExt}
}

func (rs *HttpRepoRegisterService) fetch(method string, baseUrl string, params url.Values, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, baseUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = params.Encode()
	resp, err := rs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		data, _ := io.ReadAll(resp.Body)
		var errDto Common.ErrorDTO
		_ = json.Unmarshal(data, &errDto)
		return nil, custom_erros.ServerError{Status: resp.Status, StatusCode: resp.StatusCode, Url: resp.Request.RequestURI, ErrDTO: errDto}
	}
	return io.ReadAll(resp.Body)
}

func (rs *HttpRepoRegisterService) tryType(res Common.ResultDTO[[]Common.RepoDTO]) ([]Common.RepoDTO, error) {
	result, ok := res.Result.([]Common.RepoDTO)
	if !ok {
		return nil, fmt.Errorf("Waiting []Common.RepoDTO")
	}
	return result, nil
}

func (rs *HttpRepoRegisterService) AddRepo(chatId int64, dto Common.RepoDTO) error {
	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	_, err = rs.fetch("POST", fmt.Sprintf(rs.serverUrl+rs.repoApiExt+"/%v", chatId), nil, data)
	return err
}

func (rs *HttpRepoRegisterService) GetReposByChat(chatId int64) ([]Common.RepoDTO, error) {
	data, err := rs.fetch("GET", fmt.Sprintf(rs.serverUrl+rs.repoApiExt+"/%v", chatId), nil, nil)
	if err != nil {
		return nil, err
	}
	var res Common.ResultDTO[[]Common.RepoDTO]
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return rs.tryType(res)
}

func (rs *HttpRepoRegisterService) GetReposByTag(chatId int64, tag string) ([]Common.RepoDTO, error) {
	data, err := rs.fetch("GET", fmt.Sprintf(rs.serverUrl+rs.repoApiExt+"/%v/%v", chatId, tag), nil, nil)
	if err != nil {
		return nil, err
	}
	var res Common.ResultDTO[[]Common.RepoDTO]
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return rs.tryType(res)
}

func (rs *HttpRepoRegisterService) DeleteRepo(chatId int64, link string) error {
	params := url.Values{}
	params.Add("link", link)
	_, err := rs.fetch("DELETE", fmt.Sprintf(rs.serverUrl+rs.repoApiExt+"/%v", chatId), params, nil)
	return err
}
