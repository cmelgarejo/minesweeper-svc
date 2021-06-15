package responses

type ResponseBase struct {
	Code    int    `json:"code"  example:"12345"`
	Message string `json:"message"  example:"message"`
	Details string `json:"details"  example:"some more details"`
}

type Response struct {
	ResponseBase
	Error    []ResponseError `json:"errors" example:""`
	Response interface{}     `json:"response"  example:""`
}

type ResponseError struct {
	ResponseBase
}
