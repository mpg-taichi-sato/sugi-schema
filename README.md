# What is this?
protoc-gen-genta is an plugin for protocol buffers.
this plugin can generate go, json, apidoc and csfields from .proto

from model.proto
```pb/model.proto
// コメントsyntax
syntax = "proto3"; // コメントsyntax2
// コメントsyntax3

// コメントpackage
package pb.model; // コメントpackage2

import "google/protobuf/timestamp.proto";

// this is TodoList
message TodoListResponse {
    repeated Task tasks = 1; // this is tasks
    map<int32, Task> sampleMap = 2; // mapはrepeatedできない
    Task task = 3; // this is task
}

message Task {
    string ID = 1;
    string Name = 2; // task name
    google.protobuf.Timestamp CreatedAt = 3;
}
```

to model.pb.go
```
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
```

# How To Use
make generate
please edit pb and Makefile for your project.

# parameter
go  
json  
apidoc  
csfields  

# extend
You can extend this package by create generator and process.


