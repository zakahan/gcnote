// -------------------------------------------------
// Package splitter
// Author: hanzhi
// Date: 2024/12/27
// -------------------------------------------------

package splitter

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ChunkRead(chunks []string, urlPre, indexId, kbFileId string) string {
	// 字符串数组拼接返回
	for i, chunk := range chunks {
		if strings.HasPrefix(chunk, "!") {
			// 将url部分修改
			imageParts := ExtractImageURLParts(chunk)
			// 转为对应url
			// 参考样式`![史蒂芬大教堂](/test/data/1.png "这是一个教堂")`
			if imageParts[2] == "" {
				chunks[i] = fmt.Sprintf(`![%s](%s)`,
					imageParts[0], LocalPath2WebURL(imageParts[1], urlPre, indexId, kbFileId))
			} else {
				chunks[i] = fmt.Sprintf(`![%s](%s "%s")`,
					imageParts[0], LocalPath2WebURL(imageParts[1], urlPre, indexId, kbFileId), imageParts[2])
			}

		}
	}
	// 合并chunks并return
	return strings.Join(chunks, "\n\n")
}

func LocalPath2WebURL(path, urlPre, indexId, kbFileId string) string {
	// 大概就是将类似 /test/data/1.png 或者images\\image10.png这样的处理为
	//这样的 http://127.0.0.1:8086/images/514b6721-26f9-46da-9be1-0b92261d2290/d6d6206b-f200-4715-86b4-b8a512fc401e/image2.png
	// url切分，提取后半部分
	cleanPath := filepath.Clean(path)
	filename := filepath.Base(cleanPath)

	url := fmt.Sprintf("%s/%s/%s/%s", urlPre, indexId, kbFileId, filename)
	return url
}
