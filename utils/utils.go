package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AdaptivePath 自适应相对目录和绝对目录处理 ${HOME} $HOME 变量
func AdaptivePath(path string) string {
	if strings.HasPrefix(path, "$HOME") {
		homeDir := os.Getenv("HOME")
		return strings.Replace(path, "$HOME", homeDir, 1)
	} else if strings.HasPrefix(path, "${HOME}") {
		homeDir := os.Getenv("HOME")
		return strings.Replace(path, "${HOME}", homeDir, 1)
	} else if filepath.IsAbs(path) {
		return path
	} else {
		// 处理相对路径
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Println("Error converting to absolute path:", err)
			return path
		}
		return absPath
	}
}

// GetGitRepoNameByGitAddress 从git地址 获取 git仓库名
func GetGitRepoNameByGitAddress(gitAddress string) string {
	// 使用 strings.Split 分割字符串
	parts := strings.Split(gitAddress, "/")
	// 获取最后一段
	repoName := parts[len(parts)-1]
	// 检查是否以 .git 结尾
	if strings.HasSuffix(repoName, ".git") {
		// 去掉 .git 后缀
		repoName = strings.TrimSuffix(repoName, ".git")
	}

	return repoName
}
