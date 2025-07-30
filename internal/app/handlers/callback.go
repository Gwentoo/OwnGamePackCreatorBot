package handlers

import (
	"OwnGamePack/internal/app/googleDrive"
	"OwnGamePack/internal/app/serializers/json"
	"OwnGamePack/internal/app/serializers/xml"
	"OwnGamePack/internal/app/utils"
	"OwnGamePack/internal/storage"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	roundDel = -1
	themeDel = -1
)

func RegisterCallbackHandlers(bot *tele.Bot) {

	bot.Handle(tele.OnCallback, func(c tele.Context) error {

		if strings.HasPrefix(c.Callback().Data, "ST_") {
			parts := strings.Split(c.Callback().Data, "_")
			grid := utils.BuildTagMenu(TempPack[c.Sender().ID])
			if len(parts) == 4 {
				i, _ := strconv.Atoi(string(parts[3][0]))
				j, _ := strconv.Atoi(string(parts[3][1]))

				btn := tele.InlineButton{
					Data: "ST_" + parts[2] + fmt.Sprintf("_%d%d", i, j),
					Text: parts[2],
				}
				grid[i][j] = btn
				for k, v := range TempPack[c.Sender().ID].PackTags {
					if v == parts[2] {
						TempPack[c.Sender().ID].PackTags = append(TempPack[c.Sender().ID].PackTags[:k], TempPack[c.Sender().ID].PackTags[k+1:]...)
					}
				}

			} else {
				i, _ := strconv.Atoi(string(parts[2][0]))
				j, _ := strconv.Atoi(string(parts[2][1]))
				btn := tele.InlineButton{
					Data: "ST_S_" + parts[1] + fmt.Sprintf("_%d%d", i, j),
					Text: "✅ " + parts[1],
				}
				grid[i][j] = btn
				TempPack[c.Sender().ID].PackTags = append(TempPack[c.Sender().ID].PackTags, parts[1])

			}
			_, err := c.Bot().Edit(c.Message(), &tele.ReplyMarkup{InlineKeyboard: grid})
			return err
		}

		if strings.HasPrefix(c.Callback().Data, "select_theme_add_quest_") {
			selectedTheme := strings.TrimPrefix(c.Callback().Data, "select_theme_add_quest_")
			for i, theme := range TempPack[c.Sender().ID].Rounds[userState[c.Sender().ID].GetPos(0)].Themes {
				if theme.Name == selectedTheme {
					userState[c.Sender().ID].SetPos(1, i)
					break
				}
			}
			userState[c.Sender().ID].SetState("add_quest")
			menu.Inline(
				menu.Row(btnQuestDef, btnQuestBet),
				menu.Row(btnQuestSecret),
			)
			pos0 := userState[c.Sender().ID].GetPos(0)
			_, err := c.Bot().Edit(c.Message(), fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n✏️ Укажите тип вопроса",
				TempPack[c.Sender().ID].Rounds[pos0].Name, TempPack[c.Sender().ID].Rounds[pos0].Themes[userState[c.Sender().ID].GetPos(1)].Name), menu)
			return err
		}

		if strings.HasPrefix(c.Callback().Data, "select_round_add_quest_") {
			selectedRound := strings.TrimPrefix(c.Callback().Data, "select_round_add_quest_")
			for i, round := range TempPack[c.Sender().ID].Rounds {
				if round.Name == selectedRound {
					userState[c.Sender().ID].SetPos(0, i)
					break
				}
			}
			var keyboard [][]tele.InlineButton
			for _, theme := range TempPack[c.Sender().ID].Rounds[userState[c.Sender().ID].GetPos(0)].Themes {
				themeName := theme.Name
				btn := tele.InlineButton{
					Data: "select_theme_add_quest_" + themeName,
					Text: themeName,
				}
				keyboard = append(keyboard, []tele.InlineButton{btn})
			}
			storage.SaveMessage(c.Chat().ID, c.Message().Text, c.Message().ReplyMarkup)
			err1 := c.Bot().Delete(c.Message())
			if err1 != nil {
				return err1
			}
			return c.Send(fmt.Sprintf("Выберите тему из раунда '%s' для добавления вопроса", TempPack[c.Sender().ID].Rounds[userState[c.Sender().ID].GetPos(0)].Name), &tele.SendOptions{
				ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
			})
		}

		if strings.HasPrefix(c.Callback().Data, "select_round_add_theme_") {
			selectedRound := strings.TrimPrefix(c.Callback().Data, "select_round_add_theme_")
			for i, round := range TempPack[c.Sender().ID].Rounds {
				if round.Name == selectedRound {
					userState[c.Sender().ID].SetPos(0, i)
					break
				}
			}
			userState[c.Sender().ID].SetState("add_theme")
			_, err := c.Bot().Edit(c.Message(), fmt.Sprintf("📍 Раунд: %s\n\n✏️  Введите название темы", selectedRound), &tele.ReplyMarkup{})
			return err
		}

		if strings.HasPrefix(c.Callback().Data, "select_round_btn_del_") {
			var keyboard [][]tele.InlineButton
			selectedRound := strings.TrimPrefix(c.Callback().Data, "select_round_btn_del_")
			btnConfirm.Data = "round_btn_del_confirm_" + selectedRound
			keyboard = append(keyboard, []tele.InlineButton{btnBack})
			keyboard = append(keyboard, []tele.InlineButton{btnConfirm})
			err1 := c.Bot().Delete(c.Message())
			if err1 != nil {
				return err1
			}
			return c.Send(fmt.Sprintf("Вы точно хотите удалить раунд: %s?", selectedRound), &tele.SendOptions{
				ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
			})
		}

		if strings.HasPrefix(c.Callback().Data, "SR4DT_") {
			var keyboard [][]tele.InlineButton
			selectedRound := strings.TrimPrefix(c.Callback().Data, "SR4DT_")
			for i, round := range TempPack[c.Sender().ID].Rounds {
				if round.Name == selectedRound {
					roundDel = i
					for _, theme := range round.Themes {
						btn := tele.InlineButton{
							Data: "ST4D_" + theme.Name,
							Text: theme.Name,
						}
						keyboard = append(keyboard, []tele.InlineButton{btn})
					}
				}
			}
			err1 := c.Bot().Delete(c.Message())
			if err1 != nil {
				return err1
			}
			return c.Send(fmt.Sprintf("Выберите тему из раунда '%s' для удаления", selectedRound), &tele.SendOptions{
				ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
			})
		}

		if strings.HasPrefix(c.Callback().Data, "ST4D_") {
			var keyboard [][]tele.InlineButton
			selectedTheme := strings.TrimPrefix(c.Callback().Data, "ST4D_")
			btnConfirm.Data = "TDC_" + selectedTheme
			keyboard = append(keyboard, []tele.InlineButton{btnBack})
			keyboard = append(keyboard, []tele.InlineButton{btnConfirm})
			err1 := c.Bot().Delete(c.Message())
			if err1 != nil {
				return err1
			}
			return c.Send(fmt.Sprintf("Вы точно хотите удалить тему '%s' из раунда '%s'?", selectedTheme, TempPack[c.Sender().ID].Rounds[roundDel].Name), &tele.SendOptions{
				ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
			})
		}

		if strings.HasPrefix(c.Callback().Data, "SR4DQ_") {
			selectedRound := strings.TrimPrefix(c.Callback().Data, "SR4DQ_")
			var keyboard [][]tele.InlineButton
			for i, round := range TempPack[c.Sender().ID].Rounds {
				if round.Name == selectedRound {
					roundDel = i
					for _, theme := range round.Themes {
						if len(theme.Quests) != 0 {
							btn := tele.InlineButton{
								Data: "ST4DQ_" + theme.Name,
								Text: theme.Name,
							}
							keyboard = append(keyboard, []tele.InlineButton{btn})
						}
					}
				}
			}
			err1 := c.Bot().Delete(c.Message())
			if err1 != nil {
				return err1
			}
			return c.Send(fmt.Sprintf("Выберите тему из раунда '%s' для удаления вопроса", selectedRound), &tele.SendOptions{
				ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
			})
		}

		if strings.HasPrefix(c.Callback().Data, "ST4DQ_") {
			var keyboard [][]tele.InlineButton
			selectedRound := TempPack[c.Sender().ID].Rounds[roundDel].Name
			selectedTheme := strings.TrimPrefix(c.Callback().Data, "ST4DQ_")
			for j, theme := range TempPack[c.Sender().ID].Rounds[roundDel].Themes {
				if theme.Name == selectedTheme {
					themeDel = j
					for i, quest := range theme.Quests {
						btn := tele.InlineButton{
							Data: fmt.Sprintf("QD_%d", i),
							Text: quest.Description,
						}
						keyboard = append(keyboard, []tele.InlineButton{btn})
					}
				}
			}
			err1 := c.Bot().Delete(c.Message())
			if err1 != nil {
				return err1
			}
			return c.Send(fmt.Sprintf("Выберите вопрос из темы '%s' из раунда '%s' для удаления", selectedTheme, selectedRound), &tele.SendOptions{
				ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
			})
		}

		if strings.HasPrefix(c.Callback().Data, "QD_") {
			var keyboard [][]tele.InlineButton
			selectedRound := TempPack[c.Sender().ID].Rounds[roundDel].Name
			selectedTheme := TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel].Name
			questNum, _ := strconv.Atoi(strings.TrimPrefix(c.Callback().Data, "QD_"))
			selectedQuest := TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel].Quests[questNum]
			btnConfirm.Data = fmt.Sprintf("QDC_%d", questNum)
			keyboard = append(keyboard, []tele.InlineButton{btnBack})
			keyboard = append(keyboard, []tele.InlineButton{btnConfirm})
			err1 := c.Bot().Delete(c.Message())
			if err1 != nil {
				return err1
			}
			switch selectedQuest.ContentType {
			case "audio":
				return c.Send(selectedQuest.Audio, fmt.Sprintf("%s\nВы точно хотите удалить вопрос из темы '%s' из раунда '%s'?", selectedQuest.Audio.Caption, selectedTheme, selectedRound), &tele.SendOptions{
					ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
				})
			case "video":
				return c.Send(selectedQuest.Video, fmt.Sprintf("%s\nВы точно хотите удалить вопрос из темы '%s' из раунда '%s'?", selectedQuest.Video.Caption, selectedTheme, selectedRound), &tele.SendOptions{
					ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
				})
			case "photo":
				return c.Send(selectedQuest.Photo, fmt.Sprintf("%s\nВы точно хотите удалить вопрос из темы '%s' из раунда '%s'?", selectedQuest.Photo.Caption, selectedTheme, selectedRound), &tele.SendOptions{
					ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
				})
			default:
				return c.Send(fmt.Sprintf("\"📍 Раунд: %s\\n\\n📂 Тема: %s\\n\\n❓ Тип вопроса: %s\\n\\n💵 Стоимость: %s\\n\\n📋 Описание: %s\\n\\n📄 Вопрос: %s\\n\\n🔍 Ответ: %s\"\nВы точно хотите удалить вопрос?",
					TempPack[c.Sender().ID].Rounds[roundDel].Name,
					TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel].Name,
					questType[selectedQuest.Type],
					selectedQuest.Cost,
					selectedQuest.Description,
					selectedQuest.Text,
					selectedQuest.Answer,
				), &tele.SendOptions{
					ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
				})
			}

		}

		if strings.HasPrefix(c.Callback().Data, "SP4SIQ_") {
			packName := strings.TrimPrefix(c.Callback().Data, "SP4SIQ_")
			link, err := StorageDB.GetPack(c.Sender().ID, packName)
			if err != nil {
				return err
			}
			fileID, err := googleDrive.ExtractFileIDFromURL(link)
			if err != nil {
				return err
			}
			data, err := googleDrive.DownloadFileByID(fileID)
			pack, err := json.DataToPack(data)
			if err != nil {
				return err
			}
			xmlData, err := xml.ConvertPackToXML(bot, *pack)

			if err != nil {
				return err
			}
			if err1 := os.WriteFile(fmt.Sprintf("C:\\Users\\timof\\Desktop\\Packs\\%d\\content.xml", pack.PackID), xmlData, 0644); err != nil {
				return err1
			}

			err1 := utils.CreateZipArchive(fmt.Sprintf("C:\\Users\\timof\\Desktop\\Packs\\%d\\", pack.PackID), fmt.Sprintf("%d.siq", pack.PackID))
			if err1 != nil {
				log.Fatal(err1)
			}
			file, err := os.Open(fmt.Sprintf("C:\\Users\\timof\\Desktop\\Packs\\%d\\%d.siq", pack.PackID, pack.PackID))
			if err != nil {
				return fmt.Errorf("ошибка открытия файла: %w", err)
			}

			defer file.Close()

			tgFile := &tele.Document{
				File:     tele.FromReader(file),
				FileName: fmt.Sprintf("%s.siq", pack.PackName),
				Caption:  "Файл для проведения викторины в SiGame",
			}
			err2 := c.Bot().Delete(c.Message())
			if err != nil {
				return err2
			}
			err = c.Send(tgFile)
			if err != nil {
				return err
			}
			return nil

		}

		if strings.HasPrefix(c.Callback().Data[9:], "QDC_") {
			questNum, _ := strconv.Atoi(strings.TrimPrefix(c.Callback().Data[9:], "QDC_"))
			err1 := c.Bot().Delete(c.Message())
			if err1 != nil {
				return err1
			}
			TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel].Quests = append(TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel].Quests[:questNum], TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel].Quests[questNum+1:]...)
			TempPack[c.Sender().ID].QuestsCount -= 1

			if TempPack[c.Sender().ID].QuestsCount == 0 {
				menu.Inline(
					menu.Row(btnAddRound, btnDelRound),
					menu.Row(btnAddTheme, btnDelTheme),
					menu.Row(btnAddQuest),
					menu.Row(BtnSaveTmp),
					menu.Row(btnViewStruct),
				)
			} else {
				menu.Inline(
					menu.Row(btnAddRound, btnDelRound),
					menu.Row(btnAddTheme, btnDelTheme),
					menu.Row(btnAddQuest, btnDelQuest),
					menu.Row(BtnSaveTmp),
					menu.Row(btnViewStruct),
				)
			}

			if len(TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel].Quests) == 0 {
				userState[c.Sender().ID].SetPos(0, roundDel)
				userState[c.Sender().ID].SetPos(1, themeDel)
				userState[c.Sender().ID].SetPos(2, -1)

				saveKeyboard(c.Sender().ID, menu)
				return c.Send(fmt.Sprintf("Вопрос из темы '%s' из раунда '%s' удален", TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel].Name, TempPack[c.Sender().ID].Rounds[roundDel].Name), menu)
			}
			return c.Send(fmt.Sprintf("Вопрос из темы '%s' из раунда '%s' удален", TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel].Name, TempPack[c.Sender().ID].Rounds[roundDel].Name), storage.GetMessage(c.Chat().ID).Keyboard)
		}

		if strings.HasPrefix(c.Callback().Data[9:], "round_btn_del_confirm_") {
			selectedRound := strings.TrimPrefix(c.Callback().Data[9:], "round_btn_del_confirm_")
			for i, round := range TempPack[c.Sender().ID].Rounds {
				if round.Name == selectedRound {
					for _, theme := range round.Themes {
						TempPack[c.Sender().ID].QuestsCount -= len(theme.Quests)
					}
					TempPack[c.Sender().ID].Rounds = append(TempPack[c.Sender().ID].Rounds[:i], TempPack[c.Sender().ID].Rounds[i+1:]...)
					break
				}
			}
			err1 := c.Bot().Delete(c.Message())
			if err1 != nil {
				return err1
			}
			if len(TempPack[c.Sender().ID].Rounds) == 0 {
				userState[c.Sender().ID].SetPos(0, -1)
				userState[c.Sender().ID].SetPos(1, -1)
				userState[c.Sender().ID].SetPos(2, -1)
				menu.Inline(
					menu.Row(btnAddRound),
					menu.Row(BtnSaveTmp),
				)
				saveKeyboard(c.Sender().ID, menu)
				return c.Send(fmt.Sprintf("Раунд '%s' удалён.", selectedRound), menu)
			}
			if TempPack[c.Sender().ID].ThemesCount == 0 {
				menu.Inline(
					menu.Row(btnAddRound, btnDelRound),
					menu.Row(btnAddTheme),
					menu.Row(BtnSaveTmp),
					menu.Row(btnViewStruct),
				)
			} else {
				if TempPack[c.Sender().ID].QuestsCount == 0 {
					menu.Inline(
						menu.Row(btnAddRound, btnDelRound),
						menu.Row(btnAddTheme, btnDelTheme),
						menu.Row(btnAddQuest),
						menu.Row(BtnSaveTmp),
						menu.Row(btnViewStruct),
					)
				} else {
					menu.Inline(
						menu.Row(btnAddRound, btnDelRound),
						menu.Row(btnAddTheme, btnDelTheme),
						menu.Row(btnAddQuest, btnDelQuest),
						menu.Row(BtnSaveTmp),
						menu.Row(btnViewStruct),
						menu.Row(btnPublish),
					)
				}
			}

			return c.Send(fmt.Sprintf("Раунд '%s' удалён.", selectedRound), menu)
		}

		if strings.HasPrefix(c.Callback().Data[9:], "TDC_") {
			selectedTheme := strings.TrimPrefix(c.Callback().Data[9:], "TDC_")
			for i, theme := range TempPack[c.Sender().ID].Rounds[roundDel].Themes {
				if theme.Name == selectedTheme {
					themeDel = i
					break
				}
			}

			TempPack[c.Sender().ID].QuestsCount -= len(TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel].Quests)
			TempPack[c.Sender().ID].Rounds[roundDel].Themes = append(TempPack[c.Sender().ID].Rounds[roundDel].Themes[:themeDel], TempPack[c.Sender().ID].Rounds[roundDel].Themes[themeDel+1:]...)
			TempPack[c.Sender().ID].ThemesCount -= 1
			err1 := c.Bot().Delete(c.Message())
			if err1 != nil {
				return err1
			}

			if TempPack[c.Sender().ID].ThemesCount == 0 {
				menu.Inline(
					menu.Row(btnAddRound, btnDelRound),
					menu.Row(btnAddTheme),
					menu.Row(BtnSaveTmp),
				)
			} else {
				if TempPack[c.Sender().ID].QuestsCount == 0 {
					menu.Inline(
						menu.Row(btnAddRound, btnDelRound),
						menu.Row(btnAddTheme, btnDelTheme),
						menu.Row(btnAddQuest),
						menu.Row(BtnSaveTmp),
					)
				} else {
					menu.Inline(
						menu.Row(btnAddRound, btnDelRound),
						menu.Row(btnAddTheme, btnDelTheme),
						menu.Row(btnAddQuest, btnDelQuest),
						menu.Row(BtnSaveTmp),
						menu.Row(btnViewStruct),
						menu.Row(btnPublish),
					)
				}

			}

			if len(TempPack[c.Sender().ID].Rounds[roundDel].Themes) == 0 {
				userState[c.Sender().ID].SetPos(0, roundDel)
				userState[c.Sender().ID].SetPos(1, -1)
				userState[c.Sender().ID].SetPos(2, -1)

				saveKeyboard(c.Sender().ID, menu)
				return c.Send(fmt.Sprintf("Тема '%s' из раунда '%s' удалена", TempPack[c.Sender().ID].Rounds[roundDel].Name, selectedTheme), menu)
			}

			return c.Send(fmt.Sprintf("Тема '%s' из раунда '%s' удалена", TempPack[c.Sender().ID].Rounds[roundDel].Name, selectedTheme), menu)
		}

		return nil
	})

}
