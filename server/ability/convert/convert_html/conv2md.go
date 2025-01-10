// -------------------------------------------------
// Package convert_html
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package convert_html

import (
	"fmt"
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/zakahan/docx2md/docx_parser"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
)

func HtmlConvert(documentPath string, outputDir string) (string, string, error) {
	mdPath, _, err := docx_parser.CreateMdDir(documentPath, outputDir, ".html")
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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			zap.S().Errorf("无法关闭文件：%v", err)
		}
	}(file)
	// 读取文件
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("读取文件时发生错误")
	}

	// 将数据转为字符串
	inputString := string(data)
	// 改一下，添加表格处理的功能
	markdownStr, err := htmltomarkdown.ConvertString(inputString)
	if err != nil {
		log.Fatal(err)
	}

	err = docx_parser.SaveFile(mdPath, markdownStr)
	return mdPath, markdownStr, err
}
