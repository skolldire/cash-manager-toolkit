package kit

type StatusCode struct {
	Code     string `json:"code"`
	Msg      string `json:"msg"`
	HttpCode int    `json:"http_code"`
}
