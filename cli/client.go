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
	config   *config.Args
	context  *kubernetes.Context
	logger   *log.Logger
	gcloud   *gcloud.Client
	platform platforms.Platform
}

func New(config *config.Args, cliVersion string) (*Client, error) {

	logger := log.New(
		os.Stdout, "[cli] ", log.Lmicroseconds,
	)

	logger.Printf("gcp-builder version %s booting up...", cliVersion)

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

	ctx, err := kubernetes.NewContext(prj, config.Environment, version)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:   config,
		context:  ctx,
		gcloud:   gcloud.NewClient(),
		platform: platform,
		logger:   logger,
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

			c.logger.Printf("CI/CD platform info:")
			c.logger.Printf("\tName: %s", c.platform.Name())
			c.logger.Printf("\tCurrent branch: %s", c.platform.CurrentBranch())
			c.logger.Printf("\tCurrent tag: %s", c.platform.CurrentTag())
			c.logger.Printf("\tCurrent commit: %s", c.platform.CurrentCommit())
			c.logger.Printf("\tCurrent build number: %s", c.platform.CurrentBuildNumber())

			c.logger.Printf("Project info:")
			c.logger.Printf("\tName: %s", c.context.Config.Project.Name)
			c.logger.Printf("\tDomain: %s", c.context.Config.Project.Domain)
			c.logger.Printf("\tContext: %s", c.context.Config.Project.Context)
			c.logger.Printf("\tVersion: %s", c.context.Version)

			env := c.context.Environment()

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
	env := c.context.Environment()

	c.logger.Printf("Authorizing to project %s and cluster %s...", env.Cloud.Project, env.Kubernetes.Cluster)

	if err := c.gcloud.ActivateServiceAccount(env.ServiceKey); err != nil {
		return err
	}

	err2 := c.gcloud.GetClusterCredentials(env.Cloud.Project, env.Kubernetes.Cluster, env.Kubernetes.Zone)
	if err2 != nil {
		return err2
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

		if image.Dockerfile == "" {
			image.Dockerfile = "Dockerfile"
		}

		if err := client.BuildContainer(c.context, image); err != nil {
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

	ids, err := c.gatherImagesShas()
	if err != nil {
		return err
	}

	c.context.ContainersShas = ids

	return c.context.InterpolateConfig(
		c.context.CurrentEnvironment.Kubernetes.Template,
		filename,
	)
}

func (c *Client) gatherImagesShas() (map[string]string, error) {
	images := make(map[string]string, 0)

	client, err := containers.New(c.gcloud)
	if err != nil {
		return images, err
	}

	c.logger.Printf("Checking images SHAs...")

	for _, image := range c.context.Config.Images {
		id, err := client.ContainerSha256(c.context, image)
		if err != nil {
			return images, err
		}

		images[c.context.ContainerPath(image.Name)] = id
	}

	for tag, id := range images {
		c.logger.Printf("\t%s has sum: %s", tag, id)
	}

	return images, nil
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
	env := c.context.Environment()
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
