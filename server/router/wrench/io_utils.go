// -------------------------------------------------
// Package wrench
// Author: hanzhi
// Date: 2024/12/16
// -------------------------------------------------

package wrench

import "os"

// RemoveContents 清除目录里的东西
func RemoveContents(dir string) error {
	// 删除目录及其所有内容
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}
