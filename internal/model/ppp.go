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

type PPPConnection struct {
	Tag     string
	Task    *PPPTask
	Up      bool // Created in the first place
	Healthy bool // Ping success
	IFName  string
}
