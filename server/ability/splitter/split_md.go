// -------------------------------------------------
// Package splitter
// Author: hanzhi
// Date: 2024/12/20
// -------------------------------------------------

package splitter

import (
	"regexp"
	"strings"
)

func splitMarkdown(input string, maxTextLength int) []string {
	// Regular expressions for splitting
	patterns := []string{
		`(?m)^(#{1,6} .*)`,                 // Heading: #, ##, ###, etc.
		`(?m)^\|.*\|\s*\n(?:\|.*\|\s*\n)*`, // Tables
		`!\[.*?\]\(.*?\)`,                  // Images
		`(?m)^>.*`,                         // Blockquotes
		"```[\\s\\S]*?```",                 // Code blocks
		`(?m)^-{3,}\s*$`,                   // Horizontal rules
	}

	// Combine patterns into one
	re := regexp.MustCompile(strings.Join(patterns, "|"))

	// Find matches and split
	matches := re.FindAllStringIndex(input, -1)

	var result []string
	lastIndex := 0

	for _, match := range matches {
		start, end := match[0], match[1]
		if start > lastIndex {
			result = append(result, strings.TrimSpace(input[lastIndex:start]))
		}
		result = append(result, strings.TrimSpace(input[start:end]))
		lastIndex = end
	}

	// Append the remaining part of the input
	if lastIndex < len(input) {
		result = append(result, strings.TrimSpace(input[lastIndex:]))
	}

	// Filter out empty strings
	finalResult := []string{}
	for _, item := range result {
		if item != "" {
			finalResult = append(finalResult, item)
		}
	}

	// Split large text segments
	var processedResult []string
	for _, segment := range finalResult {
		if isTextSegment(segment) && len(segment) > maxTextLength {
			processedResult = append(processedResult, splitLargeText(segment, maxTextLength)...)
		} else {
			processedResult = append(processedResult, segment)
		}
	}

	return processedResult
}

func isTextSegment(segment string) bool {
	// Check if the segment is plain text (not table, code block, etc.)
	return !strings.HasPrefix(segment, "|") && !strings.HasPrefix(segment, "```") && !strings.HasPrefix(segment, ">") && !strings.HasPrefix(segment, "#") && !strings.HasPrefix(segment, "!") && !strings.HasPrefix(segment, "---")
}

func splitLargeText(text string, maxTextLength int) []string {
	var result []string
	paragraphs := strings.Split(text, "\n\n")
	current := ""

	for _, paragraph := range paragraphs {
		if len(current)+len(paragraph) > maxTextLength {
			result = append(result, strings.TrimSpace(current))
			current = paragraph
		} else {
			if current != "" {
				current += "\n\n"
			}
			current += paragraph
		}
	}

	if current != "" {
		result = append(result, strings.TrimSpace(current))
	}

	return result
}
