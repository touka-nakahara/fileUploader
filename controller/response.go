package controller

import "fileUploader/model"

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type FilesUploadRequest struct {
	Data []model.File `json:"data"`
}

type FilesUplaodResponse struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	StatusCode int `json:"statuscode"`
	Message    string `json:"message"`
	ID         string `json:"ID"`
}
