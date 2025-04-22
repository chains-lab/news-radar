package ape

import "fmt"

var (
	ErrorTagNameAlreadyTaken = fmt.Errorf("tag with this name already exists")

	ErrTagNotFound    = fmt.Errorf("tag not found")
	ErrTagInactive    = fmt.Errorf("tag status inactive")
	ErrTooManyTags    = fmt.Errorf("too many tags")
	ErrTagReplication = fmt.Errorf("tag replication error")

	ErrAuthorNotFound    = fmt.Errorf("author not found")
	ErrAuthorInactive    = fmt.Errorf("author status inactive")
	ErrAuthorReplication = fmt.Errorf("author replication error")

	ErrArticleNotFound = fmt.Errorf("article not found")
)
