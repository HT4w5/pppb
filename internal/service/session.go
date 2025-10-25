package service

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/HT4w5/pppb/internal/model"
)

func runTask(task model.PPPTask, results chan<- model.PPPResult) {
	cmd := exec.Command(task.Comand, task.Args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("[%s] STARTING process (PID: %d).\n", task.Name, os.Getpid())

	err := cmd.Run()

	if err != nil {
		fmt.Printf("[%s] ERROR: Process failed. Error: %s\n", task.Name, err.Error())
		results <- model.PPPResult{
			Name:    task.Name,
			Success: false,
			Error:   err,
		}
		return
	}

	fmt.Printf("[%s] FINISHED successfully.\n", task.Name)
	results <- model.PPPResult{
		Name:    task.Name,
		Success: true,
		Error:   nil,
	}
}
