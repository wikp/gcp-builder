package project

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
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
		if val, err := fillWithEnvironmentVariable(env.Kubernetes.Zone, "", env.envKey("KUBERNETES_ZONE")); err != nil {
			return err
		} else {
			env.Kubernetes.Zone = val
		}

		if val, err := fillWithEnvironmentVariable(env.Kubernetes.Cluster, "", env.envKey("KUBERNETES_CLUSTER")); err != nil {
			return err
		} else {
			env.Kubernetes.Cluster = val
		}

		if val, err := fillWithEnvironmentVariable(env.Kubernetes.Template, "deployment.yml", env.envKey("KUBERNETES_TEMPLATE")); err != nil {
			return err
		} else {
			env.Kubernetes.Template = val
		}

		if val, err := fillWithEnvironmentVariable(env.Cloud.Project, "", env.envKey("GCLOUD_PROJECT")); err != nil {
			return err
		} else {
			env.Cloud.Project = val
		}

		if val, err := fillWithEnvironmentVariable(env.Cloud.Registry, "", env.envKey("GCLOUD_REGISTRY")); err != nil {
			return err
		} else {
			env.Cloud.Registry = val
		}

		if val, err := fillWithEnvironmentVariable(env.ServiceKey, "", env.envKey("SERVICE_KEY")); err != nil {
			return err
		} else {
			env.ServiceKey = val
		}
	}

	return nil
}

func fillWithEnvironmentVariable(value string, defaultValue string, key string) (string, error) {

	if value == "" {
		if value, ok := os.LookupEnv(key); !ok {
			if defaultValue == "" {
				return "", errors.New(fmt.Sprintf("ProjectConfigurationMissing(%s)", key))
			}

			return defaultValue, nil
		} else if value != "" {
			return value, nil
		}
	}

	return value, nil
}
