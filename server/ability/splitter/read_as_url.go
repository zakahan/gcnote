// -------------------------------------------------
// Package splitter
// Author: hanzhi
// Date: 2024/12/27
// -------------------------------------------------

package splitter

import (
	"fmt"
	"gcnote/server/config"
	"io"
	"os"
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

// ------------------------------------------------------- //

func ChunkReadReverse(chunks []string, urlPre, indexId, kbFileId string) (string, error) {
	// 字符串数组拼接返回
	for i, chunk := range chunks {
		if strings.HasPrefix(chunk, "!") {
			// 将url部分修改
			imageParts := ExtractImageURLParts(chunk)
			// 转为对应url
			// 参考样式`![史蒂芬大教堂](/test/data/1.png "这是一个教堂")`
			if imageParts[2] == "" {
				url, err := WebURL2LocalPath(imageParts[1], urlPre, indexId, kbFileId)
				if err != nil {
					return "", err
				}
				chunks[i] = fmt.Sprintf(`![%s](%s)`,
					imageParts[0], url)
			} else {
				url, err := WebURL2LocalPath(imageParts[1], urlPre, indexId, kbFileId)
				if err != nil {
					return "", err
				}
				chunks[i] = fmt.Sprintf(`![%s](%s "%s")`,
					imageParts[0], url, imageParts[2])
			}

		}
	}
	// 合并chunks并return
	return strings.Join(chunks, "\n\n"), nil
}

// WebURL2LocalPath 将Web URL转换为本地路径
func WebURL2LocalPath(path, urlPre, RealIndexId, RealKbFileId string) (string, error) {
	// 去除URL前缀
	urlPath := strings.TrimPrefix(path, urlPre)

	// 分割路径
	parts := strings.Split(urlPath, "/")
	var indexId string
	var kbFileId string
	// 提取index_id, kb_file_id, image_name
	if len(parts) >= 4 {
		indexId = parts[1]
		kbFileId = parts[2]
		imageName := parts[3]

		// 检查index_id和kb_file_id是否为空
		if indexId != "" && kbFileId != "" {
			return filepath.Join("images", imageName), nil
		} else {
			err := change(imageName, RealIndexId, RealKbFileId)
			if err != nil {
				return "", err
			}
			return filepath.Join("images", imageName), nil
		}
	} else {
		return "", fmt.Errorf("提取路径出现问题，URL路径设置有问题。")
	}
}

func change(imageName, RealIndexId, realKBFileId string) error {
	tmpImagePath := filepath.Join(config.PathCfg.KnowledgeBasePath, "images")
	// 将图片保存到指定目录
	//
	settingImageDir := filepath.Join(config.PathCfg.KnowledgeBasePath, RealIndexId, realKBFileId, "images")
	sourceImagePath := filepath.Join(tmpImagePath, imageName)

	// 确保目标目录存在
	if err := os.MkdirAll(settingImageDir, os.ModePerm); err != nil {
		return err
	}

	// 构建目标文件路径
	targetImagePath := filepath.Join(settingImageDir, imageName)

	// 打开源文件
	sourceFile, err := os.Open(sourceImagePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// 创建目标文件
	targetFile, err := os.Create(targetImagePath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	// 复制文件内容
	_, err = io.Copy(targetFile, sourceFile)
	if err != nil {
		return err
	}

	//// 删除源文件
	//if err := os.Remove(sourceImagePath); err != nil {
	//	return err
	//}

	return nil
}
