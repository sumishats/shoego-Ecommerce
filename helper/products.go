package helper

import (
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

func SaveProductImage(fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	// uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// decode the image- convert the file to image format
	img, err := imaging.Decode(file)
	if err != nil {
		return "", err
	}

	resized := imaging.Resize(img, 800, 800, imaging.Lanczos) 
	cropped := imaging.CropCenter(resized, 800, 800) 

	//create folder path if not already exist
	err = os.MkdirAll(uploadDir, os.ModePerm)  //os.modeperm is set permission for create folder
	if err != nil {
		return "", err
	}

	//create unique filename and full path 
	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	fullPath := filepath.Join(uploadDir, filename)

	//save image to path 
	err = imaging.Save(cropped, fullPath) 
	if err != nil {
		return "", err
	}

	return "/uploads/products/" + filename, nil
}

//remove that file if exist when update product image or delete product
func DeleteFileIfExists(path string) {
	_ = os.Remove("." + path)
}
