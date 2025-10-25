package model

const (
	pppdCommand = "pppd"
)

type PPPTask struct {
	Tag    string
	Comand string
	Args   []string
}

type PPPResult struct {
	Tag     string
	Success bool
	Error   error
}
