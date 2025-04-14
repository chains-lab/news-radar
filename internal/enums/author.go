package enums

type AuthorStatus string

const (
	AuthorStatusActive   AuthorStatus = "active"
	AuthorStatusInactive AuthorStatus = "inactive"
)

func ParseAuthorStatus(s string) (AuthorStatus, bool) {
	switch s {
	case "active":
		return AuthorStatusActive, true
	case "inactive":
		return AuthorStatusInactive, true
	default:
		return "", false
	}
}
