package requests

import (
	"fmt"

	"github.com/hs-zavet/news-radar/internal/content"
	"github.com/hs-zavet/news-radar/internal/enums"
	"github.com/hs-zavet/news-radar/resources"
)

func UpdateSection(res resources.Content) (content.Section, error) {
	secType, ok := enums.SectionTypeParse(res.Type)
	if !ok {
		return content.Section{}, fmt.Errorf("invalid section type: %q", res.Type)
	}

	sec := content.Section{
		ID:   res.Id,
		Type: secType,
	}

	// Маппинг медиа (image/video)
	if res.Media != nil {
		sec.Media = &content.Media{
			URL:     res.Media.Url,
			Caption: res.Media.Caption,
			Alt:     res.Media.Alt,
			Width:   int(res.Media.Width),
			Height:  int(res.Media.Height),
			Source:  res.Media.Source,
		}
	}

	// Маппинг аудио
	if res.Audio != nil {
		sec.Audio = &content.Audio{
			URL:      res.Audio.Url,
			Duration: int(res.Audio.Duration),
			Caption:  res.Audio.Caption,
			Icon:     res.Audio.Icon,
		}
	}

	// Маппинг текстовых блоков
	for _, ti := range res.Text {
		var marks []enums.TextMark
		for _, m := range ti.Marks {
			mark, ok := enums.TextMarkParse(m)
			if !ok {
				return content.Section{}, fmt.Errorf("invalid text mark: %q", m)
			}
			marks = append(marks, mark)
		}

		sec.Text = append(sec.Text, content.TextBlock{
			Text:  ti.Text,
			Marks: marks,
			Color: ti.Color,
			Link:  ti.Link,
		})
	}

	return sec, nil
}
