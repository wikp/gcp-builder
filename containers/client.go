package containers

import (
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

	tag := context.Container(image.Name)

	if image.Dockerfile == "" {
		image.Dockerfile = "Dockerfile"
	}

	dockerfile := fmt.Sprintf("%s-%s", image.Dockerfile, context.CurrentEnvironment.Name)

	c.logger.Printf("Building container %s [%s] from build context %s", image.Name, tag, image.Build)

	if err := kubernetes.InterpolateConfig(context, image.Dockerfile, dockerfile); err != nil {
		return err
	}

	args := []string{
		"docker",
		fmt.Sprintf("--docker-host=%s", os.Getenv("DOCKER_HOST")),
		"--",
		"build",
		"--file",
		dockerfile,
		"-t",
		tag,
		image.Build,
	}

	if err := c.gcloud.RunCommand("gcloud", args); err != nil {
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
