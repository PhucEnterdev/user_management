package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var allowedExtension = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/jpg":  true,
}

const maxSize = 5 << 20

func ValidateAndSaveFile(fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
	// check extension in filename
	extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtension[extension] {
		return "", errors.New("unsupported file extension")
	}

	// check size
	if fileHeader.Size > maxSize {
		return "", errors.New("file too large (max 5MB)")
	}

	// check content type
	//  (tránh giả mạo file ảnh
	//  nhưng thực chất là 1 loại file khác có thể là mã độc)
	file, err := fileHeader.Open()
	if err != nil {
		return "", errors.New("cannot open file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", errors.New("cannot read file")
	}

	mimetype := http.DetectContentType(buffer)
	if !allowedMimeTypes[mimetype] {
		return "", fmt.Errorf("invalid MIME type: %s", mimetype)
	}

	// change file name
	filename := fmt.Sprintf("%s%s", uuid.New().String(), extension)
	// create folder if not exist
	err = os.MkdirAll("./upload", os.ModePerm)
	if err != nil {
		return "", errors.New("cannot create upload folder")
	}

	// uploadDir = "./upload" + filename: "golang.jpg"
	savePath := filepath.Join(uploadDir, filename)
	if err := saveFile(fileHeader, savePath); err != nil {
		return "", err
	}

	return filename, nil
}

func saveFile(fileHeader *multipart.FileHeader, destination string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}
