package gcloud

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
)

const keyFilename = "key.json"

func (i *Client) ActivateServiceAccount(key string) error {

	if key == "" {
		return errors.New("ServiceKeyEmpty")
	}

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(keyFilename, decodedKey, os.ModePerm); err != nil {
		return err
	}

	defer func() {
		os.Remove(keyFilename)
	}()

	args := []string{"auth", "activate-service-account", "--key-file", keyFilename}

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
