package common

import (
	"regexp"
	"strings"
)

func Slug(str string) string {
	str = strings.ToLower(strings.TrimSpace(str))
	str = regexp.MustCompile("[^a-z0-9-]").ReplaceAllString(str, "-")
	return str
}

func PanicOn(err error) {
	if err != nil {
		panic(err.Error())
	}
}
