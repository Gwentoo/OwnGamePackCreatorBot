package handlers

import (
	"OwnGamePack/internal/app/validators"
	"fmt"
	tele "gopkg.in/telebot.v3"
)

func RegisterPhotoHandlers(bot *tele.Bot) {
	bot.Handle(tele.OnPhoto, func(c tele.Context) error {

		switch userState[c.Sender().ID].GetState() {
		case "awaiting_photo":
			return handleNewPhoto(c)
		default:
			return nil
		}
	})

}

func handleNewPhoto(c tele.Context) error {
	photo := c.Message().Photo
	err := validators.IsValidPhoto(photo)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ %s\nПопробуйте еще раз:", err))
	}

	pos0 := userState[c.Sender().ID].GetPos(0)
	pos1 := userState[c.Sender().ID].GetPos(1)
	pos2 := userState[c.Sender().ID].GetPos(2)
	TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Photo = photo

	userState[c.Sender().ID].SetState("awaiting_answer")

	return c.Send("✅ Фото добавлено в вопрос.\n\n✏️ Укажите ответ на вопрос")

}
