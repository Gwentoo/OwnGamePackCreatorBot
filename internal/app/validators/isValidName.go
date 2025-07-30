package validators

import (
	"OwnGamePack/internal/structs"
	"errors"
	"unicode"
	"unicode/utf8"
)

func IsValidName(name string) error {

	if utf8.RuneCountInString(name) > 25 {
		return errors.New("название не должно быть длиннее 25 символов")
	}
	for _, r := range name {

		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || r == '/') {
			return errors.New("название может состоять только из букв, цифр и знаков препинания")
		}

	}
	return nil
}

func IsValidRoundName(name string, pack *structs.Pack) error {
	err := IsValidName(name)
	if err != nil {
		return err
	}
	for _, round := range pack.Rounds {
		if round.Name == name {
			return errors.New("такое название раунда уже есть")
		}
	}
	return nil
}

func IsValidThemeName(name string, pack *structs.Pack, roundNum int) error {
	err := IsValidName(name)
	if err != nil {
		return err
	}
	for _, theme := range pack.Rounds[roundNum].Themes {
		if theme.Name == name {
			return errors.New("такое название темы уже есть в этом раунде")
		}
	}
	return nil
}
