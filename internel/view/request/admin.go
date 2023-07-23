package request

type (
	ReqMeetingCreate struct {
		BeginTime int64 `json:"begin_time"`
		EndTime   int64 `json:"end_time"`
	}
)
