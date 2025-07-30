package structs

import tele "gopkg.in/telebot.v3"

type Quest struct {
	Description string      `json:"questDisc"`
	Cost        string      `json:"cost"`
	Type        string      `json:"type"`
	Text        string      `json:"text"`
	Answer      string      `json:"answer"`
	ContentType string      `json:"contentType"`
	Audio       *tele.Audio `json:"Audio"`
	Video       *tele.Video `json:"Video"`
	Photo       *tele.Photo `json:"Photo"`
}

func NewQuest(Type string) Quest {
	return Quest{Type: Type}
}
