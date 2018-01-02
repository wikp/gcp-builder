package project

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

	return config, nil
}
