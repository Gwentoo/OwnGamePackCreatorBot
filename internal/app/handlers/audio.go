package handlers

import (
	"OwnGamePack/internal/app/validators"
	"fmt"
	tele "gopkg.in/telebot.v3"
)

func RegisterAudioHandlers(bot *tele.Bot) {
	bot.Handle(tele.OnAudio, func(c tele.Context) error {

		switch userState[c.Sender().ID].GetState() {
		case "awaiting_audio":
			return handleNewAudio(c)
		default:
			return nil
		}
	})

}

func handleNewAudio(c tele.Context) error {
	audio := c.Message().Audio
	err := validators.IsValidAudio(audio)
	if err != nil {
		return c.Send(fmt.Sprintf("❌ %s\nПопробуйте еще раз:", err))
	}

	pos0 := userState[c.Sender().ID].GetPos(0)
	pos1 := userState[c.Sender().ID].GetPos(1)
	pos2 := userState[c.Sender().ID].GetPos(2)
	TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Audio = audio

	userState[c.Sender().ID].SetState("awaiting_answer")

	return c.Send("✅ Аудио добавлено в вопрос.\n\n✏️ Укажите ответ на вопрос")

}
