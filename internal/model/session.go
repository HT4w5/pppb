package model

const (
	pppdCommand = "pppd"
)

type PPPTask struct {
	Name   string   `yaml:"name"`
	Comand string   `yaml:"command"`
	Args   []string `yaml:"args"`
}

type PPPResult struct {
	Name    string
	Success bool
	Error   error
}
