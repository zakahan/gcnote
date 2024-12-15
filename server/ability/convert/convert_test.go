// -------------------------------------------------
// Package convert
// Author: hanzhi
// Date: 2024/12/15
// -------------------------------------------------

package convert

import (
	"fmt"
	"gcnote/server/config"
	"path/filepath"
	"testing"
)

func TestAutoConvert(t *testing.T) {
	path := filepath.Join(config.PathCfg.BaseProjectPath, "test/docs/23年统计公报-节选.docx")
	tmpPath := filepath.Join(config.PathCfg.TempDirPath)
	s, _, err := AutoConvert(path, tmpPath, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}
