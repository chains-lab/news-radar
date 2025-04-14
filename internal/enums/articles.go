package enums

type ArticleStatus string

const (
	ArticleStatusPending  ArticleStatus = "pending"
	ArticleStatusActive   ArticleStatus = "active"
	ArticleStatusInactive ArticleStatus = "inactive"
)

func ParseArticleStatus(s string) (ArticleStatus, bool) {
	switch s {
	case "pending":
		return ArticleStatusPending, true
	case "active":
		return ArticleStatusActive, true
	case "inactive":
		return ArticleStatusInactive, true
	default:
		return "", false
	}
}
