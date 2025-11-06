package model

const (
	PPPDefaultCommand = "pppd"
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

type PPPLink struct {
	Tag    string
	Task   *PPPTask
	Up     bool
	IFName string
}
