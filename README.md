# protobuf-code-generator
generate .go from .proto

WIP csharp


# sample
```model.proto
syntax = "proto3";
package model;

import "google/protobuf/timestamp.proto";

message TodoListResponse {
    repeated Task tasks = 1;
}

message Task {
    string ID = 1;
    string Name = 2;
    google.protobuf.Timestamp CreatedAt = 3;
}
```

```model.pb.go
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
```

# HowToUse
make generate

# parameter
go
json
