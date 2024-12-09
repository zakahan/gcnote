package config

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
)

func getPackagePath(pkgName string) string {
	// 使用 go list 获取包信息
	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", pkgName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(out.String())
}

type PathConfig struct {
	BaseProjectPath string
	EtcConfigPath   string
}

func Constructor() PathConfig {
	baseProjectPath := getPackagePath("gcnote")
	etcConfigPath := filepath.Join(baseProjectPath, "server/etc/config.yaml")

	return PathConfig{
		BaseProjectPath: baseProjectPath,
		EtcConfigPath:   etcConfigPath,
	}
}

var PathCfg PathConfig = Constructor()
