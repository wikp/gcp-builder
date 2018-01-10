package containers

type inspect []inspectedImage

type inspectedImage struct {
	RepoDigests []string `json:"RepoDigests"`
}
