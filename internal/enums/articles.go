package enums

type ArticleStatus string

const (
	ArticleStatusPending   ArticleStatus = "pending"
	ArticleStatusPublished ArticleStatus = "published"
	ArticleStatusInactive  ArticleStatus = "inactive"
)

func ParseArticleStatus(s string) (ArticleStatus, bool) {
	switch s {
	case "pending":
		return ArticleStatusPending, true
	case "published":
		return ArticleStatusPublished, true
	case "inactive":
		return ArticleStatusInactive, true
	default:
		return "", false
	}
}
