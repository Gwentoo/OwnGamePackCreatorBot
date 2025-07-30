package handlers

import (
	"OwnGamePack/internal/app/validators"
	"fmt"
	tele "gopkg.in/telebot.v3"
)

func RegisterVideoHandlers(bot *tele.Bot) {
	bot.Handle(tele.OnVideo, func(c tele.Context) error {

		switch userState[c.Sender().ID].GetState() {
		case "awaiting_video":
			return handleNewVideo(c)
		default:
			return nil
		}
	})

}

func handleNewVideo(c tele.Context) error {

	video := c.Message().Video
	err := validators.IsValidVideo(video)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ %s\nПопробуйте еще раз:", err))
	}

	pos0 := userState[c.Sender().ID].GetPos(0)
	pos1 := userState[c.Sender().ID].GetPos(1)
	pos2 := userState[c.Sender().ID].GetPos(2)
	TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Video = video

	userState[c.Sender().ID].SetState("awaiting_answer")

	return c.Send("✅ Видео добавлено в вопрос.\n\n✏️ Укажите ответ на вопрос")

}
