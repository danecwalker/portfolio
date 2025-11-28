package notion

import "time"

type FileObject struct {
	Type       string     `json:"type"`
	File       File       `json:"file"`
	FileUpload FileUpload `json:"file_upload "`
	External   External   `json:"external"`
	Emoji      string     `json:"emoji"`
}

type File struct {
	URL        string    `json:"url"`
	ExpiryTime time.Time `json:"expiry_time"`
}

type FileUpload struct {
	ID string `json:"id"`
}

type External struct {
	URL string `json:"url"`
}

func (fo *FileObject) GetURL() string {
	switch fo.Type {
	case "file":
		return fo.File.URL
	case "external":
		return fo.External.URL
	default:
		return ""
	}
}
