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
|param||*SampleSub|
# GET v1/todolist
 Todo list
##### Response  
|Parameter|Description|Data Type|
|:--|:--|:--|
|tasks| this is tasks |[]*Task|
|sampleMap| mapはrepeatedできない |map[int]*Task|
|task| this is task |*Task|


# SubStructs
#### SampleSub
|Parameter|Description|Data Type|
|:--|:--|:--|
|num||int|
