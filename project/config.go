package project

import (
	"errors"
	"fmt"
	"strings"
)

type Variables []Variable

type Configuration struct {
	Project      Project        `yaml:"project"`
	Environments []*Environment `yaml:"environments"`
	Images       []Image        `yaml:"images"`
	Variables    Variables      `yaml:"variables"`
}

type Project struct {
	Name          string `yaml:"name"`
	Domain        string `yaml:"domain"`
	Context       string `yaml:"context"`
	VersionPrefix string `yaml:"versionPrefix"`
}

func (p Project) FullName() string {
	return fmt.Sprintf("%s-%s-%s", p.Domain, p.Context, p.Name)
}

type Environment struct {
	Name       string      `yaml:"name"`
	ServiceKey string      `yaml:"key"`
	Kubernetes Kubernetes  `yaml:"kubernetes"`
	Cloud      GoogleCloud `yaml:"gcloud"`
}

func (e *Environment) envKey(key string) string {
	return fmt.Sprintf("%s_%s", key, strings.ToUpper(e.Name))
}

type GoogleCloud struct {
	Registry string `yaml:"registry"`
	Project  string `yaml:"project"`
}

type Kubernetes struct {
	Cluster   string    `yaml:"cluster"`
	Zone      string    `yaml:"zone"`
	Template  string    `yaml:"template"`
	Variables Variables `yaml:"variables"`
}

type Variable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Image struct {
	Build      string `yaml:"build"`
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
}

func (vars Variables) FindByName(key string) (string, error) {
	for _, v := range vars {
		if v.Name == key {
			return v.Value, nil
		}
	}

	return "", errors.New(fmt.Sprintf("VariableNotFound(%s)", key))
}
