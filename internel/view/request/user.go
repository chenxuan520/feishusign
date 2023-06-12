package request

type (
	ReqSignin struct {
		Code   string `json:"code"`
		Status string `json:"status"`
	}
)
