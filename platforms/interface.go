package platforms

type Platform interface {
	CurrentCommit() string
	CurrentTag() string
	CurrentBranch() string
	CurrentBuildNumber() string
	IsDetected() bool
	Name() string
	BuildUrl() string
	RepositoryUrl() string
}
