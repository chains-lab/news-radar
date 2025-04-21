package ape

import "fmt"

var (
	ErrorTagNameAlreadyTaken = fmt.Errorf("tag with this name already exists")
	ErrTagNotFound           = fmt.Errorf("tag not found")
	ErrAuthorNotFound        = fmt.Errorf("author not found")
	ErrArticleNotFound       = fmt.Errorf("article not found")
)
