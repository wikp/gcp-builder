package platforms

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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
	tag := ""

	if err != nil {
		return "unknown"
	}

	tags, err := l.repo.TagObjects()

	tags.ForEach(func(t *object.Tag) error {
		if t.Target.String() == ref.Hash().String() {
			tag = t.Name
		}

		return nil
	})

	return tag
}

func (l LocalGitRepositoryPlatform) CurrentBranch() string {
	ref, err := l.repo.Head()

	if err != nil {
		return "unknown"
	}

	return strings.TrimPrefix(ref.Name().String(), "refs/heads/")
}

func (l LocalGitRepositoryPlatform) CurrentBuildNumber() string {
	return "0"
}

func (b LocalGitRepositoryPlatform) Name() string {
	return "Git"
}
