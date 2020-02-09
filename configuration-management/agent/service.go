package main

import (
	"fmt"
	"os/exec"
)

func ModuleService(task task) (result result) {
	err := ServiceCtl(task.Name, task.State)
	if err != nil {
		result.Success = false
		result.Message = err.Error()
	} else {
		result.Success = true
		result.Message = fmt.Sprintf("service %s state is: %s", task.Name, task.State)
	}
	log.Infof("result.success=%v result.message=%s", result.Success, result.Message)
	return
}

func ServiceCtl(service, operation string) error {
	if operation != "start" &&  operation != "stop" &&  operation != "restart" {
		return fmt.Errorf("service operation not supported: %s", operation)
	}
	cmd := exec.Command("systemctl", operation, service)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, out)
	}
	return nil
}