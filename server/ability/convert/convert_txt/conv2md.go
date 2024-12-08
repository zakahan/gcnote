// -------------------------------------------------
// Package convert_txt
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package convert_txt

import (
	"fmt"
	"github.com/zakahan/docx2md/docx_parser"
	"io"
	"log"
	"os"
)

func TxtConvert(documentPath string, outputDir string) (string, string, error) {
	mdPath, _, err := docx_parser.CreateMdDir(documentPath, outputDir, ".txt")
	if err != nil {
		fmt.Println("Error:", err)
		return "", "", err
	}
	// 打开文件
	file, err := os.Open(documentPath)
	if err != nil {
		log.Fatalf("无法打开文件：%v", err)
		return "", "", err
	}
	defer file.Close()
	// 读取文件
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("读取文件时发生错误")
		return "", "", err
	}
	inputString := string(data)
	// 保存
	err = docx_parser.SaveFile(mdPath, inputString)
	return mdPath, inputString, err
}
