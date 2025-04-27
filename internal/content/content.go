package content

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/news-radar/resources"
	"github.com/sirupsen/logrus"
)

type Section struct {
	ID    int               `json:"id" bson:"id"`
	Type  enums.SectionType `json:"type" bson:"type"`
	Text  []TextBlock       `json:"text,omitempty" bson:"text,omitempty"`
	Media []Media           `json:"media,omitempty" bson:"media,omitempty"`
	Audio []Audio           `json:"audio,omitempty" bson:"audio,omitempty"`
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
}

func ParseContentSection(req resources.Section) (Section, error) {
	sec := Section{
		ID: int(req.Id),
	}

	logrus.Info("1")

	sectionType, ok := enums.SectionTypeParse(req.Type)
	if !ok {
		return sec, validation.Errors{
			"type": validation.NewError("invalid", "invalid section type"),
		}
	}

	logrus.Info("2")

	if sectionType == enums.SectionTypeAudio {
		if req.Audio == nil || len(req.Audio) == 0 {
			return sec, validation.Errors{
				"audio": validation.NewError("invalid", "audio section must have at least one audio"),
			}
		}
		audio := make([]Audio, 0)
		for _, a := range req.Audio {
			audio = append(audio, Audio{
				URL:      a.Url,
				Duration: int(a.Duration),
				Caption:  a.Caption,
				Icon:     a.Icon,
			})
		}

		sec.Type = enums.SectionTypeAudio
		sec.Audio = audio
	}

	logrus.Info("3")

	if sectionType == enums.SectionTypeText {
		if req.Text == nil || len(req.Text) == 0 {
			return sec, validation.Errors{
				"text": validation.NewError("invalid", "text section must have at least one text"),
			}
		}
		text := make([]TextBlock, 0)
		for _, t := range req.Text {
			text = append(text, TextBlock{
				Text: t.Text,
			})
		}
		sec.Type = enums.SectionTypeText
		sec.Text = text
	}

	logrus.Info("4")

	if sectionType == enums.SectionTypeMedia {
		if req.Media == nil || len(req.Media) == 0 {
			return sec, validation.Errors{
				"media": validation.NewError("invalid", "media section must have at least one media"),
			}
		}
		media := make([]Media, 0)
		for _, m := range req.Media {
			media = append(media, Media{
				URL:     m.Url,
				Caption: m.Caption,
				Alt:     m.Alt,
				Width:   int(m.Width),
				Height:  int(m.Height),
				Source:  m.Source,
			})
		}
		sec.Type = enums.SectionTypeMedia
		sec.Media = media
	}

	logrus.Info("5")

	return sec, nil
}
