package xml

import "encoding/xml"

type XMLPackage struct {
	XMLName    xml.Name `xml:"package"`
	Name       string   `xml:"name,attr"`
	Version    string   `xml:"version,attr"`
	ID         string   `xml:"id,attr"`
	Date       string   `xml:"date,attr"`
	Difficulty string   `xml:"difficulty,attr"`
	XMLNS      string   `xml:"xmlns,attr"`
	Tags       struct {
		Tag []string `xml:"tag"`
	} `xml:"tags"`
	Info struct {
		Authors struct {
			Author []string `xml:"author"`
		} `xml:"authors"`
	} `xml:"info"`
	Rounds []XMLRound `xml:"rounds>round"`
}

type XMLRound struct {
	Name   string     `xml:"name,attr"`
	Type   string     `xml:"type,attr,omitempty"`
	Themes []XMLTheme `xml:"themes>theme"`
}

type XMLTheme struct {
	Name string `xml:"name,attr"`
	Info *struct {
		Comments string `xml:"comments,omitempty"`
	} `xml:"info,omitempty"`
	Questions []XMLQuestion `xml:"questions>question"`
}

type XMLQuestion struct {
	Price  string `xml:"price,attr"`
	Type   string `xml:"type,attr,omitempty"`
	Params struct {
		Params []XMLParam `xml:"param"`
	} `xml:"params"`
	Right struct {
		Answer string `xml:"answer"`
	} `xml:"right"`
}

type XMLNumberSet struct {
	Minimum string `xml:"minimum,attr"`
	Maximum string `xml:"maximum,attr"`
	Step    string `xml:"step,attr"`
}

type XMLParam struct {
	Name      string        `xml:"name,attr"`
	Type      string        `xml:"type,attr,omitempty"`
	Item      *XMLItem      `xml:"item,omitempty"`
	NumberSet *XMLNumberSet `xml:"numberSet,omitempty"`
	Content   string        `xml:",innerxml"`
}

type XMLItem struct {
	Type      string `xml:"type,attr"`
	IsRef     bool   `xml:"isRef,attr"`
	Placement string `xml:"placement,attr,omitempty"`
	Content   string `xml:",chardata"`
}
