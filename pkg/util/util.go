package util

import (
	"regexp"
	"strings"
)

var zoomURLRegex = regexp.MustCompile(`https?://[^\s]*zoom\.(us|com|gov)/[^\s]*`)

func HasZoomURL(message string) bool {
	return zoomURLRegex.MatchString(message)
}

func BuildPostURL(baseUrl string, postID string) string {
	return strings.TrimSuffix(baseUrl, "/") + "/pl/" + postID
}
