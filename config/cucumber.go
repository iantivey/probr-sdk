package config

import "strings"

// CucumberTagsListToString will parse the tags specified in Vars.Tags
func CucumberTagsListToString(tags []string) string {
	var configTags []string
	for _, tag := range tags {
		configTags = append(configTags, "@"+tag)
	}
	return strings.Join(configTags, ",")
}
