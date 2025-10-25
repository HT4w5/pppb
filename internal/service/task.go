package service

import (
	"log"
	"os"
	"os/exec"

	"github.com/HT4w5/pppb/internal/model"
)

func runTask(task model.PPPTask, results chan<- model.PPPResult) {
	cmd := exec.Command(task.Comand, task.Args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("[%s] STARTING process (PID: %d).\n", task.Tag, os.Getpid())

	err := cmd.Run()

	if err != nil {
		log.Printf("[%s] ERROR: Process failed. Error: %s\n", task.Tag, err.Error())
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
