package model

import "time"

type FileID int

type File struct {
	ID          FileID `json:"id"`
	Name        string `json:"name"`
	Size        int    `json:"size"`
	Extension   string `json:"extension"`
	Description string `json:"description"`

	Uuid     string
	Password string `json:"password"`

	Thumbnail []byte `json:"thumnbnail"`

	IsAvailable time.Time `json:"is_available"`
	UpdateDate  time.Time `json:"update_date"`
	UploadDate  time.Time `json:"upload_date"`

	HasPassword bool `json:"has_password"`
}

type FileBlob struct {
	ID   FileID `json:"id"`
	Data []byte `json:"data"`
}

type GetQueryParam struct {
	Extension    string
	Is_available string
	Search       string
	Sort         string
	Ordered      string
	Page         int
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type UploadResponse struct {
	ID FileID `json:"id"`
}
