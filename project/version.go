package project

import (
	"errors"
	"fmt"
	"github.com/wendigo/gcp-builder/platforms"
	"strings"
)

const snapshotPrefix = "snapshot-"
const releasePrefix = "release-"

func DetectVersion(project *Configuration, platform platforms.Platform) (string, error) {

	if strings.HasPrefix(platform.CurrentTag(), project.Project.VersionPrefix) {
		return fmt.Sprintf("%s%s", releasePrefix, strings.TrimPrefix(platform.CurrentTag(), project.Project.VersionPrefix)), nil
	}

	if platform.CurrentCommit() == "" {
		return "", errors.New("CouldNotResolveVersion")
	}

	return fmt.Sprintf("%s%s-%s", snapshotPrefix, platform.CurrentBranch(), platform.CurrentBuildNumber()), nil
}
