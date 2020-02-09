package main

import (
	"fmt"
	"os/exec"
)

func ModuleShellCmd(task task) (result result) {
	err := RunShellCmd(task.Name)
	if err != nil {
		result.Success = false
		result.Message = err.Error()
	} else {
		result.Success = true
		result.Message = fmt.Sprintf("command executed successfully: %s", task.Name)
	}
	log.Infof("result.success=%v result.message=%s", result.Success, result.Message)
	return
}

func RunShellCmd(shellCmd string) error {
	cmd := exec.Command("sh", "-c", shellCmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, out)
	}
	return nil
}