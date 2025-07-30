package validators

import (
	"errors"
	"gopkg.in/telebot.v3"
)

func IsValidPhoto(photo *telebot.Photo) error {
	if photo.FileSize > 5<<20 {
		return errors.New("слишком большая фотография")
	}
	return nil
}

func IsValidVideo(video *telebot.Video) error {
	if video.Duration > 20 {
		return errors.New("видео не должно быть дольше 15 секунд")
	}
	return nil
}

func IsValidAudio(audio *telebot.Audio) error {
	if audio.Duration > 20 {
		return errors.New("аудио не должно быть дольше 15 секунд")
	}
	return nil
}
