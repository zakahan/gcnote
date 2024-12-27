// -------------------------------------------------
// Package splitter
// Author: hanzhi
// Date: 2024/12/20
// -------------------------------------------------

package splitter

import "regexp"

// ExtractImageURL extracts the URL from a Markdown image tag (![]())
func ExtractImageURL(markdown string) string {
	// Regular expression to match ![]() style Markdown image syntax
	re := regexp.MustCompile(`!\[(.*?)]\((.*?)(?: "(.*?)")?\)`)
	matches := re.FindAllStringSubmatch(markdown, -1)

	if matches != nil && len(matches) > 0 && matches[0] != nil && len(matches[0]) > 2 {
		return matches[0][2]
	} else {
		return ""
	}
}

func ExtractImageURLParts(markdown string) []string {
	// Regular expression to match ![]() style Markdown image syntax
	re := regexp.MustCompile(`!\[(.*?)]\((.*?)(?: "(.*?)")?\)`)
	matches := re.FindAllStringSubmatch(markdown, -1)

	if matches != nil && len(matches) > 0 && matches[0] != nil && len(matches[0]) > 3 {
		return []string{matches[0][1], matches[0][2], matches[0][3]}
	} else if matches != nil && len(matches) > 0 && matches[0] != nil && len(matches[0]) > 2 {
		return []string{matches[0][1], matches[0][2], ""}
	} else {
		return []string{"", "", ""}
	}
}
