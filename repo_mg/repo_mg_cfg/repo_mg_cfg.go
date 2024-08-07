package repo_mg_cfg

type RepoInfo struct {
	SSHKeyPath string `json:"ssh_key_path"`
	Address    string `json:"address"`
}
type RepoManagerConfig struct {
	SourceRepo  RepoInfo   `json:"source_repo"`
	TargetRepos []RepoInfo `json:"target_repos"`
}
