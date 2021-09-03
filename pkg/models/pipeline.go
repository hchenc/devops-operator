package models

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
