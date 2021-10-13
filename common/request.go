package common

type Request struct {
	ChannelID string   `json:"channelID" form:"channelID"`
	Contract  string   `json:"contract" form:"contract"`
	Method    string   `json:"method" form:"method"`
	Params    []string `json:"params" form:"params"`
}

type XARequest struct {
	XAID     string     `json:"xaID"`
	Requests []*Request `json:"requests"`
}

type XADetailRequst struct {
	XAID string `json:"xaID"  form:"xaID"`
}
