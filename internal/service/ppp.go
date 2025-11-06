package service

import (
	"os"
	"os/exec"

	"github.com/HT4w5/pppcm/internal/model"
)

func runPPPTask(task model.PPPTask, startSignal <-chan struct{}, results chan<- model.PPPResult) {
	// Block until start signal
	<-startSignal

	cmd := exec.Command(task.Comand, task.Args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		results <- model.PPPResult{
			Tag:     task.Tag,
			Success: false,
			Error:   err,
		}
		return
	}

	results <- model.PPPResult{
		Tag:     task.Tag,
		Success: true,
		Error:   nil,
	}
}
