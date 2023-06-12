package request

type (
	ReqSignin struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}
)
