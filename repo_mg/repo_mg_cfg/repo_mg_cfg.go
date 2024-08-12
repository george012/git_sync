package repo_mg_cfg

type RepoInfoBase struct {
	SSHKeyPath string `json:"ssh_key_path"`
	Address    string `json:"address"`
}

type RepoInfo struct {
	SourceRepo  RepoInfoBase   `json:"source_repo"`
	TargetRepos []RepoInfoBase `json:"target_repos"`
}

type RepoManagerConfig struct {
	Repos []RepoInfo `json:"repos"`
}
