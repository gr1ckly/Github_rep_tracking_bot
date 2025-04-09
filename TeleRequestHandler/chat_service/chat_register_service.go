package chat_service

import (
	"Common"
)

type ChatRegisterService interface {
	RegisterChat(dto Common.ChatDTO) error
	DeleteChat(chatId int) error
}
