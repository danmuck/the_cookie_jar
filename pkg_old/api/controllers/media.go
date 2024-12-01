package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/danmuck/the_cookie_jar/pkg/api/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	MAX_FILE_SIZE = 100 << 20
	UPLOAD_DIR    = "uploads"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"video/mp4":  true,
}

// This creates the upload directory if it doesn't exist.
func init() {
	absPath, err := filepath.Abs(UPLOAD_DIR)
	if err != nil {
		panic(fmt.Sprintf("Failed to get path for uploads: %v", err))
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
}

// validateFile checks the file's MIME type and returns the correct extension
func validateFile(file multipart.File) error {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("invalid file headers %w", err)
	}

	//Resetting the file to the beginning
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file position: %w", err)
	}

	contentType := http.DetectContentType(buffer[:n])
	if !allowedMimeTypes[contentType] {
		return fmt.Errorf("unsupported file type: %s", contentType)
	}

	return nil
}

func saveFile(file multipart.File, path string) error {
	dst, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// UploadMedia handles file uploads
func UploadMedia(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	if header.Size > MAX_FILE_SIZE {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large, maximum size is 100MB"})
		return
	}

	if err := validateFile(file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", id, ext)

	absPath, err := filepath.Abs(UPLOAD_DIR)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get upload directory path"})
		return
	}

	fullPath := filepath.Join(absPath, filename)
	relativePath := filepath.Join(UPLOAD_DIR, filename)

	if err := saveFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	if err := database.AddMedia(id, c.GetString("username"), relativePath, header.Size); err != nil {
		os.Remove(fullPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"path": relativePath,
		"size": header.Size,
	})
}

// GetMedia retrieves and serves the media file
func GetMedia(c *gin.Context) {
	id := c.Param("id")

	media, err := database.GetMedia(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	absPath, err := filepath.Abs(media.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file path"})
		return
	}

	_, err = os.Stat(absPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found on disk"})
		return
	}

	c.File(absPath)
}

// DeleteMedia removes media from both database and disk
func DeleteMedia(c *gin.Context) {
	id := c.Param("id")

	media, err := database.GetMedia(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	absPath, err := filepath.Abs(media.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file path"})
		return
	}

	if err := database.DeleteMedia(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from database"})
		return
	}

	if err := os.Remove(absPath); err != nil {
		fmt.Printf("Failed to delete file %s: %v\n", absPath, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}
