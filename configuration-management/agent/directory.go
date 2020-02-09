package main

import (
	"fmt"
	"os"
)

func ModuleDirectory(task task) (result result) {
	switch task.State {
	case "present":
		err := CreateDirectory(task.Name)
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "directory created: " + task.Name
		}
	case "absent":
		err := DeleteDirectory(task.Name)
		if err != nil {
			result.Success = false
			result.Message = err.Error()
		} else {
			result.Success = true
			result.Message = "directory deleted: " + task.Name
		}
	default:
		result.Success = false
		result.Message = "state value is invalid: " + task.State
	}
	log.Infof("result.success=%v result.message=%s", result.Success, result.Message)
	return
}

func CreateDirectory(path string) error {
	if path == "" || path == "/" {
		return fmt.Errorf("invalid path specified: %s", path)
	}
	return os.MkdirAll(path, 0755)
	//return os.Chmod(path, os.FileMode(mode))
}

func DeleteDirectory(path string) error {
	return os.RemoveAll(path)
}
