package cli

import (
	"errors"
	"fmt"
	"github.com/wendigo/gcp-builder/config"
	"github.com/wendigo/gcp-builder/containers"
	"github.com/wendigo/gcp-builder/gcloud"
	"github.com/wendigo/gcp-builder/kubernetes"
	"github.com/wendigo/gcp-builder/platforms"
	"github.com/wendigo/gcp-builder/project"
	"log"
	"os"
	"reflect"
)

type Client struct {
	config  *config.Args
	context kubernetes.Context
	logger  *log.Logger
	gcloud  *gcloud.Client
}

func New(config *config.Args) (*Client, error) {

	prj, err := project.FromFile(config.ProjectConfig)
	if err != nil {
		return nil, err
	}

	platform, err := platforms.Detect()
	if err != nil {
		return nil, err
	}

	version, err := project.DetectVersion(prj, platform)
	if err != nil {
		return nil, err
	}

	ctx := kubernetes.Context{
		Config:  prj,
		Env:     config.Environment,
		Version: version,
	}

	return &Client{
		config:  config,
		context: ctx,
		gcloud:  gcloud.NewClient(),
		logger: log.New(
			os.Stdout, "[cli] ", log.Lmicroseconds,
		),
	}, nil
}

func (c *Client) Run() error {
	if err := c.init(); err != nil {
		return err
	}

	if reflect.DeepEqual(c.config.Steps, []string{"all"}) {
		c.config.Steps = []string{
			"info",
			"auth",
			"build",
			"push",
			"deploy-config",
			"deploy",
			"wait-for-deploy",
		}
	}

	for _, step := range c.config.Steps {
		switch step {
		case "info":
			c.logger.Printf("Project info:")
			c.logger.Printf("\tName: %s", c.context.Config.Project.Name)
			c.logger.Printf("\tDomain: %s", c.context.Config.Project.Domain)
			c.logger.Printf("\tContext: %s", c.context.Config.Project.Context)
			c.logger.Printf("\tVersion: %s", c.context.Version)

			env, err := c.context.Environment()
			if err != nil {
				return err
			}

			c.logger.Printf("Environment info:")
			c.logger.Printf("\tName: %s", env.Name)
			c.logger.Printf("\tProject: %s", env.Cloud.Project)
			c.logger.Printf("\tRegistry: %s", env.Cloud.Registry)
			c.logger.Printf("\tCluster: %s", env.Kubernetes.Cluster)
			c.logger.Printf("\tZone: %s", env.Kubernetes.Zone)
		case "auth":
			if err := c.authorize(); err != nil {
				return err
			}
		case "build":
			if err := c.buildContainers(); err != nil {
				return err
			}

		case "push":
			if err := c.pushContainers(); err != nil {
				return err
			}

		case "deploy-config":
			if err := c.buildDeployment(); err != nil {
				return err
			}

		case "deploy":
			if err := c.deploy(); err != nil {
				return err
			}

		case "wait-for-deploy":

		default:
			return errors.New(fmt.Sprintf("UnrecognizedStep(%s)", step))
		}
	}

	return nil
}

func (c *Client) authorize() error {
	c.logger.Printf("Authorizing...")

	if err := c.gcloud.ActivateServiceAccount(c.config.KeyFile); err != nil {
		return err
	}

	env, err := c.context.Environment()
	if err != nil {
		return err
	}

	err2 := c.gcloud.GetClusterCredentials(env.Cloud.Project, env.Kubernetes.Cluster, env.Kubernetes.Zone)
	if err2 != nil {
		return err
	}

	return nil
}

func (c *Client) init() error {
	c.logger.Printf("Installing dependencies...")

	if err := c.gcloud.Install(); err != nil {
		return err
	}

	return nil
}

func (c *Client) buildContainers() error {

	client, err := containers.New(c.gcloud)
	if err != nil {
		return err
	}

	c.logger.Printf("Building containers")

	for _, image := range c.context.Config.Images {
		if err := client.BuildContainer(image.Name, image.Build, c.context.Container(image.Name)); err != nil {
			c.logger.Printf("Error building container: %s", err)
			return err
		}
	}

	return nil
}

func (c *Client) buildDeployment() error {
	filename, err := c.deploymentFile()
	if err != nil {
		return err
	}

	return kubernetes.InterpolateConfig(c.context, filename)
}

func (c *Client) deploy() error {
	filename, err := c.deploymentFile()
	if err != nil {
		return err
	}

	if err := c.gcloud.RunCommand("kubectl", []string{"apply", "-f", filename}); err != nil {
		return err
	}

	if err := os.Remove(filename); err != nil {
		return err
	}

	return nil
}

func (c *Client) deploymentFile() (string, error) {
	env, err := c.context.Environment()
	if err != nil {
		return "", err
	}

	projectName := c.context.Config.Project.FullName()

	return fmt.Sprintf("deployment-%s-%s.yml", projectName, env.Name), nil
}

func (c *Client) pushContainers() error {

	client, err := containers.New(c.gcloud)
	if err != nil {
		return err
	}

	c.logger.Printf("Pushing containers")

	for _, image := range c.context.Config.Images {
		if err := client.PushContainer(c.context.Container(image.Name)); err != nil {
			c.logger.Printf("Error pushing container: %s", err)
			return err
		}
	}

	return nil
}
