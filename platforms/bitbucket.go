package platforms

import "os"

type BitbucketPlatform struct {
}

func (b BitbucketPlatform) IsDetected() bool {
	if _, ok := os.LookupEnv("BITBUCKET_REPO_SLUG"); ok {
		return true
	}

	return false
}

func (b BitbucketPlatform) CurrentCommit() string {
	if value, ok := os.LookupEnv("BITBUCKET_COMMIT"); ok {
		return value
	}

	return ""
}

func (b BitbucketPlatform) CurrentTag() string {
	if value, ok := os.LookupEnv("BITBUCKET_TAG"); ok {
		return value
	}

	return ""
}

func (b BitbucketPlatform) CurrentBranch() string {
	if value, ok := os.LookupEnv("BITBUCKET_BRANCH"); ok {
		return value
	}

	return ""
}

func (b BitbucketPlatform) CurrentBuildNumber() string {
	if value, ok := os.LookupEnv("BITBUCKET_BUILD_NUMBER"); ok {
		return value
	}

	return ""
}
