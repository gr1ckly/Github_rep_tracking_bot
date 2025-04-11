package state_machine

import "TeleRequestHandler/bot"

func GetTransitions(commandName string, bot bot.Bot[any, string, int64]) map[StateName]*State {
	ans := make(map[StateName]*State)
	switch commandName {

	}
	return ans
}
