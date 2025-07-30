package bot

import (
	"OwnGamePack/config"
	"OwnGamePack/internal/app/generatePackID"
	"OwnGamePack/internal/app/handlers"
	"OwnGamePack/internal/storage"
	"fmt"
	_ "github.com/lib/pq"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
)

type MyBot struct {
	Bot     *tele.Bot
	Storage *storage.Storage
}

func registerHandlers(bot *tele.Bot) {
	handlers.RegisterCommonHandlers(bot)
	handlers.RegisterTextHandlers(bot)
	handlers.RegisterButtonHandlers(bot)
	handlers.RegisterPhotoHandlers(bot)
	handlers.RegisterVideoHandlers(bot)
	handlers.RegisterAudioHandlers(bot)
	handlers.RegisterCallbackHandlers(bot)
}

func NewBot(cfg *config.Config, storage *storage.Storage) (*MyBot, error) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  cfg.Bot.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	b := &MyBot{
		Bot:     bot,
		Storage: storage,
	}

	if err != nil {
		panic(err)
	}

	registerHandlers(b.Bot)
	b.registerHandlers()
	commands := []tele.Command{
		{
			Text:        "info",
			Description: "Информация",
		},
		{
			Text:        "newpack",
			Description: "Создание нового пака",
		},
		{
			Text:        "siq",
			Description: "Сохранить пак в .siq",
		},
		{
			Text:        "packs",
			Description: "Список паков",
		},
		{
			Text:        "support",
			Description: "Обратиться в поддержку",
		},
	}
	err2 := b.Bot.SetCommands(commands)
	if err2 != nil {
		return nil, err2
	}

	return b, nil
}

func (b *MyBot) registerHandlers() {
	b.Bot.Handle("/start", b.handleStart)
	b.Bot.Handle("/siq", b.handleSiq)
	b.Bot.Handle("/packs", b.handlePacks)
	b.Bot.Handle("/save", b.HandleSave)
	b.Bot.Handle("/publish", b.handlePublish)
	b.Bot.Handle(&handlers.BtnSaveTmp, func(c tele.Context) error {
		defer func() {
			if err := c.Respond(); err != nil {
				log.Printf("Ошибка при ответе: %v", err)
			}
		}()
		if handlers.TempPack[c.Sender().ID].PackID == 0 {
			packID, err := generatePackID.GeneratePackID()
			if err != nil {
				return err
			}
			handlers.TempPack[c.Sender().ID].PackID = packID
		}

		err := b.Storage.SavePack(handlers.TempPack[c.Sender().ID], false)
		if err != nil {
			return err
		}

		return c.Send("✅ Изменения сохранены")

	})

}

func (b *MyBot) handleStart(c tele.Context) error {

	userID := c.Sender().ID
	exists, err := b.Storage.CheckUserExists(userID)
	if err != nil {
		log.Printf("Error checking user: %v", err)
		return err
	}

	if !exists {
		err = b.Storage.SaveUser(userID)
		if err != nil {
			log.Printf("Error adding user: %v", err)
			return err
		}
		err2 := c.Send("Добро пожаловать!")
		if err2 != nil {
			return err2
		}
	} else {
		err3 := c.Send("С возвращением!")
		if err3 != nil {
			return err3
		}
	}

	return err
}

func (b *MyBot) handleSiq(c tele.Context) error {
	packs, err := b.Storage.GetPacksName(c.Sender().ID)
	if err != nil {
		log.Println(err)
	}
	if len(packs) == 0 {
		return c.Send("У вас ещё нет опубликованных паков")
	}
	var keyboard [][]tele.InlineButton
	for k, v := range packs {
		if v {
			btn := tele.InlineButton{

				Text: k,
				Data: fmt.Sprintf("SP4SIQ_%s", k),
			}
			keyboard = append(keyboard, []tele.InlineButton{btn})
		}
	}

	return c.Send("Выберите пак для сохранения в .siq", &tele.ReplyMarkup{
		InlineKeyboard: keyboard,
	})
}

func (b *MyBot) handlePacks(c tele.Context) error {
	names, err := b.Storage.GetPacksName(c.Sender().ID)
	if err != nil {
		return err
	}
	published := ""
	unPublished := ""
	for name, public := range names {
		if public {
			published += name + "\n"
		} else {
			unPublished += name + "\n"
		}
	}

	if published == "" && unPublished == "" {
		err2 := c.Send("Вы ещё не создавали паки. Чтобы начать, напишите /newpack")
		return err2
	}
	err2 := c.Send(fmt.Sprintf("Опубликованные паки:\n%s\nСоздающиеся паки:\n%s", published, unPublished))
	return err2
}

func (b *MyBot) HandleSave(c tele.Context) error {

	err := b.Storage.SavePack(handlers.TempPack[c.Sender().ID], false)
	if err != nil {
		return err
	}
	return c.Send("Пак сохранён.")
}

func (b *MyBot) handlePublish(c tele.Context) error {

	err2 := b.Storage.SavePack(handlers.TempPack[c.Sender().ID], true)
	if err2 != nil {
		return err2
	}
	err3 := c.Bot().Delete(c.Message())
	if err3 != nil {
		return err3
	}
	err := c.Send("Пак опубликован. Спасибо, что пользуетесь ботом ❤️")
	return err
}
