package utils

import (
	"mime/multipart"
	"os"
	"runtime"
	"strings"

	"github.com/disintegration/imaging"
)

func CreateFolder(id string) error {
	_, err := os.Stat(id)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(id, 0755)
		if errDir != nil {
			return err
		}
	}
	return nil
}

//	CheckIfValidFileFormat gets *multipart.FileHeader and returns true and file type if file format is png,jpeg,pdf,tiff or jpg.
func CheckIfValidFileFormat(file *multipart.FileHeader) (bool, string) {
	b := file.Header["Content-Disposition"]
	delimiter := "."
	extention := strings.Join(strings.Split(b[0], delimiter)[1:], delimiter)
	if extention == "png\"" || extention == "jpeg\"" || extention == "pdf\"" || extention == "tiff\"" || extention == "jpg\"" {
		val := len(extention) - 1
		return true, extention[:val]
	}
	return false, ""
}

//	CheckIfValidFileSize gets image path and check the size of image, Returns true if size is lower than 1mb
func CheckIfValidFileSize(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	var bytes int64
	bytes = stat.Size()
	var kilobytes int64
	kilobytes = (bytes / 1024)
	if kilobytes < 1024 {
		return true
	} else {
		return false
	}
}

//	ResizeDimension gets image path and rezise image as given arguments
func ResizeDimension(path string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	img, err := imaging.Open(path)
	if err != nil {
		return err
	}

	dstimg := imaging.Resize(img, 500, 0, imaging.Box)
	err = imaging.Save(dstimg, path)
	if err != nil {
		return err
	}
	return nil
}
