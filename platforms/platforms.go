package platforms

import "errors"

func GetAll() []Platform {
	platforms := []Platform{
		&BitbucketPlatform{},
	}

	if platform, err := NewLocalGitRepositoryPlatform(); err == nil {
		platforms = append(platforms, *platform)
	}

	return platforms
}

func Detect() (Platform, error) {
	for _, platform := range GetAll() {
		if platform.IsDetected() {
			return platform, nil
		}
	}

	return nil, errors.New("PlatformNotRecognized")
}
