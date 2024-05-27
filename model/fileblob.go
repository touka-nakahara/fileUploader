package model

type FileBlob struct {
	ID   FileID `json:"id"`
	Data []byte `json:"data"`
}
