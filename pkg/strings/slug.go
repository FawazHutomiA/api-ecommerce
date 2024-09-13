package strings

import (
	"regexp"
	"strings"
)

func Slug(s string) string {
	// Convert the string to lowercase
	slug := strings.ToLower(s)

	// Replace non-alphanumeric characters with dashes
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing dashes
	slug = strings.Trim(slug, "-")
	slug = strings.Trim(slug, "#")

	return slug
}
