package service

type ResFilesList struct {
	Access AccessConf `json:"auth"`
	Files  []FileInfo `json:"files"`
}
