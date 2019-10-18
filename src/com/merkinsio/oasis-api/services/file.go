package services

import (
	"com/merkinsio/oasis-api/config"
	"com/merkinsio/oasis-api/domain"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
)

// CreateFile creates a file and stores it in the disk
func CreateFile(fileName string, uploadedFile multipart.File) (string, *domain.File, error) {

	file := domain.File{
		ID:         uuid.NewV4().String(),
		UploadedAt: time.Now().UTC(),
		FileName:   fileName,
	}

	filePath := fmt.Sprintf("%s/%s", config.Config.GetString("uploads.path"), file.ID)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Error in services.StoreFile -> error: %v", err.Error())
		return "", nil, err
	}

	defer f.Close()
	io.Copy(f, uploadedFile)

	return filePath, &file, nil
}
