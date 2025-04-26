package requests

import (
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/resources"
)

func UpdateSection(section resources.Section) (content.Section, error) {
	sec := content.Section{
		ID: int(section.Id),
	}

	if section.Audio != nil {
		audio := make([]resources.SectionAudioInner, 0)
		for _, a := range section.Audio {
			audio = append(audio, resources.SectionAudioInner{
				Url:      a.Url,
				Duration: a.Duration,
				Caption:  a.Caption,
				Icon:     a.Icon,
			})
		}
	}

	if section.Text != nil {
		text := make([]resources.SectionTextInner, 0)
		for _, t := range section.Text {
			text = append(text, resources.SectionTextInner{
				Text: t.Text,
			})
		}
		section.Text = text
	}

	for _, ti := range section.Text {
		sec.Text = append(sec.Text, content.TextBlock{
			Text: *ti.Text,
		})
	}

	for _, m := range section.Media {
		sec.Media = append(sec.Media, content.Media{
			URL:     m.Url,
			Caption: m.Caption,
			Alt:     m.Alt,
			Width:   int(m.Width),
			Height:  int(m.Height),
			Source:  m.Source,
		})
	}

	return sec, nil
}
