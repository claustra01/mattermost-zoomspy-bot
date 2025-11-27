package util

import (
	"regexp"
	"strings"
)

var zoomURLRegex = regexp.MustCompile(`https?://[^\s]*zoom\.(us|com|gov)/[^\s]*`)
var zoomArchiveRegex = regexp.MustCompile(`https?://[^\s]*zoom\.(us|com|gov)/(rec|recording)/[^\s]*`)

func HasZoomURL(message string) bool {
	return zoomURLRegex.MatchString(message) && !zoomArchiveRegex.MatchString(message)
}

func BuildPostURL(baseUrl string, postID string) string {
	return strings.TrimSuffix(baseUrl, "/") + "/pl/" + postID
}
