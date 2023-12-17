package main

import "path/filepath"

type Tag struct {
	Name string `json:"name"`
}

type Media struct {
	Url      string
	FilePath string
}

// region - Download

type Download struct {
	Url       string
	FilePath  string
	Error     error
	IsSuccess bool
	Hash      string
}

func (d *Download) MediaType() MediaType {
	extension := filepath.Ext(d.FilePath)

	if extension == ".jpg" || extension == ".jpeg" || extension == ".png" {
		return Image
	} else if extension == ".gif" || extension == ".mp4" || extension == ".m4v" {
		return Video
	}

	return Unknown
}

type ByFilePath []Download

func (a ByFilePath) Len() int           { return len(a) }
func (a ByFilePath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFilePath) Less(i, j int) bool { return a[i].FilePath < a[j].FilePath }

// endregion
