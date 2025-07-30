package utils

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
)

func DownloadFile(bot *tele.Bot, fileID string, savePath string) error {
	file, err := bot.FileByID(fileID)
	if err != nil {
		return fmt.Errorf("ошибка получения файла: %w", err)
	}
	err2 := bot.Download(&file, savePath)
	if err2 != nil {
		return fmt.Errorf("ошибка скачивания: %w", err2)
	}

	return nil
}
