package validator

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type fileConfig struct {
	MaxSize int64 //max size for the file
	AllowedExt []string  // allowed extenstion
}

var YmlDefaultCofig = fileConfig{
	MaxSize: 1 << 20, // 1 mb limit for yml file
	AllowedExt: []string{
		".yaml",
	},
}

func ValidateYML(file *multipart.FileHeader) error {
	
	if file.Size > YmlDefaultCofig.MaxSize {
		return errors.New("Yml file exceeded")
	}

	ext:= strings.ToLower(filepath.Ext(file.Filename))
	extFlag := false

	for _,i := range YmlDefaultCofig.AllowedExt {
		if i == ext {
			extFlag = true
			break
		}
	}
	if !extFlag {
		return errors.New("File is not supported")
	}
	return nil
}