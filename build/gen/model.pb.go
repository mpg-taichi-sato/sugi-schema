// コメントsyntax
// コメントsyntax2
// コメントpackage
// コメントsyntax2
package model

import (
	"time"
)

//TodoListResponse  this is TodoList
type TodoListResponse struct {
	tasks     []*Task       // this is tasks
	sampleMap map[int]*Task // mapはrepeatedできない
	task      *Task         // this is task
}

type Task struct {
	ID        string
	Name      string // task name
	CreatedAt time.Time
}
