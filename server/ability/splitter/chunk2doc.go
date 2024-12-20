// -------------------------------------------------
// Package splitter
// Author: hanzhi
// Date: 2024/12/20
// -------------------------------------------------

package splitter

import (
	"gcnote/server/ability/document"
	"strconv"
	"strings"
)

func Chunk2Doc(chunks []string, kbFileId, indexId string) []*document.Document {
	/*
	   暂时先不考虑图片的问题了，图片就放在那了，不加理解之类的东西了。
	*/
	l := len(chunks)
	var docList []*document.Document = make([]*document.Document, l)
	for i, chunk := range chunks {
		var imagePath = ""
		if strings.HasPrefix(chunk, "!") {
			imagePath = ExtractImageURL(chunk)
		}
		docList[i] = &document.Document{
			PageContent: chunk,
			Metadata: map[string]string{
				"doc_id":     strconv.Itoa(i), // 自增更合适
				"kb_file_id": kbFileId,
				"index_id":   indexId,
				"type":       getType(chunk),
				"image_path": imagePath,
			},
		}
	}
	return docList
}

func getType(chunk string) string {
	switch {
	case strings.HasPrefix(chunk, "|"):
		return document.TABLE.String()
	case strings.HasPrefix(chunk, "!"):
		return document.IMAGE.String()
	default:
		// 如果没有匹配的情况发生，可以返回一个默认值或空字符串
		return document.TEXT.String()
	}
}
