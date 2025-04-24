package state_machine

import (
	"TeleRequestHandler/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetTransitions(commandName string, bot bot.Bot[tgbotapi.UpdatesChannel, tgbotapi.MessageConfig]) map[StateName]State {
	ans := make(map[StateName]State)
	switch commandName {
	case "add":
		ans[NONE] = NewState(WAIT_LINK, NewWaitLinkState(bot))
		ans[WAIT_LINK] = NewState(WAIT_EVENTS, NewWaitEventsState(bot))
		ans[WAIT_EVENTS] = NewState(WAIT_TAGS, NewWaitTagsState(bot))
		ans[WAIT_TAGS] = NewState(NONE, NewNoneState(bot))
	case "del":
		ans[NONE] = NewState(WAIT_LINK, NewWaitLinkState(bot))
		ans[WAIT_LINK] = NewState(NONE, NewNoneState(bot))
	case "repos":
		ans[NONE] = NewState(WAIT_TAGS, NewWaitTagsState(bot))
		ans[WAIT_TAGS] = NewState(NONE, NewNoneState(bot))
	}
	return ans
}
