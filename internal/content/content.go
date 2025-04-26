package content

type Section struct {
	ID    int         `json:"id" bson:"_id"`
	Text  []TextBlock `json:"blocks,omitempty" bson:"blocks,omitempty"`
	Media []Media     `json:"media,omitempty" bson:"media,omitempty"`
	Audio []Audio     `json:"audio,omitempty" bson:"audio,omitempty"`
}

type Audio struct {
	URL      string `json:"url" bson:"url"`
	Duration int    `json:"duration" bson:"duration"`
	Caption  string `json:"caption" bson:"caption"`
	Icon     string `json:"icon" bson:"icon"`
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
	Text string `json:"text,omitempty" bson:"text,omitempty"`
	//TODO in future customize text block
}
