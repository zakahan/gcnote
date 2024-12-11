// -------------------------------------------------
// Package wrench
// Author: hanzhi
// Date: 2024/12/11
// -------------------------------------------------

package wrench

import "regexp"

// ？好像没啥用，因为index name应该是个uuid
func validateIndexName(indexName string) bool {
	// 定义正则表达式，匹配小写字母、数字、连字符和下划线以外的字符
	re := regexp.MustCompile(`[^a-z0-9-_]`)
	// 检查字符串中是否存在匹配的字符
	return !re.MatchString(indexName)
}
