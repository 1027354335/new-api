package controller

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting/storage_setting"
	"github.com/gin-gonic/gin"
)

// GetPlaygroundImage proxies requests for images stored in MinIO.
func GetPlaygroundImage(c *gin.Context) {
	// The path parameter from *path starts with a leading slash, e.g. "/[Bucket]/1/uuid.png"
	pathParam := c.Param("path")
	objectName := strings.TrimPrefix(pathParam, "/")

	// Validate path structure: [Bucket]/{userID}/{filename}
	parts := strings.Split(objectName, "/")
	if len(parts) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image path format"})
		return
	}

	targetUserID, err := strconv.Atoi(parts[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in path"})
		return
	}

	// Security check: only allow owners or admins to retrieve the image
	currentUserID := c.GetInt("id")
	currentUserRole := c.GetInt("role")

	if currentUserRole < common.RoleAdminUser && currentUserID != targetUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this image"})
		return
	}

	// The objectName inside MinIO bucket is the path after the bucket name segment (parts[0])
	actualObjectName := strings.Join(parts[1:], "/")

	// Fetch image stream from MinIO
	reader, contentLength, contentType, err := service.GetPlaygroundImageReader(c.Request.Context(), actualObjectName)
	if err != nil {
		// Fallback: try fetching with the full objectName (including bucket name segment)
		// in case of legacy objects stored with "playground/" prefix key
		var err2 error
		reader, contentLength, contentType, err2 = service.GetPlaygroundImageReader(c.Request.Context(), objectName)
		if err2 != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
	}
	defer reader.Close()

	// Stream to client
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))
	c.Header("Cache-Control", "private, max-age=31536000") // Cache locally for privacy-isolated requests

	c.Writer.WriteHeader(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}

// TestStorageConnection tests the configured storage settings connection
func TestStorageConnection(c *gin.Context) {
	// Only accessible by admins (routes should enforce middleware.AdminAuth(), but let's be safe)
	role := c.GetInt("role")
	if role < common.RoleAdminUser {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Forbidden"})
		return
	}

	err := storage_setting.TestConnection()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Connection tested successfully!",
	})
}

// UploadPlaygroundImage handles POST /api/playground/images/upload.
func UploadPlaygroundImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}
	defer file.Close()

	// Read all file bytes
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Verify size limits (e.g. 10MB)
	if len(data) > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
		return
	}

	// Detect content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	// Base64 encode it and call UploadPlaygroundImageFromBase64
	base64Str := base64.StdEncoding.EncodeToString(data)
	base64DataUri := fmt.Sprintf("data:%s;base64,%s", contentType, base64Str)

	userID := c.GetInt("id")
	relativePath, err := service.UploadPlaygroundImageFromBase64(c.Request.Context(), userID, base64DataUri)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format to complete URL
	fullURL := service.FormatPlaygroundImageURL(relativePath)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"url":     fullURL,
	})
}
