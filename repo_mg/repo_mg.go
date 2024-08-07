package repo_mg

import (
	"github.com/george012/git_sync/repo_mg/repo_mg_cfg"
	"github.com/george012/gtbox/gtbox_log"
)

type RepoManager struct {
	Config *repo_mg_cfg.RepoManagerConfig `json:"config"`
}

func newRepoManager(cfg *repo_mg_cfg.RepoManagerConfig) *RepoManager {
	repMgr := &RepoManager{
		Config: cfg,
	}
	return repMgr
}

func (repoMgr *RepoManager) startAutoSyncGitRepo() {

}

func StartAutoRepoSyncService(cfg *repo_mg_cfg.RepoManagerConfig) {
	go func() {
		repoMgr := newRepoManager(cfg)
		repoMgr.startAutoSyncGitRepo()
		gtbox_log.LogInfof("Started AutoSync Git Repo")
	}()
}
