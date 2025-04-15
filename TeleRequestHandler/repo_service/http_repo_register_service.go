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
	client     *http.Client
	serverUrl  string
	repoApiExt string
}

func NewHttpRepoRegisterService(serverUrl string, chatApiExt string, timeout int) *HttpRepoRegisterService {
	return &HttpRepoRegisterService{client: &http.Client{Timeout: time.Duration(timeout) * time.Second}, serverUrl: serverUrl, repoApiExt: chatApiExt}
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
	var result []Common.RepoDTO
	items, ok := res.Result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("ожидался []interface{}")
	}

	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("ожидался map[string]interface{}")
		}

		dto := Common.RepoDTO{
			Link:   m["link"].(string),
			Tags:   rs.convertToStringSlice(m["tags"]),
			Events: rs.convertToStringSlice(m["events"]),
		}
		result = append(result, dto)
	}

	return result, nil
}

func (rs *HttpRepoRegisterService) convertToStringSlice(v interface{}) []string {
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	var result []string
	for _, item := range raw {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func (rs *HttpRepoRegisterService) AddRepo(chatId int64, dto Common.RepoDTO) error {
	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	_, err = rs.fetch("POST", fmt.Sprintf("%v/%v/%v", rs.serverUrl, rs.repoApiExt, chatId), nil, data)
	return err
}

func (rs *HttpRepoRegisterService) GetReposByChat(chatId int64) ([]Common.RepoDTO, error) {
	data, err := rs.fetch("GET", fmt.Sprintf("%v/%v/%v", rs.serverUrl, rs.repoApiExt, chatId), nil, nil)
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
	data, err := rs.fetch("GET", fmt.Sprintf("%v/%v/%v/%v", rs.serverUrl, rs.repoApiExt, chatId, tag), nil, nil)
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
	_, err := rs.fetch("DELETE", fmt.Sprintf("%v/%v/%v", rs.serverUrl, rs.repoApiExt, chatId), params, nil)
	return err
}
