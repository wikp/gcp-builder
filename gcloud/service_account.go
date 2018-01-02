package gcloud

func (i *Client) ActivateServiceAccount(key string) error {
	args := []string{"auth", "activate-service-account", "--key-file", key}

	if _, err := i.CaptureCommand("gcloud", args); err != nil {
		return err
	}

	return nil
}

func (i *Client) GetClusterCredentials(project, cluster, zone string) error {
	args := []string{
		"container",
		"clusters",
		"get-credentials",
		cluster,
		"--project",
		project,
		"--zone",
		zone,
	}

	if _, err := i.CaptureCommand("gcloud", args); err != nil {
		return err
	}

	return nil
}
