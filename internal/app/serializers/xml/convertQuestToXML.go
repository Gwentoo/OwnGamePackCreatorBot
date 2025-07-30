package xml

import (
	"OwnGamePack/internal/app/utils"
	"OwnGamePack/internal/structs"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"strconv"
)

var (
	AudioCount int
	VideoCount int
	PhotoCount int
)

func convertQuestToXML(bot *tele.Bot, quest structs.Quest, pack structs.Pack) XMLQuestion {
	dirName := "C:\\Users\\timof\\Desktop\\Packs\\" + strconv.FormatInt(pack.PackID, 10)
	xmlQuestion := XMLQuestion{
		Price: quest.Cost,
		Type:  "",
		Right: struct {
			Answer string `xml:"answer"`
		}{
			Answer: quest.Answer,
		},
	}

	questionParam := XMLParam{
		Name: "question",
		Type: "content",
	}

	var item XMLItem
	switch quest.ContentType {
	case "audio":
		AudioCount += 1
		utils.DownloadFile(bot, quest.Audio.FileID, dirName+fmt.Sprintf("\\Audio\\%d.mp3", AudioCount))
		item = XMLItem{Type: "audio", IsRef: true, Content: fmt.Sprintf("%d.mp3", AudioCount)}
	case "video":
		VideoCount += 1
		utils.DownloadFile(bot, quest.Video.FileID, dirName+fmt.Sprintf("\\Video\\%d.mp4", VideoCount))
		item = XMLItem{Type: "video", IsRef: true, Content: fmt.Sprintf("%d.mp4", VideoCount)}
	case "photo":
		PhotoCount += 1
		utils.DownloadFile(bot, quest.Photo.FileID, dirName+fmt.Sprintf("\\Images\\%d.png", PhotoCount))
		item = XMLItem{Type: "image", IsRef: true, Content: fmt.Sprintf("%d.png", PhotoCount)}
	default:
		item = XMLItem{Type: "text", Content: quest.Text}
	}
	questionParam.Item = &item

	xmlQuestion.Params.Params = append(xmlQuestion.Params.Params, questionParam)

	if quest.Type == "secret" {
		xmlQuestion.Type = "secret"

		selectionModeParam := XMLParam{
			Name:    "selectionMode",
			Content: "exceptCurrent",
		}

		priceParam := XMLParam{
			Name:      "price",
			Type:      "numberSet",
			NumberSet: &XMLNumberSet{Minimum: "0", Maximum: "0", Step: "0"},
		}

		themeParam := XMLParam{
			Name: "theme",
		}

		xmlQuestion.Params.Params = append(xmlQuestion.Params.Params,
			selectionModeParam,
			priceParam,
			themeParam,
		)
	}

	if quest.Type == "bet" {
		xmlQuestion.Type = "stake"
		xmlQuestion.Price = "5000"
	}

	return xmlQuestion
}
