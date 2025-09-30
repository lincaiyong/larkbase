package larkbase

type Record struct {
	Id     string            `json:"id"`
	Fields map[string]IField `json:"fields"`
}
