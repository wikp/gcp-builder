package containers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wendigo/gcp-builder/gcloud"
	"github.com/wendigo/gcp-builder/kubernetes"
	"github.com/wendigo/gcp-builder/project"
	"log"
	"os"
)

type Client struct {
	gcloud *gcloud.Client
	logger *log.Logger
}

func New(gcloud *gcloud.Client) (*Client, error) {
	return &Client{
		gcloud: gcloud,
		logger: log.New(
			os.Stdout, "[containers] ", log.Lmicroseconds,
		),
	}, nil
}

func (c *Client) BuildContainer(context *kubernetes.Context, image project.Image) error {
	tag := context.ContainerPath(image.Name)

	if image.Dockerfile == "" {
		image.Dockerfile = "Dockerfile"
	}

	dockerfile := fmt.Sprintf("%s-%s", image.Dockerfile, context.CurrentEnvironment.Name)

	c.logger.Printf("Building container %s [%s] from build context %s", image.Name, tag, image.Build)

	if err := context.InterpolateConfig(image.Dockerfile, dockerfile); err != nil {
		return err
	}

	args := []string{
		"build",
		"--file",
		dockerfile,
		"--compress",
		"-t",
		tag,
		image.Build,
	}

	if err := c.gcloud.RunCommand("docker", args); err != nil {
		return err
	}

	if err := os.Remove(dockerfile); err != nil {
		return err
	}

	return nil
}

func (c *Client) PushContainer(tag string) error {
	c.logger.Printf("Pushing container %s", tag)

	args := []string{
		"docker",
		fmt.Sprintf("--docker-host=%s", os.Getenv("DOCKER_HOST")),
		"--",
		"push",
		tag,
	}

	if err := c.gcloud.RunCommand("gcloud", args); err != nil {
		return err
	}

	return nil
}

func (c *Client) ContainerSha256(context *kubernetes.Context, image project.Image) (string, error) {
	tag := context.ContainerPath(image.Name)

	args := []string{
		"inspect",
		tag,
	}

	output, err := c.gcloud.CaptureCommand("docker", args)
	if err != nil {
		return "", err
	}

	inspect := make(inspect, 0)

	if err := json.Unmarshal(output, &inspect); err != nil {
		return "", err
	}

	if len(inspect) == 0 {
		return "", errors.New(fmt.Sprintf("Could not find image with tag %s to inspect", tag))
	}

	return inspect[0].Id, nil
}
