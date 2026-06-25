package models

type Workflow struct{
	Name string `yaml:"name"`
	Version float64 `yaml:"version"`
	Description string `yaml:"description"`
	Tasks []Task	`yaml:"tasks"`
}

type Task struct{
	Id string `yaml:"id"`
	Type string	`yaml:"type"`
	Retries int	`yaml:"retries"`
	Parameters any `yaml:"parameters"`//need alternative for this
	DependOn []string `yaml:"dependOn"`
}

