package validators

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

func IsValidCost(cost string) error {

	for _, r := range cost {
		if !(unicode.IsDigit(r)) {
			return errors.New("цена должна состоять из цифр")
		}
	}

	if utf8.RuneCountInString(cost) > 4 {
		return errors.New("цена не может быть выше 9999")
	}
	return nil
}
