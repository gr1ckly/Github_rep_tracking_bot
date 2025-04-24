package converters

import (
	"Common"
	"TeleRequestHandler/controller/state_machine"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func ConvertChat(usrCtx state_machine.UserContext, update tgbotapi.Update) Common.ChatDTO {
	return Common.ChatDTO{
		ChatID: usrCtx.ChatId,
		Type:   update.Message.Chat.Type,
	}
}

func ConvertRepo(usrCtx state_machine.UserContext) Common.RepoDTO {
	return Common.RepoDTO{
		Link:   usrCtx.Link,
		Tags:   usrCtx.Tags,
		Events: usrCtx.Events,
	}
}

func ConvertToMessage(repos []Common.RepoDTO) string {
	builder := strings.Builder{}
	if len(repos) == 0 {
		builder.WriteString("У вас пока нет отслеживаемых репозиториев")
	}
	for _, repo := range repos {
		builder.WriteString("Ссылка: ")
		builder.WriteString(repo.Link)
		builder.WriteRune('\n')
		builder.WriteString("События: ")
		for _, event := range repo.Events {
			builder.WriteString(event)
			builder.WriteString(" ")
		}
		builder.WriteRune('\n')
		builder.WriteString("Теги: ")
		for _, tag := range repo.Tags {
			builder.WriteRune('#')
			builder.WriteString(tag)
			builder.WriteString(" ")
		}
		builder.WriteRune('\n')
		builder.WriteRune('\n')
	}
	return builder.String()
}
