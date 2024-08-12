package repo_mg

import (
	"fmt"
	"github.com/george012/git_sync/repo_mg/repo_mg_cfg"
	"github.com/george012/gtbox/gtbox_log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func gitPushRepo(aRepo *repo_mg_cfg.RepoInfo) error {
	repoDir := getRepoDir(aRepo.SourceRepo.Address)

	var wg sync.WaitGroup

	for i, tRepo := range aRepo.TargetRepos {
		wg.Add(1)
		go func(i int, tgRepo repo_mg_cfg.RepoInfoBase) {
			defer wg.Done()

			newOriginName := fmt.Sprintf("new-origin%d", i)

			// 添加新的远程地址
			addRemoteCmd := exec.Command("git", "remote", "add", newOriginName, tgRepo.Address)
			addRemoteCmd.Dir = repoDir // 指定工作目录
			addRemoteCmd.Stdout = os.Stdout
			addRemoteCmd.Stderr = os.Stderr

			if err := addRemoteCmd.Run(); err != nil {
				gtbox_log.LogErrorf("failed to add remote to repo: [%s] [%s]", err.Error(), tgRepo.Address)
				return
			}

			currentBranch := getCurrentBranch(repoDir)

			if currentBranch == "" {
				return
			}

			// 获取远程仓库的最新提交哈希
			remoteHash := getRemoteLatestCommitHash(newOriginName, currentBranch, repoDir)
			if remoteHash == "" {
				gtbox_log.LogErrorf("failed to get remote commit hash from [%s]", tgRepo.Address)
				return
			}

			// 获取本地仓库的最新提交哈希
			localHash := getLocalLatestCommitHash(repoDir)
			if localHash == "" {
				gtbox_log.LogErrorf("failed to get local commit hash from [%s]", repoDir)
				return
			}

			// 比较本地和远程的哈希
			if localHash == remoteHash {
				gtbox_log.LogInfof("The remote repository [%s] is already up-to-date. Skipping push.", tgRepo.Address)
				return
			}

			// 推送代码到新的远程仓库
			pushCmd := exec.Command("git", "push", newOriginName, currentBranch)
			pushCmd.Dir = repoDir // 指定工作目录
			pushCmd.Stdout = os.Stdout
			pushCmd.Stderr = os.Stderr

			if err := pushCmd.Run(); err != nil {
				gtbox_log.LogErrorf("failed to push to new remote: [%s] [%v]", tgRepo.Address, err)
			}
		}(i, tRepo)
	}

	wg.Wait()

	// 返回 nil，因为不再需要处理单个目标失败的情况
	return nil
}

func getRemoteLatestCommitHash(remoteName, branchName, repoDir string) string {
	// 获取远程仓库最新提交的哈希
	cmd := exec.Command("git", "ls-remote", remoteName, branchName)
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		gtbox_log.LogErrorf("failed to get remote commit hash: %v", err)
		return ""
	}

	// 解析输出，提取哈希值
	hash := strings.Fields(string(output))
	if len(hash) > 0 {
		return hash[0]
	}
	return ""
}

func getLocalLatestCommitHash(repoDir string) string {
	// 获取本地仓库最新提交的哈希
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		gtbox_log.LogErrorf("failed to get local commit hash: %v", err)
		return ""
	}
	return strings.TrimSpace(string(output))
}
