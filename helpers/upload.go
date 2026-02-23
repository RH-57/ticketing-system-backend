package helpers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadAvatar(c *gin.Context, file *multipart.FileHeader) (string, error) {
	uploadPath := "uploads/avatars"

	// pastikan folder ada
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(uploadPath, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	return fullPath, nil
}
