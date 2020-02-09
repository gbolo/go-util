package main

import (
	"fmt"
	"os/exec"
)

func ModuleApt(task task) (result result) {
	switch task.State {
	case "update":
		err := AptUpdate()
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "apt was updated"
		}
	case "present":
		err := AptInstall(task.Name)
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "package installed: " + task.Name
		}
	case "absent":
		err := AptRemove(task.Name)
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "package removed: " + task.Name
		}
	default:
		result.Success = false
		result.Message = "state value is invalid: " + task.State
	}
	log.Infof("result.success=%v result.message=%s", result.Success, result.Message)
	return
}

func AptInstall(pkg string) error {
	cmd := exec.Command("apt-get", "-y", "install", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, out)
	}
	return nil
}

func AptRemove(pkg string) error {
	cmd := exec.Command("apt-get", "-y", "remove", pkg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, out)
	}
	return nil
}

func AptUpdate() error {
	cmd := exec.Command("apt-get", "update")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, out)
	}
	return nil
}
