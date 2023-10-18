package convert

import (
	"io/fs"
	"net/http"
	"os"
	"testing"
)

func TestGetImageFiles(t *testing.T) {
	var mockFileNames []string
	mockFileNames = append(mockFileNames, "test.jpeg")
	mockFileNames = append(mockFileNames, "test.jpg")

	type args struct {
		imgFilePath string
	}
	tests := []struct {
		name            string
		args            args
		wantedFileNames []string
		wantErr         bool
	}{
		{
			name: "should throw error when non existent dir is passed",
			args: args{
				imgFilePath: "./jpgdddd",
			},
			wantErr: true,
		},
		{
			name: "should return correct file names",
			args: args{
				imgFilePath: "./mocks/jpg",
			},
			wantErr:         false,
			wantedFileNames: mockFileNames,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetImageFiles(tt.args.imgFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageFileNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) > 0 {
				for index, file := range got {
					if file.Name() != tt.wantedFileNames[index] {
						t.Errorf("GetImageFileNames() expected filename - %v not equal to actual filename - %v", file.Name(), tt.wantedFileNames[index])
						return
					}
				}
			}
		})
	}
}

func TestValidateImgFileExt(t *testing.T) {
	mockJpgs, err := os.ReadDir("./mocks/jpg")
	if err != nil {
		t.Fatal(err)
	}
	mockPngs, err := os.ReadDir("./mocks/png")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		expectedFormat string
		imgFiles       []fs.DirEntry
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should throw an error when png files are not expected format",
			args: args{
				expectedFormat: "png",
				imgFiles:       mockJpgs,
			},
			wantErr: true,
		},
		{
			name: "should throw an error when jpg files are not expected format",
			args: args{
				expectedFormat: "jpg",
				imgFiles:       mockPngs,
			},
			wantErr: true,
		},
		{
			name: "should accept jpeg and jpg",
			args: args{
				expectedFormat: "jpg",
				imgFiles:       mockJpgs,
			},
			wantErr: false,
		},
		{
			name: "should accept png",
			args: args{
				expectedFormat: "png",
				imgFiles:       mockPngs,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateImgFileExt(tt.args.expectedFormat, tt.args.imgFiles); (err != nil) != tt.wantErr {
				t.Errorf("ValidateImgFileFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateImgFileType(t *testing.T) {
	type args struct {
		filePath     string
		expectedType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should not return an error if data type is correct",
			args: args{
				filePath:     "./mocks/jpg/test.jpeg",
				expectedType: "jpg",
			},
			wantErr: false,
		},
		{
			name: "should return an error if data type if file is the wrong type",
			args: args{
				filePath:     "./mocks/test.jpg",
				expectedType: "png",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateImgFileType(tt.args.filePath, tt.args.expectedType); (err != nil) != tt.wantErr {
				t.Errorf("ValidateImgFileType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertJpgToPng(t *testing.T) {
	imgData, err := os.ReadFile("./mocks/jpg/test.jpeg")
	if err != nil {
		t.Fatal("error reading image data")
	}

	invalidData, err := os.ReadFile("./mocks/png/test.png")
	if err != nil {
		t.Fatal("error reading image data")
	}

	type args struct {
		imageBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			args: args{
				imageBytes: imgData,
			},
			name:    "should convert jpeg image data to png",
			wantErr: false,
		},
		{
			args: args{
				imageBytes: invalidData,
			},
			name:    "should return error for invalid file",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertJpgToPng(tt.args.imageBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertJpgToPng() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			actualGotContentType := http.DetectContentType(got)
			if actualGotContentType != "image/png" && tt.wantErr != true {
				t.Errorf("image is wrong type expected: image/png, got %v", actualGotContentType)
				return
			}

		})
	}
}
