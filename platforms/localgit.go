package platforms

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"strings"
)

type LocalGitRepositoryPlatform struct {
	repo *git.Repository
}

func NewLocalGitRepositoryPlatform() (*LocalGitRepositoryPlatform, error) {
	repository, err := git.PlainOpen(".")
	if err != nil {
		fmt.Printf("Error is :%s", err)
		return nil, err
	}

	return &LocalGitRepositoryPlatform{
		repository,
	}, nil
}

func (l LocalGitRepositoryPlatform) IsDetected() bool {
	return true
}

func (l LocalGitRepositoryPlatform) CurrentCommit() string {
	ref, err := l.repo.Head()

	if err != nil {
		return "unknown"
	}

	return ref.Hash().String()
}

func (l LocalGitRepositoryPlatform) CurrentTag() string {
	ref, err := l.repo.Head()

	if err != nil {
		return "unknown"
	}

	tag, err := l.repo.TagObject(ref.Hash())
	if err != nil {
		return ""
	}

	return tag.Name
}

func (l LocalGitRepositoryPlatform) CurrentBranch() string {
	ref, err := l.repo.Head()

	if err != nil {
		return "unknown"
	}

	return strings.TrimLeft(ref.Name().String(), "refs/heads/")
}

func (l LocalGitRepositoryPlatform) CurrentBuildNumber() string {
	return "0"
}

func (b LocalGitRepositoryPlatform) Name() string {
	return "Git"
}

