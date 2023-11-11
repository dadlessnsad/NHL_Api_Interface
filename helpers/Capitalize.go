package helpers

import (
	"strings"
)

func Capitalize(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str)
}
