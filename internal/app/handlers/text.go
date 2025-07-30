package handlers

import (
	tele "gopkg.in/telebot.v3"
)

func RegisterTextHandlers(bot *tele.Bot) {
	bot.Handle(tele.OnText, func(c tele.Context) error {

		switch userState[c.Sender().ID].GetState() {
		case "awaiting_pack_name":
			return handlePackName(c)
		case "add_desc":
			return handlePackDesc(c)
		case "add_round":
			return handleNewRound(c)
		case "add_theme":
			return handleNewTheme(c)
		case "awaiting_price":
			return handleNewCost(c)
		case "awaiting_quest_disc":
			return handleNewQuestDisc(c)
		case "awaiting_text":
			return handleNewText(c)
		case "awaiting_answer":
			return handleNewAnswer(c)
		default:
			return nil
		}
	})

}
