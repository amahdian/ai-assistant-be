package global

import "strings"

func TrimApiPrefix(url string) string {
	return strings.TrimPrefix(url, ApiPrefix)
}
