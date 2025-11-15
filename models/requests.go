package models

type Query struct {
	Query string `json:"query"`
}

type Url struct {
	Url string `json:"url"`
}

type Message struct {
	Msg string `json:"msg"`
}

type QueryResponse struct {
	Results []string `json:"results"`
}
