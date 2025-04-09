package chat_service

import (
	"Common"
	"TeleRequestHandler/custom_erros"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type HttpChatRegisterService struct {
	client     *http.Client
	serverUrl  string
	chatApiExt string
}

func NewHttpChatRegisterService(serverUrl string, chatApiExt string, timeout int) *HttpChatRegisterService {
	return &HttpChatRegisterService{client: &http.Client{Timeout: time.Duration(timeout) * time.Second}, serverUrl: serverUrl, chatApiExt: chatApiExt}
}

func (cr *HttpChatRegisterService) fetch(method string, baseUrl string, params url.Values, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, baseUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = params.Encode()
	resp, err := cr.client.Do(req)
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

func (cr *HttpChatRegisterService) RegisterChat(dto Common.ChatDTO) error {
	data, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	_, err = cr.fetch("POST", cr.serverUrl+"/"+cr.chatApiExt, nil, data)
	return err
}

func (cr *HttpChatRegisterService) DeleteChat(chatId int) error {
	params := url.Values{}
	params.Add("chat_id", strconv.Itoa(chatId))
	_, err := cr.fetch("DELETE", cr.serverUrl+"/"+cr.chatApiExt, params, nil)
	return err
}
