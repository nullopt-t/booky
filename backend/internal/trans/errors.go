package trans

type ApiErr struct {
	Code    ApiErrCode `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
}
