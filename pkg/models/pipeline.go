package models

type Config struct {
	Devops Devops `yaml:"Devops"`
}

type Devops struct {
	Harbor    Harbor      `yaml:"Harbor"`
	Gitlab    Gitlab      `yaml:"Gitlab"`
	Pipelines []Pipelines `yaml:"Pipelines"`
}

type Harbor struct {
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Host     string `yaml:"Host"`
}

type Gitlab struct {
	Password string `yaml:"Password"`
	Port     string `yaml:"Port"`
	Token    string `yaml:"Token"`
	User     string `yaml:"User"`
	Version  string `yaml:"Version"`
	Host     string `yaml:"Host"`
}

type Pipelines struct {
	Ci       string `yaml:"Ci"`
	Pipeline string `yaml:"Pipeline"`
	Template string `yaml:"Template"`
}
