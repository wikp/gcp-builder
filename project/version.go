package project

import (
	"errors"
	"fmt"
	"github.com/wendigo/gcp-builder/platforms"
	"strings"
)

const snapshotSuffix = "snapshot"

func DetectVersion(project *Configuration, platform platforms.Platform) (string, error) {

	if strings.HasPrefix(platform.CurrentTag(), project.Project.VersionPrefix) {
		return strings.TrimPrefix(platform.CurrentTag(), project.Project.VersionPrefix), nil
	}

	if platform.CurrentCommit() == "" {
		return "", errors.New("CouldNotResolveVersion")
	}

	return fmt.Sprintf("%s-%s", platform.CurrentBranch(), snapshotSuffix), nil
}

func IsSnapshotVersion(version string) bool {
	return strings.HasSuffix(version, snapshotSuffix)
}
