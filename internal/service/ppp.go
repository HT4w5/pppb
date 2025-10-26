package service

import (
	"log"
	"os"
	"os/exec"

	"github.com/HT4w5/pppb/internal/model"
)

func runPPPTask(task model.PPPTask, startSignal <-chan struct{}, results chan<- model.PPPResult) {
	// Block until start signal
	<-startSignal

	cmd := exec.Command(task.Comand, task.Args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("[%s] STARTING dial\n", task.Tag)

	err := cmd.Run()

	if err != nil {
		log.Printf("[%s] ERROR: Process failed: %s\n", task.Tag, err.Error())
		results <- model.PPPResult{
			Tag:     task.Tag,
			Success: false,
			Error:   err,
		}
		return
	}

	log.Printf("[%s] FINISHED successfully.\n", task.Tag)
	results <- model.PPPResult{
		Tag:     task.Tag,
		Success: true,
		Error:   nil,
	}
}
