package validators

import (
	"OwnGamePack/internal/structs"
	"errors"
	"fmt"
)

func IsValidPack(pack *structs.Pack) error {
	var roundWOTheme []string
	var themeWOQuest = make(map[string]string, 1)
	for _, round := range pack.Rounds {
		if len(round.Themes) == 0 {
			roundWOTheme = append(roundWOTheme, round.Name)
		}
		for _, theme := range round.Themes {
			if len(theme.Quests) == 0 {
				themeWOQuest[round.Name] = theme.Name
			}
		}
	}

	if len(roundWOTheme) > 0 {
		var rounds = ""
		for _, roundName := range roundWOTheme {
			rounds += roundName + "\n"
		}
		return errors.New(fmt.Sprintf("❌ В каждом раунде должна быть хотя бы 1 тема.\nРаунды без тем:\n%s", rounds))
	}
	if len(themeWOQuest) > 0 {
		var themes = ""
		for roundName, themeName := range themeWOQuest {
			themes += fmt.Sprintf("Раунд: %s, тема: %s", roundName, themeName)
		}
		return errors.New(fmt.Sprintf("❌ В каждой тема должен быть хотя бы 1 вопрос.\nТемы без вопросов:\n%s", themes))
	}
	return nil
}
