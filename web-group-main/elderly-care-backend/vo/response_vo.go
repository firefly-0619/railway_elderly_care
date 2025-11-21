package vo

type ResponseVO struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Msg     string      `json:"msg"`
}

func Success(data interface{}) *ResponseVO {
	return &ResponseVO{
		Success: true,
		Data:    data,
	}
}

func Fail(msg string) *ResponseVO {
	return &ResponseVO{
		Success: false,
		Msg:     msg,
	}
}
