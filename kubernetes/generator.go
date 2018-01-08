package kubernetes

import "github.com/wendigo/gcp-builder/project"
import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

type Context struct {
	Config             *project.Configuration
	Env                string
	Version            string
	CurrentEnvironment *project.Environment
}

func NewContext(prj *project.Configuration, environment string, version string) (*Context, error) {

	var currentEnvironment *project.Environment = nil

	for _, env := range prj.Environments {
		if env.Name == environment {
			currentEnvironment = env
		}
	}

	if currentEnvironment == nil {
		return nil, errors.New(fmt.Sprintf("UnrecognizedEnvironment(%s)", environment))
	}

	return &Context{
		Config:             prj,
		Env:                environment,
		Version:            version,
		CurrentEnvironment: currentEnvironment,
	}, nil
}

func (c Context) Environment() (*project.Environment, error) {
	for _, env := range c.Config.Environments {
		if env.Name == c.Env {
			return env, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("UnrecognizedEnvironment(%s)", c.Env))
}

func (c Context) EnvironmentName() string {
	return c.Env
}

func (c Context) Variable(name string) string {
	if env, err := c.Environment(); err == nil {
		if val, err := env.Kubernetes.Variables.FindByName(name); err == nil {
			return val
		}
	}

	if val, err := c.Config.Variables.FindByName(name); err == nil {
		return val
	}

	return fmt.Sprintf("VariableNotFound(%s)", name)
}

func (c Context) Container(name string) string {
	if env, err := c.Environment(); err == nil {
		return fmt.Sprintf("%s/%s/%s:%s", env.Cloud.Registry, c.Config.Project.FullName(), name, c.Version)
	} else {
		return fmt.Sprintf("ContainerNotFound(%s)", name)
	}
}

func (c Context) ContainerVersion(name string, version string) string {
	if env, err := c.Environment(); err == nil {
		return fmt.Sprintf("%s/%s/%s:%s", env.Cloud.Registry, c.Config.Project.FullName(), name, version)
	} else {
		return fmt.Sprintf("ContainerNotFound(%s)", name)
	}
}

func InterpolateConfig(context *Context, input string, output string) error {

	inputTemplate, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	tmpl, err := template.New(context.Env).Parse(string(inputTemplate))
	if err != nil {
		return err
	}

	buffer := &bytes.Buffer{}

	if err := tmpl.Execute(buffer, context); err != nil {
		return err
	}

	log.Printf("Generating '%s' from template '%s' for environment '%s'",
		output,
		input,
		context.CurrentEnvironment.Name,
	)

	return ioutil.WriteFile(output, buffer.Bytes(), os.ModePerm)
}
