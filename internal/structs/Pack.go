package structs

type Pack struct {
	UserID      int64    `json:"userID"`
	UserName    string   `json:"userName"`
	PackID      int64    `json:"packID"`
	PackName    string   `json:"packName"`
	PackDesc    string   `json:"packDesc"`
	ThemesCount int      `json:"themesCount"`
	QuestsCount int      `json:"questsCount"`
	PackTags    []string `json:"packTags"`
	Rounds      []Round  `json:"roundName"`
}

func NewPack() Pack {
	return Pack{}
}
