package main

// represents a task that the remote agent should execute
type task struct {
	// required
	Module string `json:"module"`
	Name   string `json:"name"`
	State  string `json:"state"`

	// used for file module
	Content string `json:"content,omitempty"`
}

// represents the result of a task that has been returned by the remote agent
type result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	// TODO: maybe add a boolean for changed, which would represent if the resource has changed
}
