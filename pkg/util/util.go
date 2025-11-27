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

func BuildPostURL(baseUrl string, teamName string, postID string) string {
	base := strings.TrimSuffix(baseUrl, "/")
	if teamName == "" {
		return base + "/pl/" + postID
	}
	return base + "/" + strings.Trim(teamName, "/") + "/pl/" + postID
}
