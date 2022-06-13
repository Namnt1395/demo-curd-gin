package response

type Response struct {
	Data         interface{}          `json:"data,omitempty"`
	ErrorCode    string               `json:"error_code,omitempty"`
	ErrorMessage string               `json:"error_msg,omitempty"`
	ErrorFields  []ResponseErrorField `json:"error_fields,omitempty"`
}

type ResponseErrorField struct {
	Field        string `json:"field,omitempty"`
	Tag          string `json:"tag,omitempty"`
	ErrorMessage string `json:"error_msg,omitempty"`
}
