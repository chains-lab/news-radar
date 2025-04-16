package responses

import (
	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/resources"
)

func Content(section content.Section) resources.Content {
	res := resources.Content{
		Id:   section.ID,
		Type: string(section.Type),
	}

	if section.Media != nil {
		res.Media = &resources.ContentMedia{
			Url:     section.Media.URL,
			Caption: section.Media.Caption,
			Alt:     section.Media.Alt,
			Width:   int32(section.Media.Width),
			Height:  int32(section.Media.Height),
			Source:  section.Media.Source,
		}
	}

	if section.Audio != nil {
		res.Audio = &resources.ContentAudio{
			Url:      section.Audio.URL,
			Caption:  section.Audio.Caption,
			Duration: int32(section.Audio.Duration),
			Icon:     section.Audio.Icon,
		}
	}

	if section.Text != nil {
		text := make([]resources.ContentTextInner, 0)
		for _, t := range section.Text {
			if t.Text != nil {
				marks := make([]string, 0)
				for _, mark := range t.Marks {
					marks = append(marks, string(mark))
				}

				text = append(text, resources.ContentTextInner{
					Text:  t.Text,
					Color: t.Color,
					Link:  t.Link,
					Marks: marks,
				})
			}
		}
		res.Text = text
	}

	return res
}
