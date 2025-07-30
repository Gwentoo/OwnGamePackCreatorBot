package utils

import (
	"OwnGamePack/internal/structs"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"slices"
)

var QuizTags = []string{
	"Аниме",
	"Манга",
	"Видеоигры",
	"Флаги",
	"Мемы",
	"История",
	"Наука",
	"Кино",
	"Музыка",
	"Литература",
	"География",
	"Искусство",
	"Спорт",
	"Технологии",
	"Математика",
	"Политика",
	"Экономика",
	"Медицина",
	"Космос",
	"Языки",
	"Архитектура",
	"Автомобили",
	"Изобретения",
	"Праздники",
}

func BuildTagMenu(pack *structs.Pack) [][]tele.InlineButton {
	buttons := make([]tele.InlineButton, 24)
	for i, tag := range QuizTags {
		if !slices.Contains(pack.PackTags, tag) {
			btn := tele.InlineButton{
				Data: "ST_" + tag + fmt.Sprintf("_%d%d", i/3, i%3),
				Text: tag,
			}
			buttons[i] = btn
		} else {
			btn := tele.InlineButton{
				Data: "ST_S_" + tag + fmt.Sprintf("_%d%d", i/3, i%3),
				Text: "✅ " + tag,
			}
			buttons[i] = btn
		}
	}
	grid := make([][]tele.InlineButton, 8)
	for i := 0; i < 8; i++ {
		grid[i] = buttons[i*3 : (i+1)*3]
	}
	return grid
}
