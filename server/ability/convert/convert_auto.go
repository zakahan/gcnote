// -------------------------------------------------
// Package convert
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package convert

import (
	"errors"
	"gcnote/server/ability/convert/convert_docx"
	"gcnote/server/ability/convert/convert_html"
	"gcnote/server/ability/convert/convert_md"
	"gcnote/server/ability/convert/convert_pdf"
	"gcnote/server/ability/convert/convert_txt"
	"os"
	"path/filepath"
)

func AutoConvert(documentPath string, outputDir string, suffix string) (string, string, error) {
	/*
		返回值：mdPath, mdString, err
	*/
	// 首先根据输入的suffix来处理
	var suffixList = []string{".docx", ".html", ".txt", ".md", ".pdf"}
	// 如果suffix不为空，直接匹配是否在suffixList里，否则需要判断
	if suffix == "" {
		//suffix :=
		suffix = filepath.Ext(documentPath)
	}
	// 查看是否允许这个suffix
	err := errors.New("不支持格式的文件格式：" + suffix +
		"，目前仅支持'.md', '.docx', '.html'和'.txt'类型的文件")
	for _, suf := range suffixList {

		if suf == suffix {
			err = nil
			break
		}
	}
	if err != nil {
		return "", "", err
	}
	if suffix == ".docx" {
		return convert_docx.DocxConvert(documentPath, outputDir)
	} else if suffix == ".pdf" {
		return convert_pdf.PdfConvert(documentPath, outputDir)
	} else if suffix == ".html" {
		// 在outputDir里创建一个名为images的文件夹
		err = os.Mkdir(filepath.Join(outputDir, "images"), os.ModePerm)
		if err != nil {
			return "", "", err
		}
		return convert_html.HtmlConvert(documentPath, outputDir)
	} else if suffix == ".txt" {
		err = os.Mkdir(filepath.Join(outputDir, "images"), os.ModePerm)
		if err != nil {
			return "", "", err
		}
		return convert_txt.TxtConvert(documentPath, outputDir)
	} else if suffix == ".md" {
		err = os.Mkdir(filepath.Join(outputDir, "images"), os.ModePerm)
		if err != nil {
			return "", "", err
		}
		return convert_md.MdConvert(documentPath, outputDir)
	} else {
		return "", "", errors.New("未知的文件格式：" + suffix)
	}

}
