package responses

import (
	"github.com/chains-lab/news-radar/internal/content"
	"github.com/chains-lab/news-radar/internal/enums"
	"github.com/chains-lab/news-radar/resources"
)

func ContentSection(section content.Section) resources.Section {
	res := resources.Section{
		Id: int32(section.ID),
	}

	if section.Media != nil {
		media := make([]resources.SectionMediaInner, 0)
		for _, m := range section.Media {
			media = append(media, resources.SectionMediaInner{
				Url:     m.URL,
				Caption: m.Caption,
				Alt:     m.Alt,
				Width:   int32(m.Width),
				Height:  int32(m.Height),
				Source:  m.Source,
			})
		}
		res.Media = media
		res.Type = string(enums.SectionTypeMedia)
	}

	if section.Audio != nil {
		audio := make([]resources.SectionAudioInner, 0)
		for _, a := range section.Audio {
			audio = append(audio, resources.SectionAudioInner{
				Url:      a.URL,
				Duration: int32(a.Duration),
				Caption:  a.Caption,
				Icon:     a.Icon,
			})
		}
		res.Audio = audio
		res.Type = string(enums.SectionTypeAudio)
	}

	if section.Text != nil {
		text := make([]resources.SectionTextInner, 0)
		for _, t := range section.Text {
			text = append(text, resources.SectionTextInner{
				Text: t.Text,
			})
		}
		res.Text = text
		res.Type = string(enums.SectionTypeText)
	}

	return res
}
