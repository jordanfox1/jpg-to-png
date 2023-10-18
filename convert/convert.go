package convert

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func ConvertJpgToPng(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/png":
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode jpeg")
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, errors.Wrap(err, "unable to encode png")
		}

		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("unable to convert %#v to png", contentType)
}

func GetImageFiles(imgFilePath string) ([]fs.DirEntry, error) {
	filesInPath, err := os.ReadDir(imgFilePath)
	if err != nil {
		return nil, err
	}

	return filesInPath, nil
}

func ValidateImgFileExt(expectedFormat string, imgFiles []fs.DirEntry) error {
	if expectedFormat == "png" {
		for _, file := range imgFiles {
			t := file.Type()
			fmt.Print(t)
			if !strings.HasSuffix(file.Name(), ".png") {
				return errors.New("file name is not correct, it should end in .png")
			}
		}
	}

	if expectedFormat == "jpg" {
		for _, file := range imgFiles {
			t := file.Type()
			fmt.Print(t)
			if !strings.HasSuffix(file.Name(), ".jpg") && !strings.HasSuffix(file.Name(), ".jpeg") {
				return errors.New("file name is not correct, it should end in .jpg or jpeg")
			}
		}
	}

	return nil
}
