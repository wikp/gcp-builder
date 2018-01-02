package containers

import (
	"github.com/wendigo/gcp-builder/gcloud"
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

func (c *Client) BuildContainer(name, path, tag string) error {
	c.logger.Printf("Building container %s [%s] from build context %s", name, tag, path)

	args := []string{
		"docker",
		"--",
		"build",
		"-t",
		tag,
		path,
	}

	if err := c.gcloud.RunCommand("gcloud", args); err != nil {
		return err
	}

	return nil
}

func (c *Client) PushContainer(tag string) error {
	c.logger.Printf("Pushing container %s", tag)

	args := []string{
		"docker",
		"--",
		"push",
		tag,
	}

	if err := c.gcloud.RunCommand("gcloud", args); err != nil {
		return err
	}

	return nil
}
