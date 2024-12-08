// -------------------------------------------------
// Package convert_docx
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package convert_docx

import "github.com/zakahan/docx2md"

func DocxConvert(documentPath string, outputDir string) (string, string, error) {
	return docx2md.DocxConvert(documentPath, outputDir)
}
