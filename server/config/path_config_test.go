// -------------------------------------------------
// Package config
// Author: hanzhi
// Date: 2024/12/9
// -------------------------------------------------

package config

import (
	"fmt"
	"testing"
)

func TestPathCfg(t *testing.T) {
	fmt.Println(PathCfg.EtcConfigPath)
	fmt.Println(PathCfg.BaseProjectPath)
}
