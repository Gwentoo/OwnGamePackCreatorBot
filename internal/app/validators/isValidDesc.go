package validators

import (
	"errors"
	"unicode/utf8"
)

func IsValidQuestDesc(desc string) error {
	if utf8.RuneCountInString(desc) > 100 {
		return errors.New("описание вопроса должно быть не длиннее 100 символов")
	}
	return nil
}

func IsValidPackDesc(desc string) error {
	if utf8.RuneCountInString(desc) > 150 {
		return errors.New("описание пака должно быть не длиннее 150 символов")
	}
	return nil
}
