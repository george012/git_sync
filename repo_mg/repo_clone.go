package repo_mg

import (
	"errors"
	"fmt"
	"github.com/george012/git_sync/config"
	"github.com/george012/git_sync/repo_mg/repo_mg_cfg"
	"github.com/george012/git_sync/utils"
	"github.com/george012/gtbox/gtbox_files"
	"github.com/george012/gtbox/gtbox_log"
	"os"
	"os/exec"
	"sync"
)

var sshAgentLock sync.Mutex

func checkSSHKey(aRepo *repo_mg_cfg.RepoInfo) error {
	tmpSSHKeyPath := utils.AdaptivePath(aRepo.SourceRepo.SSHKeyPath)
	_, err := gtbox_files.GTToolsFileRead(tmpSSHKeyPath)
	return err
}

func getRepoDir(repoAddr string) string {
	// 仓库名
	repoName := utils.GetGitRepoNameByGitAddress(repoAddr)

	repoDir := fmt.Sprintf("%s/%s", config.CurrentApp.ReposDir, repoName)

	return repoDir
}

// getCurrentBranch 获取当前分支
func getCurrentBranch(repoDir string) string {
	getBranchCmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	getBranchCmd.Dir = repoDir // 指定工作目录
	branchNameBytes, err := getBranchCmd.Output()
	if err != nil {
		gtbox_log.LogErrorf("failed to get current branch: %v", err)
		return ""
	}
	branchName := string(branchNameBytes)
	return branchName[:len(branchName)-1] // 去掉结尾的换行符
}

func gitCloneRepo(aRepo *repo_mg_cfg.RepoInfo) error {
	sshAgentLock.Lock()
	defer sshAgentLock.Unlock()

	err := checkSSHKey(aRepo)
	if err != nil {
		return errors.New(fmt.Sprintf("read ssh-key err [%s]", err.Error()))
	}

	// 启动 ssh-agent
	agentCmd := exec.Command("ssh-agent")
	if err := agentCmd.Run(); err != nil {
		return fmt.Errorf("failed to start ssh-agent: %v", err)
	}
	defer exec.Command("ssh-agent", "-k").Run()

	// 添加指定的 SSH 私钥
	addCmd := exec.Command("ssh-add", aRepo.SourceRepo.SSHKeyPath)
	if err = addCmd.Run(); err != nil {
		if err.Error() != "exit status 1" {
			return fmt.Errorf("failed to add ssh key: %v", err)
		}
	}

	// 在操作完成后，删除指定的 SSH 私钥
	defer func() {
		delCmd := exec.Command("ssh-add", "-d", aRepo.SourceRepo.SSHKeyPath)
		if err = delCmd.Run(); err != nil {
			fmt.Printf("failed to delete ssh key: %v\n", err)
		}
	}()

	repoDir := getRepoDir(aRepo.SourceRepo.Address)

	// 执行 git clone 命令
	cmd := exec.Command("git", "clone", aRepo.SourceRepo.Address, repoDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed with error: %v", err)
	}

	return nil
}
