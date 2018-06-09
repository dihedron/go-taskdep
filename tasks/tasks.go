package tasks

import "fmt"

type Task struct {
	ID           string   `json:"id" yaml:"id"`
	Instructions []string `json:"instructions" yaml:"instructions"`
	Dependencies []string `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
}

func NewFromPath(path string) (*Task, error) {
	return nil, fmt.Errorf("not implemented")
}
