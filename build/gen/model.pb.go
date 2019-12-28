package model

import (
	"time"
)

type TodoListResponse struct {
	tasks []Task
}

type Task struct {
	ID        string
	Name      string
	CreatedAt time.Time
}
