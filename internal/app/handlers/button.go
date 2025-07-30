package handlers

import (
	"OwnGamePack/internal/app/utils"
	"OwnGamePack/internal/app/validators"
	"OwnGamePack/internal/storage"
	"OwnGamePack/internal/structs"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
)

var (
	questType      = map[string]string{"default": "обычный", "bet": "со ставкой", "secret": "с секретом"}
	BtnSaveTmp     = menu.Data("📌 Сохранить пак", "saveTmp")
	btnViewStruct  = menu.Data("🗂️ Показать структуру пака", "viewStruct")
	btnAddRound    = menu.Data("🟢 Добавить раунд", "addRound")
	btnAddTheme    = menu.Data("🟢 Добавить тему", "addTheme")
	btnAddQuest    = menu.Data("🟢 Добавить вопрос", "addQuest")
	btnDelRound    = menu.Data("🗑️ Удалить раунд", "delRound")
	btnDelTheme    = menu.Data("🗑️ Удалить тему", "delTheme")
	btnDelQuest    = menu.Data("🗑️ Удалить вопрос", "delQuest")
	btnQuestDef    = menu.Data("🔘 Обычный", "questDef")
	btnQuestBet    = menu.Data("🎰 Со ставкой", "questBet")
	btnQuestSecret = menu.Data("🔑 С секретом", "questSec")
	btnPublish     = menu.Data("📢 Опубликовать пак", "publish")
	btnBack        = tele.InlineButton{
		Unique: "back",
		Text:   "⬅️ Вернуться назад",
	}
	btnConfirm = tele.InlineButton{
		Unique: "confirm",
		Text:   "✅ Подтвердить",
		Data:   "",
	}
)

func RegisterButtonHandlers(bot *tele.Bot) {

	bot.Handle(&btnPublish, func(c tele.Context) error {
		saveKeyboard(c.Sender().ID, menu)
		lastMenu := getKeyboard(c.Sender().ID)
		err := validators.IsValidPack(TempPack[c.Sender().ID])
		if err != nil {
			return c.Send(err.Error(), lastMenu)
		}

		grid := utils.BuildTagMenu(TempPack[c.Sender().ID])
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}
		err2 := c.Send("Выберите тэги для пака. После этого пропишите команду /publish\n⚠️ПОСЛЕ КОМАНДЫ ПАК НЕЛЬЗЯ РЕДАКТИРОВАТЬ⚠️", &tele.ReplyMarkup{InlineKeyboard: grid})
		return err2

	})

	bot.Handle(&btnDelQuest, func(c tele.Context) error {
		var keyboard [][]tele.InlineButton
		for _, round := range TempPack[c.Sender().ID].Rounds {
			for _, theme := range round.Themes {
				if len(theme.Quests) != 0 {
					roundName := round.Name
					btn := tele.InlineButton{
						Data: "SR4DQ_" + roundName,
						Text: roundName,
					}
					keyboard = append(keyboard, []tele.InlineButton{btn})
				}
			}

		}
		keyboard = append(keyboard, []tele.InlineButton{btnBack})
		storage.SaveMessage(c.Chat().ID, c.Message().Text, c.Message().ReplyMarkup)
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}
		return c.Send("Выберите раунд, в котором хотите удалить вопрос", &tele.SendOptions{
			ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
		})
	})

	bot.Handle(&btnBack, func(c tele.Context) error {
		msg := storage.GetMessage(c.Chat().ID)
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}

		err := c.Send(msg.Text, msg.Keyboard)
		return err
	})

	bot.Handle(&btnAddRound, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		userState[c.Sender().ID].SetState("add_round")
		_, err := c.Bot().Edit(c.Message(), "✏️ Введите название раунда", &tele.ReplyMarkup{})
		return err
	})

	bot.Handle(&btnAddTheme, func(c tele.Context) error {

		var keyboard [][]tele.InlineButton
		for _, round := range TempPack[c.Sender().ID].Rounds {
			roundName := round.Name
			btn := tele.InlineButton{
				Data: "select_round_add_theme_" + roundName,
				Text: roundName,
			}
			keyboard = append(keyboard, []tele.InlineButton{btn})
		}
		storage.SaveMessage(c.Chat().ID, c.Message().Text, c.Message().ReplyMarkup)
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}
		return c.Send("Выберите раунд для добавления темы", &tele.SendOptions{
			ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
		})

	})

	bot.Handle(&btnAddQuest, func(c tele.Context) error {
		var keyboard [][]tele.InlineButton
		for _, round := range TempPack[c.Sender().ID].Rounds {

			if len(round.Themes) > 0 {
				roundName := round.Name
				btn := tele.InlineButton{
					Data: "select_round_add_quest_" + roundName,
					Text: roundName,
				}
				keyboard = append(keyboard, []tele.InlineButton{btn})
			}

		}
		storage.SaveMessage(c.Chat().ID, c.Message().Text, c.Message().ReplyMarkup)
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}
		return c.Send("Выберите раунд для добавления вопроса", &tele.SendOptions{
			ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
		})
	})

	bot.Handle(&btnQuestDef, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		quest := structs.NewQuest("default")
		pos0 := userState[c.Sender().ID].GetPos(0)
		pos1 := userState[c.Sender().ID].GetPos(1)
		pos2 := userState[c.Sender().ID].GetPos(2)

		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests = append(TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests, quest)
		userState[c.Sender().ID].SetPos(2, len(TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests)-1)
		pos2 = userState[c.Sender().ID].GetPos(2)
		_, err := c.Bot().Edit(c.Message(),
			fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n❓ Тип вопроса: %s\n\n✏️ Укажите стоимость вопроса",
				TempPack[c.Sender().ID].Rounds[pos0].Name,
				TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Name,
				questType[TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Type]),
			&tele.ReplyMarkup{})
		userState[c.Sender().ID].SetState("awaiting_price")
		return err
	})

	bot.Handle(&btnQuestBet, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		quest := structs.NewQuest("bet")
		pos0 := userState[c.Sender().ID].GetPos(0)
		pos1 := userState[c.Sender().ID].GetPos(1)
		pos2 := userState[c.Sender().ID].GetPos(2)
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests = append(TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests, quest)
		userState[c.Sender().ID].SetPos(2, len(TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests)-1)
		pos2 = userState[c.Sender().ID].GetPos(2)
		_, err := c.Bot().Edit(c.Message(),
			fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n❓ Тип вопроса: %s\n\n✏️ Укажите стоимость вопроса",
				TempPack[c.Sender().ID].Rounds[pos0].Name,
				TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Name,
				questType[TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Type]),
			&tele.ReplyMarkup{})
		userState[c.Sender().ID].SetState("awaiting_price")
		return err
	})

	bot.Handle(&btnQuestSecret, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		pos0 := userState[c.Sender().ID].GetPos(0)
		pos1 := userState[c.Sender().ID].GetPos(1)
		pos2 := userState[c.Sender().ID].GetPos(2)
		quest := structs.NewQuest("secret")
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests = append(TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests, quest)
		userState[c.Sender().ID].SetPos(2, len(TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests)-1)
		pos2 = userState[c.Sender().ID].GetPos(2)
		_, err := c.Bot().Edit(c.Message(), fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n❓ Тип вопроса: %s\n\n✏️ Укажите стоимость вопроса",
			TempPack[c.Sender().ID].Rounds[pos0].Name,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Name,
			questType[TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Type]), &tele.ReplyMarkup{})
		userState[c.Sender().ID].SetState("awaiting_price")
		return err
	})

	bot.Handle(&btnDelRound, func(c tele.Context) error {
		var keyboard [][]tele.InlineButton
		for _, round := range TempPack[c.Sender().ID].Rounds {
			roundName := round.Name
			btn := tele.InlineButton{
				Data: "select_round_btn_del_" + roundName,
				Text: roundName,
			}
			keyboard = append(keyboard, []tele.InlineButton{btn})
		}
		keyboard = append(keyboard, []tele.InlineButton{btnBack})
		storage.SaveMessage(c.Chat().ID, c.Message().Text, c.Message().ReplyMarkup)
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}
		return c.Send("Выберите раунд для удаления", &tele.SendOptions{
			ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
		})
	})

	bot.Handle(&btnDelTheme, func(c tele.Context) error {
		var keyboard [][]tele.InlineButton
		for _, round := range TempPack[c.Sender().ID].Rounds {
			if len(round.Themes) > 0 {
				roundName := round.Name
				btn := tele.InlineButton{
					Data: "SR4DT_" + roundName,
					Text: roundName,
				}
				keyboard = append(keyboard, []tele.InlineButton{btn})
			}
		}
		keyboard = append(keyboard, []tele.InlineButton{btnBack})
		storage.SaveMessage(c.Chat().ID, c.Message().Text, c.Message().ReplyMarkup)
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}
		return c.Send("Выберите раунд из которого Вы хотите удалить тему", &tele.SendOptions{
			ReplyMarkup: &tele.ReplyMarkup{InlineKeyboard: keyboard},
		})

	})

	bot.Handle(&btnViewStruct, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()

		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}

		mess := ""
		if len(savePack.Rounds) != 0 {
			for i, round := range savePack.Rounds {
				mess += fmt.Sprintf("📌 %s:\n", round.Name)
				if len(savePack.Rounds[i].Themes) != 0 {
					for j, theme := range savePack.Rounds[i].Themes {
						mess += fmt.Sprintf("     📖%s\n", theme.Name)
						if len(savePack.Rounds[i].Themes[j].Quests) != 0 {
							for _, quest := range savePack.Rounds[i].Themes[j].Quests {
								mess += fmt.Sprintf("           %s\n", questType[quest.Type]+", цена: "+quest.Cost)
							}
						}
					}

				}
			}
		} else {

			return c.Send("❌ Пока ничего не добавлено. Возможно Вы не сохранили изменения.", getKeyboard(c.Sender().ID))
		}

		return c.Send(fmt.Sprintf("🗂️ Структура пака:\n\n%s", mess), getKeyboard(c.Sender().ID))

	})

}
