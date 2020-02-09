package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ModuleFile(task task) (result result) {
	switch task.State {
	case "present":
		err := WriteFile(task.Name, task.Content)
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "file created: " + task.Name
		}
	case "absent":
		err := DeleteFile(task.Name)
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "file deleted: " + task.Name
		}
	default:
		result.Success = false
		result.Message = "state value is invalid: " + task.State
	}
	log.Infof("result.success=%v result.message=%s", result.Success, result.Message)
	return
}

func WriteFile(path, content string) error {
	if path == "" || path == "/" {
		return fmt.Errorf("invalid path specified: %s", path)
	}
	contentBytes := []byte(content)
	return ioutil.WriteFile(path, contentBytes, 0644)
}

func DeleteFile(path string) error {
	return os.Remove(path)
}
