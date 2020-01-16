package api

type SampleResponse struct {
	param *SampleSub `json:"param"`
}

type SampleSub struct {
	num int `json:"num"`
}
