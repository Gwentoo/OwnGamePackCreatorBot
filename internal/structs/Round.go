package structs

type Round struct {
	Name   string  `json:"roundName"`
	Themes []Theme `json:"themes"`
}

func NewRound(name string) Round {
	return Round{Name: name}
}
