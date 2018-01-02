package project

import (
	"errors"
	"fmt"
	"github.com/wendigo/gcp-builder/config"
	"strings"
)

const snapshotPrefix = "snapshot-"
const releasePrefix = "release-"

func DetectVersion(project *Configuration, cfg *config.Args) (string, error) {

	prefix := fmt.Sprintf("%s-", project.Project.Name)

	if strings.HasPrefix(cfg.Branch, prefix) {
		return fmt.Sprintf("%s%s", releasePrefix, strings.TrimPrefix(cfg.Branch, prefix)), nil
	}

	if cfg.CommitSha == "" {
		return "", errors.New("CouldNotResolveVersion")
	}

	return fmt.Sprintf("%s%s-%s", snapshotPrefix, cfg.Branch, cfg.CommitSha), nil
}
