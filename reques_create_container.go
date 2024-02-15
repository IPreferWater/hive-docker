package main

type RequestCreateContainer struct {
	FolderLocation string `json:"folderLocation" binding:"required"`
	ImageName      string `json:"imageName" binding:"required"`
}