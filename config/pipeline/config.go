package pipeline

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

type Config struct {
	Devops Devops `yaml:"devops"`
}

type Devops struct {
	Cis []Cis `yaml:"cis"`
}

type Cis struct {
	Type string `yaml:"type"`
	Ci   string `yaml:"ci"`
}

func WriteTo(config Config, fpath string) error {
	return nil
}

func GetConfigFrom(fpath string) (*Config, error){
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	config := &Config{}

	if err := yaml.Unmarshal(data, config);err != nil {
		return nil, err
	}

	return config, nil
}

func loadConfigFromFile(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return &Config{}, err
	}

	return loadConfig(data)
}

func loadConfig(data []byte) (*Config, error) {
	config := &Config{}

	yaml.Unmarshal(data, config)

	return config, nil
}
