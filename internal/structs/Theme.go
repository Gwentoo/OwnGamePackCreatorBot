package structs

type Theme struct {
	Name   string  `json:"themeName"`
	Quests []Quest `json:"quests"`
}

func NewTheme(name string) Theme {
	return Theme{Name: name}
}
