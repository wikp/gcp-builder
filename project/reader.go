package project

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func FromFile(filename string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Configuration{}

	if err := yaml.Unmarshal([]byte(bytes), config); err != nil {
		return nil, err
	}

	if err := fillWithEnvironmentVariables(config); err != nil {
		return nil, err
	}

	return config, nil
}

func fillWithEnvironmentVariables(conf *Configuration) (err error) {

	for _, env := range conf.Environments {
		if val, err := fillWithEnvironmentVariable(env.Kubernetes.Zone, env, "KUBERNETES_ZONE"); err != nil {
			return err
		} else {
			env.Kubernetes.Zone = val
		}

		if val, err := fillWithEnvironmentVariable(env.Kubernetes.Cluster, env, "KUBERNETES_CLUSTER"); err != nil {
			return err
		} else {
			env.Kubernetes.Cluster = val
		}

		if val, err := fillWithEnvironmentVariable(env.Kubernetes.Template, env, "KUBERNETES_TEMPLATE"); err != nil {
			return err
		} else {
			env.Kubernetes.Template = val
		}

		if val, err := fillWithEnvironmentVariable(env.Cloud.Project, env, "GCLOUD_PROJECT"); err != nil {
			return err
		} else {
			env.Cloud.Project = val
		}

		if val, err := fillWithEnvironmentVariable(env.Cloud.Registry, env, "GCLOUD_REGISTRY"); err != nil {
			return err
		} else {
			env.Cloud.Registry = val
		}

		if val, err := fillWithEnvironmentVariable(env.ServiceKey, env, "SERVICE_KEY"); err != nil {
			return err
		} else {
			env.ServiceKey = val
		}
	}

	return nil
}

func environmentKey(environment *Environment, name string) string {
	return fmt.Sprintf("%s_%s", name, strings.ToUpper(environment.Name))
}

func fillWithEnvironmentVariable(value string, env *Environment, key string) (string, error) {

	if value == "" {
		if value, ok := os.LookupEnv(environmentKey(env, key)); !ok {
			return "", errors.New(fmt.Sprintf("ProjectConfigurationMissing(%s,%s)", key, env.Name))
		} else {
			return value, nil
		}
	}

	return value, nil
}
