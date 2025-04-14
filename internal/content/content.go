package content

import "github.com/hs-zavet/news-radar/internal/enums"

type Section struct {
	ID    string            `json:"id" bson:"_id"`
	Type  enums.SectionType `json:"type" bson:"type"`
	Media *Media            `json:"media,omitempty" bson:"media,omitempty"`
	Text  []TextBlock       `json:"blocks,omitempty" bson:"blocks,omitempty"`
}

type Media struct {
	URL     string `json:"url" bson:"url"`
	Caption string `json:"caption" bson:"caption"`
	Alt     string `json:"alt" bson:"alt"`
	Width   int    `json:"width" bson:"width"`
	Height  int    `json:"height" bson:"height"`
	Source  string `json:"source" bson:"source"`
}

type TextBlock struct {
	Text  *string          `json:"text,omitempty" bson:"text,omitempty"`
	Marks []enums.TextMark `json:"marks,omitempty" bson:"marks,omitempty"`
	Color *string          `json:"color,omitempty" bson:"color,omitempty"`
	Link  *string          `json:"link,omitempty" bson:"link,omitempty"`
}
