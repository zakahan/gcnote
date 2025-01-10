// -------------------------------------------------
// Package wrench
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package wrench

import "regexp"

// ValidateIndexName ？好像没啥用，因为index name应该是个uuid
func ValidateIndexName(indexName string) bool {
	// 不允许出现以下符号，这些符号可能被用来作为indexName和uuid的间隔的
	//re := regexp.MustCompile(`[^a-z0-9-_]`)
	re := regexp.MustCompile(`[?,"/\\*<>|]`)
	// 检查字符串中是否存在匹配的字符
	return !re.MatchString(indexName)
}

func ValidateKBName(indexName string) bool {
	// 不允许出现以下符号，这些符号可能被用来作为indexName和uuid的间隔的
	//re := regexp.MustCompile(`[^a-z0-9-_]`)
	re := regexp.MustCompile(`[?,"/\\*<>|]`)
	// 检查字符串中是否存在匹配的字符
	return !re.MatchString(indexName)
}
