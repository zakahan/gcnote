// -------------------------------------------------
// Package wrench
// Author: hanzhi
// Date: 2024/12/16
// -------------------------------------------------

package wrench

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
)

// RemoveContents 清除目录里的东西
func RemoveContents(dir string) error {
	// 删除目录及其所有内容
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}

// CopyFile 复制单个文件
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			zap.S().Errorf("无法关闭文件：%v", err)
		}
	}(sourceFile)

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destinationFile *os.File) {
		err := destinationFile.Close()
		if err != nil {
			zap.S().Errorf("无法关闭文件：%v", err)
		}
	}(destinationFile)

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// 保留文件权限
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// CopyDir 递归复制目录
func CopyDir(srcDir, dstDir string) error {
	srcInfo, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", srcDir)
	}

	// 创建目标目录
	err = os.MkdirAll(dstDir, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			// 递归复制子目录
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// 复制文件
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
