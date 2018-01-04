package platforms

import "errors"

func GetAll() []Platform {
	return []Platform{
		BitbucketPlatform{},
	}
}

func Detect() (Platform, error) {
	for _, platform := range GetAll() {
		if platform.IsDetected() {
			return platform, nil
		}
	}

	return nil, errors.New("PlatformNotRecognized")
}
