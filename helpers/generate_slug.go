package helpers

import (
	"strings"

	"github.com/gosimple/slug"
)

func GenerateSlug(name string) string {
	return slug.Make(strings.ToLower(name))
}
