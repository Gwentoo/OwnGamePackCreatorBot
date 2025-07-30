package handlers

import (
	"OwnGamePack/internal/app/generatePackID"
	"OwnGamePack/internal/app/validators"
	"OwnGamePack/internal/storage"
	"OwnGamePack/internal/structs"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"sync"
)

var StorageDB *storage.Storage

var (
	userKeyboards = make(map[int64]*tele.ReplyMarkup)
	keyboardMutex sync.Mutex
)

func saveKeyboard(userID int64, menu *tele.ReplyMarkup) {
	keyboardMutex.Lock()
	defer keyboardMutex.Unlock()
	userKeyboards[userID] = menu
}

func getKeyboard(userID int64) *tele.ReplyMarkup {
	keyboardMutex.Lock()
	defer keyboardMutex.Unlock()
	return userKeyboards[userID]
}

func RegisterCommonHandlers(bot *tele.Bot) {
	bot.Handle("/info", handleInfo)
	bot.Handle("/newpack", handleNewPack)
	bot.Handle("/support", handleSupport)
}

var TempPack = make(map[int64]*structs.Pack, 10)
var userState = make(map[int64]*structs.UserState, 10)

var (
	savePack = structs.NewPack()
	menu     = &tele.ReplyMarkup{}
)

func handleSupport(c tele.Context) error {
	return nil
}

func handleInfo(c tele.Context) error {
	return c.Send(
		"Этот бот создан для создания/редактирования паков для викторины \"Своя игра\"\n" +
			"/newpack - Начать создание нового пака\n" +
			//"/edit - Редактирование паков, которые еще не были добавлены в общий доступ\n" +
			"/packs - Посмотреть список всех созданных Вами паков\n" +
			"/siq - Преобразование уже опубликованного Вами пака в файл для проведения SiGame\n" +
			"/support - Обратиться в поддержку, если возникли какие-либо проблемы")
}

func handleNewPack(c tele.Context) error {

	btnCancelPack := menu.Data("Отмена", "cancelPack")
	btnContinuePack := menu.Data("Продолжить", "continuePack")
	c.Bot().Handle(&btnContinuePack, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		newPack := structs.NewPack()
		newUserState := structs.NewUserState()
		TempPack[c.Sender().ID] = &newPack
		TempPack[c.Sender().ID].UserID = c.Sender().ID
		userState[c.Sender().ID] = &newUserState
		userState[c.Sender().ID].SetState("awaiting_pack_name")
		_, err := c.Bot().Edit(c.Message(), "✏️ Введите название пака", &tele.ReplyMarkup{})
		return err
	})
	c.Bot().Handle(&btnCancelPack, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		_, err := c.Bot().Edit(c.Message(), "❌ Создание пакета отменено.", &tele.ReplyMarkup{})
		return err
	})

	menu.Inline(
		menu.Row(btnCancelPack, btnContinuePack),
	)
	return c.Send("Вы начали создание нового пака", menu)
}

func handlePackName(c tele.Context) error {

	err := validators.IsValidName(c.Text())

	names, err1 := StorageDB.GetPacksName(c.Sender().ID)
	if err1 != nil {
		log.Println(err1)
	}
	for k := range names {
		if c.Text() == k {
			return c.Send("❌ У вас уже есть пак с таким названием\nПопробуйте еще раз:")
		}
	}

	if err != nil {
		return c.Send(fmt.Sprintf("❌ %s\nПопробуйте еще раз:", err))
	}
	packID, err := generatePackID.GeneratePackID()
	if err != nil {
		return err
	}
	TempPack[c.Sender().ID].PackID = packID
	TempPack[c.Sender().ID].PackName = c.Text()
	TempPack[c.Sender().ID].UserName = "@" + c.Sender().Username
	userState[c.Sender().ID].SetState("add_desc")

	return c.Send("✏️ Введите описание пака", &tele.ReplyMarkup{})
}

func handlePackDesc(c tele.Context) error {
	err := validators.IsValidPackDesc(c.Text())
	if err != nil {
		return c.Send(fmt.Sprintf("❌ %s\nПопробуйте еще раз:", err))
	}
	TempPack[c.Sender().ID].PackDesc = c.Text()

	userState[c.Sender().ID].SetState("add_round")
	menu.Inline(
		menu.Row(btnAddRound),
		menu.Row(BtnSaveTmp),
	)
	saveKeyboard(c.Sender().ID, menu)
	return c.Send(fmt.Sprintf("Название пака: '%s'\nОписание пака: '%s'", TempPack[c.Sender().ID].PackName, TempPack[c.Sender().ID].PackDesc), menu)
}

func handleNewRound(c tele.Context) error {

	err := validators.IsValidRoundName(c.Text(), TempPack[c.Sender().ID])
	if err != nil {
		return c.Send(fmt.Sprintf("❌ %s\nПопробуйте еще раз:", err))
	}

	round := structs.NewRound(c.Text())
	TempPack[c.Sender().ID].Rounds = append(TempPack[c.Sender().ID].Rounds, round)
	pos0 := userState[c.Sender().ID].GetPos(0)
	userState[c.Sender().ID].SetPos(0, pos0+1)
	userState[c.Sender().ID].SetPos(1, -1)
	userState[c.Sender().ID].SetPos(2, -1)
	userState[c.Sender().ID].SetState("")

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

	saveKeyboard(c.Sender().ID, menu)

	return c.Send(fmt.Sprintf("✅ Добавлен раунд '%s'", c.Text()), menu)

}

func handleNewTheme(c tele.Context) error {

	err := validators.IsValidThemeName(c.Text(), TempPack[c.Sender().ID], userState[c.Sender().ID].GetPos(0))
	if err != nil {
		return c.Send(fmt.Sprintf("❌ %s\nПопробуйте еще раз:", err))
	}

	theme := structs.NewTheme(c.Text())
	pos0 := userState[c.Sender().ID].GetPos(0)
	TempPack[c.Sender().ID].Rounds[pos0].Themes = append(TempPack[c.Sender().ID].Rounds[pos0].Themes, theme)
	pos1 := userState[c.Sender().ID].GetPos(1)
	userState[c.Sender().ID].SetPos(1, pos1+1)
	userState[c.Sender().ID].SetPos(2, -1)
	userState[c.Sender().ID].SetState("")

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

	saveKeyboard(c.Sender().ID, menu)
	TempPack[c.Sender().ID].ThemesCount += 1
	return c.Send(fmt.Sprintf("✅ Добавлена тема '%s' в раунд '%s'", c.Text(), TempPack[c.Sender().ID].Rounds[pos0].Name), menu)

}

func handleNewCost(c tele.Context) error {

	err := validators.IsValidCost(c.Text())
	if err != nil {
		return c.Send(fmt.Sprintf("❌ %s\nПопробуйте еще раз:", err))
	}

	userState[c.Sender().ID].SetState("awaiting_quest_disc")
	pos0 := userState[c.Sender().ID].GetPos(0)
	pos1 := userState[c.Sender().ID].GetPos(1)
	pos2 := userState[c.Sender().ID].GetPos(2)
	TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Cost = c.Text()
	return c.Send(
		fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n❓ Тип вопроса: %s\n\n💵 Стоимость: %s\n\n✏️ Укажите описание вопроса\n(Например, назовите изображённое животное)",
			TempPack[c.Sender().ID].Rounds[pos0].Name,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Name,
			questType[TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Type],
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Cost,
		))

}

func handleNewQuestDisc(c tele.Context) error {

	err := validators.IsValidQuestDesc(c.Text())
	if err != nil {
		return c.Send(fmt.Sprintf("❌ %s\nПопробуйте еще раз:", err))
	}

	pos0 := userState[c.Sender().ID].GetPos(0)
	pos1 := userState[c.Sender().ID].GetPos(1)
	pos2 := userState[c.Sender().ID].GetPos(2)
	TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Description = c.Text()
	userState[c.Sender().ID].SetState("awaiting_quest_content")

	btnNewText := menu.Data("📄 Текст", "newText")
	btnNewAudio := menu.Data("🎵 Аудио", "newAudio")
	btnNewVideo := menu.Data("🎞️ Видео", "newVideo")
	btnNewPhoto := menu.Data("📷 Фото", "newPhoto")

	c.Bot().Handle(&btnNewText, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		userState[c.Sender().ID].SetState("awaiting_text")
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].ContentType = "text"
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}
		return c.Send("✏️ Введите текст вопроса", &tele.ReplyMarkup{})
	})
	c.Bot().Handle(&btnNewAudio, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		userState[c.Sender().ID].SetState("awaiting_audio")
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].ContentType = "audio"
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}
		return c.Send("🎵 Отправьте аудиозапись (продолжительность не более 20с)", &tele.ReplyMarkup{})
	})
	c.Bot().Handle(&btnNewVideo, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		userState[c.Sender().ID].SetState("awaiting_video")
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].ContentType = "video"
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}
		return c.Send("🎞️ Отправьте видеозапись (продолжительность не более 20с)", &tele.ReplyMarkup{})
	})
	c.Bot().Handle(&btnNewPhoto, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		userState[c.Sender().ID].SetState("awaiting_photo")
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].ContentType = "photo"
		err1 := c.Bot().Delete(c.Message())
		if err1 != nil {
			return err1
		}
		return c.Send("📷 Отправьте изображение", &tele.ReplyMarkup{})
	})

	menu.Inline(
		menu.Row(btnNewText, btnNewAudio),
		menu.Row(btnNewVideo, btnNewPhoto),
	)

	return c.Send(fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n❓ Тип вопроса: %s\n\n💵 Стоимость: %s\n\n📋 Описание: %s\n\n✏️ Выберите тип содержимого вопроса",
		TempPack[c.Sender().ID].Rounds[pos0].Name,
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Name,
		questType[TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Type],
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Cost,
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Description,
	),
		menu)

}

func handleNewText(c tele.Context) error {

	pos0 := userState[c.Sender().ID].GetPos(0)
	pos1 := userState[c.Sender().ID].GetPos(1)
	pos2 := userState[c.Sender().ID].GetPos(2)

	TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Text = c.Text()
	userState[c.Sender().ID].SetState("awaiting_answer")
	return c.Send(fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n❓ Тип вопроса: %s\n\n💵 Стоимость: %s\n\n📋 Описание: %s\n\n📄 Текст вопроса: %s\n\n ✏️ Укажите ответ на вопрос",
		TempPack[c.Sender().ID].Rounds[pos0].Name,
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Name,
		questType[TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Type],
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Cost,
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Description,
		TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Text,
	))
}

func handleNewAnswer(c tele.Context) error {
	pos0 := userState[c.Sender().ID].GetPos(0)
	pos1 := userState[c.Sender().ID].GetPos(1)
	pos2 := userState[c.Sender().ID].GetPos(2)
	TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Answer = c.Text()

	menu.Inline(
		menu.Row(btnAddRound, btnDelRound),
		menu.Row(btnAddTheme, btnDelTheme),
		menu.Row(btnAddQuest, btnDelQuest),
		menu.Row(BtnSaveTmp),
		menu.Row(btnViewStruct),
		menu.Row(btnPublish),
	)
	userState[c.Sender().ID].SetState("")

	switch TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].ContentType {

	case "photo":
		photo := TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Photo
		photo.Caption = fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n❓ Тип вопроса: %s\n\n💵 Стоимость: %s\n\n📋 Описание: %s\n\n🔍 Ответ: %s",
			TempPack[c.Sender().ID].Rounds[pos0].Name,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Name,
			questType[TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Type],
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Cost,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Description,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Answer,
		)
		err1 := c.Send(photo)
		if err1 != nil {
			return err1
		}
	case "video":
		video := TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Video
		video.Caption = fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n❓ Тип вопроса: %s\n\n💵 Стоимость: %s\n\n📋 Описание: %s\n\n🔍 Ответ: %s",
			TempPack[c.Sender().ID].Rounds[pos0].Name,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Name,
			questType[TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Type],
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Cost,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Description,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Answer,
		)
		err1 := c.Send(video)
		if err1 != nil {
			return err1
		}
	case "audio":
		audio := TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Audio
		audio.Caption = fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n❓ Тип вопроса: %s\n\n💵 Стоимость: %s\n\n📋 Описание: %s\n\n🔍 Ответ: %s",
			TempPack[c.Sender().ID].Rounds[pos0].Name,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Name,
			questType[TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Type],
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Cost,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Description,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Answer,
		)
		err1 := c.Send(audio)
		if err1 != nil {
			return err1
		}
	default:
		text := TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Text

		err1 := c.Send(fmt.Sprintf("📍 Раунд: %s\n\n📂 Тема: %s\n\n❓ Тип вопроса: %s\n\n💵 Стоимость: %s\n\n📋 Описание: %s\n\n📄 Вопрос: %s\n\n🔍 Ответ: %s",
			TempPack[c.Sender().ID].Rounds[pos0].Name,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Name,
			questType[TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Type],
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Cost,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Description,
			text,
			TempPack[c.Sender().ID].Rounds[pos0].Themes[pos1].Quests[pos2].Answer))
		if err1 != nil {
			return err1
		}
	}

	saveKeyboard(c.Sender().ID, menu)
	TempPack[c.Sender().ID].QuestsCount += 1
	return c.Send("✅ Вопрос добавлен", menu)
}
