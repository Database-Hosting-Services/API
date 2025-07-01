package sqleditor

import "encoding/json"

type RequestBody struct {
	Query string `json:"query"`
}

type ResponseBody struct {
	Result        json.RawMessage `json:"result"`
	ColumnNames   []string        `json:"column_names"`
	ExecutionTime float64         `json:"execution_time"`
}
