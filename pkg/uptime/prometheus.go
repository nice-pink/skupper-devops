package uptime

// prometheus metrics

type Data struct {
	ResultType string `json:"resultType"`
	Result     []Metric
}

type Response struct {
	Status string
	Data   Data
}

type Metric struct {
	Info  Info `json:"metric"`
	Value []interface{}
}

type Info struct {
	Name     string `json:"__name__"`
	Instance string
	Job      string
	Cluster  string
}

type Value struct {
	Timestamp int
	Value     string
}
