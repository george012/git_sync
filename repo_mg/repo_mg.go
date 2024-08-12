package repo_mg

import (
	"github.com/george012/git_sync/repo_mg/repo_mg_cfg"
	"github.com/george012/gtbox/gtbox_log"
	"time"
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
	go func() {
		for {
			for _, aRepo := range repoMgr.Config.Repos {
				err := gitCloneRepo(&aRepo)

				if err != nil {
					gtbox_log.LogErrorf("cone error src-repo[%s] targetretpo [%v]", aRepo.SourceRepo.Address, aRepo.TargetRepos)
					continue
				}

				err = gitPushRepo(&aRepo)

			}

			time.Sleep(30 * time.Second)
		}
	}()
}

func StartAutoRepoSyncService(cfg *repo_mg_cfg.RepoManagerConfig) {
	repoMgr := newRepoManager(cfg)
	repoMgr.startAutoSyncGitRepo()
	gtbox_log.LogInfof("Started AutoSync Git Repo")
}
