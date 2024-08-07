package main

import (
	"fmt"
	"github.com/george012/git_sync/api"
	"github.com/george012/git_sync/common"
	"github.com/george012/git_sync/config"
	"github.com/george012/git_sync/custom_cmd"
	"github.com/george012/git_sync/repo_mg"
	"github.com/george012/gtbox"
	"github.com/george012/gtbox/gtbox_cmd"
	"github.com/george012/gtbox/gtbox_encryption"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/george012/gtbox/gtbox_sys"
	"os"
	"runtime"
	"time"
)

var (
	mRunMode       = ""
	mGitCommitHash = ""
	mGitCommitTime = ""
	mPackageOS     = ""
	mPackageTime   = ""
	mGoVersion     = ""
)

func setupApp() {
	runMode := gtbox.RunModeDebug
	switch mRunMode {
	case "debug":
		runMode = gtbox.RunModeDebug
	case "test":
		runMode = gtbox.RunModeTest
	case "release":
		runMode = gtbox.RunModeRelease
	default:
		runMode = gtbox.RunModeDebug
	}

	config.CurrentApp = config.NewApp(
		config.ProjectName,
		config.ProjectBundleID,
		fmt.Sprintf("%s service", config.ProjectName),
		runMode,
		config.APIPortDefault,
	)

	//	TODO 初始化gtbox及log分片
	if config.CurrentApp.CurrentRunMode == gtbox.RunModeDebug {
		cmdMap := map[string]string{
			"git_commit_hash": "git show -s --format=%H",
			"git_commit_time": "git show -s --format=\"%ci\" | cut -d ' ' -f 1,2 | sed 's/ /_/'",
			"build_os":        "go env GOOS",
			"go_version":      "go version | awk '{print $3}'",
		}
		cmdRes := gtbox_cmd.RunWith(cmdMap)

		if cmdRes != nil {
			mGitCommitHash = cmdRes["git_commit_hash"]
			mGitCommitTime = cmdRes["git_commit_time"]
			mPackageOS = cmdRes["build_os"]
			mGoVersion = cmdRes["go_version"]
			mPackageTime = time.Now().UTC().Format("2006-01-02_15:04:05")
		}
	}

	config.CurrentApp.GitCommitHash = mGitCommitHash
	config.CurrentApp.GitCommitTime = mGitCommitTime
	config.CurrentApp.GoVersion = mGoVersion
	config.CurrentApp.PackageOS = mPackageOS
	config.CurrentApp.PackageTime = mPackageTime

	//	TODO 处理自定义命令
	custom_cmd.HandleCustomCmds(os.Args, config.CurrentApp)

	gtbox.SetupGTBox(config.CurrentApp.AppName,
		config.CurrentApp.CurrentRunMode,
		config.CurrentApp.AppLogPath,
		30,
		gtbox_log.GTLogSaveHours,
		int(config.CurrentApp.HTTPRequestTimeOut.Seconds()),
	)

	en_str := gtbox_encryption.GTEnc("app starting...", "hello")
	gtbox_log.LogInfof(gtbox_encryption.GTDec(en_str, "hello"))

	hard_infos := gtbox_sys.GTGetHardInfo()
	snStr := fmt.Sprintf("%s|%s|%s|%s", hard_infos.CPUNumber, hard_infos.BaseBoardNumber, hard_infos.BiosNumber, hard_infos.DiskNumber)
	snStrEnc := gtbox_encryption.GTEnc(snStr, "sn")
	config.HardSN = snStrEnc
}

func main() {
	// 锁定当前的 goroutine 到操作系统线程
	runtime.LockOSThread()

	setupApp()

	config.SyncConfigFile(config.CurrentApp.AppConfigFilePath, nil)

	api.StartAPIService(&config.GlobalConfig.Api)

	repo_mg.StartAutoRepoSyncService(&config.GlobalConfig.RepoManagerConfig)

	common.LoadSigHandle(nil, nil)
}
