package content

import "fmt"

type SectionType string

type Section struct {
	Section SectionType    `json:"section" bson:"section"`
	Content map[string]any `json:"content" bson:"content"`
}

const (
	SectionTypeText     SectionType = "text"
	SectionTypeImage    SectionType = "image"
	SectionTypeVideo    SectionType = "video"
	SectionTypeLocation SectionType = "location"
	SectionTypeQuote    SectionType = "quote"
)

func SectionTypeParse(s string) (SectionType, error) {
	switch s {
	case "text":
		return SectionTypeText, nil
	case "image":
		return SectionTypeImage, nil
	case "video":
		return SectionTypeVideo, nil
	case "location":
		return SectionTypeLocation, nil
	case "quote":
		return SectionTypeQuote, nil
	default:
		return "", fmt.Errorf("invalid section type %s", s)
	}
}
