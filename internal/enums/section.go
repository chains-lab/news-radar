package enums

type SectionType string

const (
	SectionTypeText     SectionType = "text"
	SectionTypeAudio    SectionType = "audio"
	SectionTypeImage    SectionType = "image"
	SectionTypeVideo    SectionType = "video"
	SectionTypeLocation SectionType = "location"
	SectionTypeQuote    SectionType = "quote"
)

func SectionTypeParse(s string) (SectionType, bool) {
	switch s {
	case "text":
		return SectionTypeText, true
	case "image":
		return SectionTypeImage, true
	case "audio":
		return SectionTypeAudio, true
	case "video":
		return SectionTypeVideo, true
	case "location":
		return SectionTypeLocation, true
	case "quote":
		return SectionTypeQuote, true
	default:
		return "", false
	}
}

type TextMark string

const (
	MarkBold      TextMark = "bold"
	MarkItalic    TextMark = "italic"
	MarkUnderline TextMark = "underline"
	MarkStrike    TextMark = "strike"
	MarkUppercase TextMark = "uppercase"
)

func TextMarkParse(s string) (TextMark, bool) {
	switch s {
	case "bold":
		return MarkBold, true
	case "italic":
		return MarkItalic, true
	case "underline":
		return MarkUnderline, true
	case "strike":
		return MarkStrike, true
	case "uppercase":
		return MarkUppercase, true
	default:
		return "", false
	}
}
