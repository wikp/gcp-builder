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
	
	if strings.HasPrefix(cfg.Branch, project.Project.VersionPrefix) {
		return fmt.Sprintf("%s%s", releasePrefix, strings.TrimPrefix(cfg.Branch, project.Project.VersionPrefix)), nil
	}

	if cfg.CommitSha == "" {
		return "", errors.New("CouldNotResolveVersion")
	}

	return fmt.Sprintf("%s%s-%s", snapshotPrefix, cfg.Branch, cfg.CommitSha), nil
}
