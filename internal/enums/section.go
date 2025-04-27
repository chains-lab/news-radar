package enums

type SectionType string

const (
	SectionTypeText  SectionType = "text"
	SectionTypeMedia SectionType = "media"
	SectionTypeAudio SectionType = "audio"
)

func SectionTypeParse(s string) (SectionType, bool) {
	switch s {
	case "text":
		return SectionTypeText, true
	case "media":
		return SectionTypeMedia, true
	case "audio":
		return SectionTypeAudio, true
	default:
		return "", false
	}
}
