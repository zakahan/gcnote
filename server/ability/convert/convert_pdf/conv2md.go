// -------------------------------------------------
// Package convert_pdf
// Author: hanzhi
// Date: 2025/1/8
// -------------------------------------------------

package convert_pdf

import (
	"encoding/json"
	"fmt"
	"gcnote/server/config"
	"os"
	"os/exec"
	"path/filepath"
)

type PythonOutput struct {
	Success bool   `json:"success"`
	MdPath  string `json:"md_path"`
	MdDir   string `json:"md_dir"`
	Error   string `json:"error,omitempty"`
}

func PdfConvert(documentPath string, outputDir string) (string, string, error) {
	// Python 脚本路径
	pythonScript := filepath.Join(config.PathCfg.BaseProjectPath, "component", "pdf_convert_ability", "main.py")

	// 构建命令 与python环境参数
	python_exe := os.Getenv("GC_NOTE_PYTHON_EXE")
	cmd := exec.Command(python_exe, pythonScript, "--pdf_path", documentPath, "--output_dir", outputDir)

	// 获取标准输出
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing script: %s\n", err)
		return "", "", err
	}

	// 解析 JSON 输出
	var result PythonOutput
	err = json.Unmarshal(output, &result)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		return "", "", err
	}

	// 判断执行结果
	if result.Success {
		return result.MdPath, result.MdDir, nil
	} else {
		return "", "", fmt.Errorf(result.Error)
	}
}
