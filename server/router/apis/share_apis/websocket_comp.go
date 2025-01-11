// -------------------------------------------------
// Package share_apis
// Author: hanzhi
// Date: 2025/1/10
// -------------------------------------------------

package share_apis

import (
	"bytes"
	"fmt"
	"gcnote/server/ability/splitter"
	"gcnote/server/config"
	"gcnote/server/model"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

//readInitialContent reads the initial content for a room

func readInitialContent(shareFileId string) (string, string, error) {
	// 检查当前shareFileId是否存在
	if shareFileId == "" {
		return "", "", fmt.Errorf("share file id is empty")
	}

	// sql查询
	shareFile := model.ShareFile{}
	tx := config.DB.Where("share_file_id = ?", shareFileId).First(&shareFile)
	if tx.Error != nil {
		return "", "", fmt.Errorf("share file not found")
	}

	fileDir := filepath.Join(config.PathCfg.ShareFileDirPath, shareFileId)
	filePath := filepath.Join(fileDir, shareFile.FileName+".md")

	// 首先检查fileDir和filePath是否存在
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		zap.S().Debugf("fileDir: %s is not exist", fileDir)
		return "", "", fmt.Errorf("file directory not found")
	}

	// 读取文件，转为字符串
	content, err := os.ReadFile(filePath)
	if err != nil {
		zap.S().Errorf("Failed to read file: %v", err)
		return "", "", err
	}

	// 将content转为切片然后读取
	chunks := splitter.SplitMarkdownEasy(string(content))
	resultData := splitter.ChunkRead(chunks, config.PathCfg.ImageServerURL, "share", shareFileId)

	return resultData, shareFile.FileName, nil
}

const (
	ServerID uint32 = 0 // 服务器ID设置为0
)

func initializeYDoc(content string) []byte {
	buf := new(bytes.Buffer)
	//
	writeVarUint(buf, messageSync)
	// 写入消息类型 (messageYjsUpdate = 2)
	writeVarUint(buf, 2)

	// 构建内容部分
	contentBuf := new(bytes.Buffer)

	// 版本信息
	contentBuf.Write([]byte{1, 1})
	// 时间戳/client ID
	// 就设置这个值是server id了.... 这其实是来自某次的测试日志，我不知道怎么改了
	// 我试着改成全1或者全0都不行，所以我不改了........
	contentBuf.Write([]byte{196, 248, 225, 213})
	// 结构标识符
	contentBuf.Write([]byte{10, 0})
	// 内容类型标识
	contentBuf.Write([]byte{4, 1})
	// shared-text 长度
	contentBuf.WriteByte(11)
	// shared-text 字符串
	contentBuf.WriteString("shared-text")
	// 内容长度
	writeVarUint(contentBuf, uint64(len(content)))
	// 实际内容
	contentBuf.WriteString(content)
	// 结束标记
	contentBuf.Write([]byte{10, 0})

	// 写入内容长度
	writeVarUint(buf, uint64(contentBuf.Len()))

	// 写入内容
	buf.Write(contentBuf.Bytes())

	return buf.Bytes()
}
