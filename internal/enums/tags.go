package enums

type TagStatus string

const (
	TagStatusActive   TagStatus = "active"
	TagStatusInactive TagStatus = "inactive"
)

func ParseTagStatus(s string) (TagStatus, bool) {
	switch s {
	case "active":
		return TagStatusActive, true
	case "inactive":
		return TagStatusInactive, true
	default:
		return "", false
	}
}

type TagType string

const (
	TagTypeTopic   TagType = "topic"
	TagTypeDefault TagType = "default"
)

func ParseTagType(s string) (TagType, bool) {
	switch s {
	case "topic":
		return TagTypeTopic, true
	case "default":
		return TagTypeDefault, true
	default:
		return "", false
	}
}
