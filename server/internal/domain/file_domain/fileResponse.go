package file_domain

type FileResponse struct {
	Filename  string `json:"filename"`
	IsSuccess bool   `json:"isSuccess"`
	Msg       string `json:"msg"`
}
