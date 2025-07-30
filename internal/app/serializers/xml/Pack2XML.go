package xml

import (
	"OwnGamePack/internal/app/utils"
	"OwnGamePack/internal/structs"
	"encoding/xml"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"os"
	"strconv"
	"time"
)

func ConvertPackToXML(bot *tele.Bot, pack structs.Pack) ([]byte, error) {
	xmlPack := XMLPackage{
		ID:      fmt.Sprintf("%d", pack.PackID),
		Name:    pack.PackName,
		Version: "5",
		Date:    time.Now().Format("02.01.2006"),
		Tags: struct {
			Tag []string `xml:"tag"`
		}{
			Tag: pack.PackTags,
		},
		Info: struct {
			Authors struct {
				Author []string `xml:"author"`
			} `xml:"authors"`
		}{
			Authors: struct {
				Author []string `xml:"author"`
			}{
				Author: []string{pack.UserName},
			},
		},
	}

	dirName := "C:\\Users\\timof\\Desktop\\Packs\\" + strconv.FormatInt(pack.PackID, 10)
	err4 := utils.CheckAndRemoveDir(dirName)
	if err4 != nil {
		return nil, err4
	}
	err := os.Mkdir(dirName, 0750)
	if err != nil {
		return nil, err
	}
	err1 := os.Mkdir(dirName+"\\Images", 0750)
	if err1 != nil {
		return nil, err1
	}
	err2 := os.Mkdir(dirName+"\\Audio", 0750)
	if err2 != nil {
		return nil, err2
	}
	err3 := os.Mkdir(dirName+"\\Video", 0750)
	if err3 != nil {
		return nil, err3
	}

	for _, round := range pack.Rounds {
		xmlRound := XMLRound{
			Name: round.Name,
		}

		for _, theme := range round.Themes {
			xmlTheme := XMLTheme{
				Name: theme.Name,
			}

			for _, quest := range theme.Quests {
				xmlQuestion := convertQuestToXML(bot, quest, pack)
				xmlTheme.Questions = append(xmlTheme.Questions, xmlQuestion)
			}

			xmlRound.Themes = append(xmlRound.Themes, xmlTheme)
		}

		xmlPack.Rounds = append(xmlPack.Rounds, xmlRound)
	}

	output, err := xml.MarshalIndent(xmlPack, "", "  ")
	if err != nil {
		return nil, err
	}

	result := []byte(xml.Header + string(output))
	AudioCount = 0
	VideoCount = 0
	PhotoCount = 0
	return result, nil
}
