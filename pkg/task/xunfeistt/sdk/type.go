package sdk

type XunfeiSDK struct {
	BaseUrl     string
	file_path   string
	cur_sliceid string
	taskid      string
	base        BaseReq
}

type BaseReq struct {
	AppID string `json:"app_id"`
	Signa string `json:"signa"`
	Ts    string `json:"ts"`
}

type Status struct {
	Desc   string `json:"desc"`
	Status int    `json:"status"`
}

type TaskIdReq struct {
	BaseReq
	TaskId string `json:"task_id"`
}

type PrepareReq struct {
	Language string `json:"language"`
}

type PrepareFullReq struct {
	BaseReq
	PrepareReq
	FileLen  string `json:"file_len"`
	FileName string `json:"file_name"`
	SliceNum string `json:"slice_num"`
}

type BaseResp struct {
	Ok     int64       `json:"ok"`
	ErrNo  int64       `json:"err_no"`
	Failed interface{} `json:"failed"`
	Data   string      `json:"data"`
}
