// コメント1
// コメント2
package model

type KeyValue struct {
	key   string `json:"key"`
	value string `json:"value"`
}

type Item struct {
	id   int    `json:"ID"`
	name string `json:"name"`
}
