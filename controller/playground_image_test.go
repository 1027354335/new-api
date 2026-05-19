package controller

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/setting/storage_setting"
	"github.com/gin-gonic/gin"
)

func TestUploadPlaygroundImage_Disabled(t *testing.T) {
	// Disable storage
	cfg := storage_setting.GetStorageSetting()
	cfg.Enabled = false

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create multipart request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.png")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	_, _ = part.Write([]byte("dummy image data"))
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "/api/playground/images/upload", body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	c.Request = req
	c.Set("id", 1)

	UploadPlaygroundImage(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
