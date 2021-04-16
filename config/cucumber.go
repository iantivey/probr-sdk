package config

import (
	"strings"
)

// CucumberTagsListToString will parse the tags specified in Vars.Tags
func CucumberTagsListToString(tags []string) string {
	var tagList []string
	for _, tag := range tags {
		tagList = append(tagList, "@"+tag)
	}
	return strings.Join(tagList, ",")
}

// Handle tag exclusions provided via the config vars file
func ConfigTagExclusionsListToString(tags []string) string {
	var tagList []string
	for _, tag := range tags {
		tagList = append(tagList, "~@"+tag)
	}
	return strings.Join(tagList, " && ")
}
