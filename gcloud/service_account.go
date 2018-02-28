package gcloud

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
)

const keyFilename = "key.json"

func (i *Client) SetDefaultProject(project string) error {
	if project == "" {
		return errors.New("DefaultProjectEmpty")
	}

	args := []string{"config", "set", "project", project}

	if err := i.RunCommand("gcloud", args); err != nil {
		return err
	}

	return nil
}

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

	args := []string{"auth", "activate-service-account", "--key-file", keyFilename}

	if err := i.RunCommand("gcloud", args); err != nil {
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

	if err := i.RunCommand("gcloud", args); err != nil {
		return err
	}

	return nil
}
