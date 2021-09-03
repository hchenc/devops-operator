package pipeline

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

type Config struct {
	Devops Devops `yaml:"devops"`
}

type Devops struct {
	Gitlab    Gitlab      `yaml:"gitlab"`
	Pipelines []Pipelines `yaml:"pipelines"`
}

type Gitlab struct {
	Version  string `yaml:"version"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
}

type Pipelines struct {
	Pipeline string `yaml:"pipeline"`
	Ci       string `yaml:"ci"`
	Template string `yaml:"template"`
}

func WriteConfigTo(config *Config, fpath string) error {
	data, _ := yaml.Marshal(config)
	err := ioutil.WriteFile(fpath, data, 0666)
	return err
}

func GetConfigFrom(fpath string) (*Config, error) {
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	config := &Config{}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
