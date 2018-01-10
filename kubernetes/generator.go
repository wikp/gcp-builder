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
	ContainersShas     map[string]string
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
		ContainersShas:     make(map[string]string),
	}, nil
}

func (c Context) Environment() *project.Environment {
	return c.CurrentEnvironment
}

func (c Context) EnvironmentName() string {
	return c.Env
}

func (c Context) Variable(name string) string {
	if val, err := c.CurrentEnvironment.Kubernetes.Variables.FindByName(name); err == nil {
		return val
	}

	if val, err := c.Config.Variables.FindByName(name); err == nil {
		return val
	}

	return fmt.Sprintf("VariableNotFound(%s)", name)
}

func (c Context) Container(name string) string {
	path := fmt.Sprintf("%s/%s/%s:%s", c.CurrentEnvironment.Cloud.Registry, c.Config.Project.FullName(), name, c.Version)

	if project.IsSnapshotVersion(name) {
		if id, exists := c.ContainersShas[path]; exists {
			return fmt.Sprintf("%s/%s/%s@%s", c.CurrentEnvironment.Cloud.Registry, c.Config.Project.FullName(), name, id)
		}
	}

	return path
}

func (c Context) ContainerPath(name string) string {
	return fmt.Sprintf("%s/%s/%s:%s", c.CurrentEnvironment.Cloud.Registry, c.Config.Project.FullName(), name, c.Version)
}

func (c Context) ContainerVersion(name string, version string) string {
	return fmt.Sprintf("%s/%s/%s:%s", c.CurrentEnvironment.Cloud.Registry, c.Config.Project.FullName(), name, version)
}

func (ctx *Context) InterpolateConfig(input string, output string) error {

	inputTemplate, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	tmpl, err := template.New(ctx.Env).Parse(string(inputTemplate))
	if err != nil {
		return err
	}

	buffer := &bytes.Buffer{}

	if err := tmpl.Execute(buffer, ctx); err != nil {
		return err
	}

	log.Printf("Generating '%s' from template '%s' for environment '%s'",
		output,
		input,
		ctx.CurrentEnvironment.Name,
	)

	return ioutil.WriteFile(output, buffer.Bytes(), os.ModePerm)
}
