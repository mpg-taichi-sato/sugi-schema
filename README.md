# What is this?
protoc-gen-genta is an plugin for protocol buffers.  
this plugin can generate go, json, apidoc and csfields from .proto

like this.

from 

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

```pb/api/api.proto
syntax = "proto3";
package pb.api;

import "model.proto";
import "option/http.proto";
import "google/protobuf/empty.proto";

// Serviceのコメント
service SampleService {
    
    // サンプルで作った関数
    rpc SampleRPC(pb.model.Task) returns (pb.model.TodoListResponse){
        option (pb.option.http) = {
            method: "GET"
            path: "v1/samplerpc"
        };
    }

    // Todo list
    rpc TodoList(google.protobuf.Empty) returns (pb.model.TodoListResponse){
        option (pb.option.http) = {
            method: "GET"
            path: "v1/todolist"
        };
    }
}

service SampleService2 {
    // サンプルで作った関数2222
    rpc SampleRPC2(pb.model.Task) returns (pb.model.TodoListResponse){
        option (pb.option.http) = {
            method: "GET"
            path: "v1/samplerpc2"
        };
    }
}
```

to
```build/gen/model.go
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

```build/gen/api/api.md
# GET v1/samplerpc
 サンプルで作った関数
##### Parameters  
|Parameter|Description|Data Type|
|:--|:--|:--|
|ID||string|
|Name| task name |string|
|CreatedAt||time.Time|
##### Response  
|Parameter|Description|Data Type|
|:--|:--|:--|
|tasks| this is tasks |[]*Task|
|sampleMap| mapはrepeatedできない |map[int]*Task|
|task| this is task |*Task|
# GET v1/todolist
 Todo list
##### Response  
|Parameter|Description|Data Type|
|:--|:--|:--|
|tasks| this is tasks |[]*Task|
|sampleMap| mapはrepeatedできない |map[int]*Task|
|task| this is task |*Task|
# GET v1/samplerpc2
 サンプルで作った関数2222
##### Parameters  
|Parameter|Description|Data Type|
|:--|:--|:--|
|ID||string|
|Name| task name |string|
|CreatedAt||time.Time|
##### Response  
|Parameter|Description|Data Type|
|:--|:--|:--|
|tasks| this is tasks |[]*Task|
|sampleMap| mapはrepeatedできない |map[int]*Task|
|task| this is task |*Task|

```

# How To Use
`make generate`  

please edit pb and Makefile for your project.

# parameter
go  
json  
apidoc  
csfields  

# extend
You can extend this package by create generator and process.


