package main

// represents a task that this agent should execute
type task struct {
	Module  string `json:"module"`
	Name    string `json:"name"`
	State   string `json:"state"`
	Content string `json:"content,omitempty"`
}

// represents the result of a task that has been returned by this agent
type result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
