package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/dto"
	relayconstant "github.com/QuantumNous/new-api/relay/constant"
	"github.com/QuantumNous/new-api/setting/storage_setting"
	"github.com/gin-gonic/gin"
)

func TestGetAndValidOpenAIImageRequest_Interception(t *testing.T) {
	// Disable storage setting in tests so it falls back to original URLs
	cfg := storage_setting.GetStorageSetting()
	cfg.Enabled = false

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Create a dummy JSON request body with images as a string array
	reqBody := `{"model": "dall-e-3", "prompt": "a blue sky", "images": ["/api/playground/images/playground/1/uuid.png"]}`

	req, err := http.NewRequest(http.MethodPost, "/v1/images/generations", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Call the target function
	imageReq, err := GetAndValidOpenAIImageRequest(c, relayconstant.RelayModeImagesGenerations)
	if err != nil {
		t.Fatalf("GetAndValidOpenAIImageRequest failed: %v", err)
	}

	// Verify the in-memory request struct is normalized
	if imageReq.Images == nil {
		t.Fatal("Expected images to be populated, got nil")
	}

	var imagesArr []map[string]any
	err = json.Unmarshal(imageReq.Images, &imagesArr)
	if err != nil {
		t.Fatalf("Failed to unmarshal images: %v", err)
	}

	if len(imagesArr) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(imagesArr))
	}

	if imagesArr[0]["image_url"] != "/api/playground/images/playground/1/uuid.png" {
		t.Errorf("Expected image_url to be normalized, got %v", imagesArr[0]["image_url"])
	}

	// Verify the BodyStorage cache in Gin Context has been updated
	storage, err := common.GetBodyStorage(c)
	if err != nil {
		t.Fatalf("Failed to get BodyStorage: %v", err)
	}

	bodyBytes, err := storage.Bytes()
	if err != nil {
		t.Fatalf("Failed to read body bytes: %v", err)
	}

	var parsedBody dto.ImageRequest
	err = json.Unmarshal(bodyBytes, &parsedBody)
	if err != nil {
		t.Fatalf("Failed to unmarshal rewritten body: %v", err)
	}

	var rewrittenImagesArr []map[string]any
	err = json.Unmarshal(parsedBody.Images, &rewrittenImagesArr)
	if err != nil {
		t.Fatalf("Failed to unmarshal rewritten images: %v", err)
	}

	if len(rewrittenImagesArr) != 1 {
		t.Fatalf("Expected 1 item in rewritten body, got %d", len(rewrittenImagesArr))
	}

	if rewrittenImagesArr[0]["image_url"] != "/api/playground/images/playground/1/uuid.png" {
		t.Errorf("Expected rewritten image_url to be original URL, got %v", rewrittenImagesArr[0]["image_url"])
	}

	// Verify c.Request.Body is readable and has rewritten content
	rewrittenBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		t.Fatalf("Failed to read c.Request.Body: %v", err)
	}

	var parsedReqBody dto.ImageRequest
	err = json.Unmarshal(rewrittenBytes, &parsedReqBody)
	if err != nil {
		t.Fatalf("Failed to parse request body: %v", err)
	}
}

func TestGetAndValidOpenAIImageRequest_InputFidelity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		model             string
		inputFidelity     string
		expectHasFidelity bool
	}{
		{
			model:             "gpt-image-2",
			inputFidelity:     "high",
			expectHasFidelity: false,
		},
		{
			model:             "gpt-image-1.5",
			inputFidelity:     "high",
			expectHasFidelity: true,
		},
		{
			model:             "gpt-image-1.5-pro",
			inputFidelity:     "high",
			expectHasFidelity: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.model, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			reqBody := `{"model": "` + tc.model + `", "prompt": "a star", "input_fidelity": "` + tc.inputFidelity + `"}`
			req, err := http.NewRequest(http.MethodPost, "/v1/images/generations", bytes.NewBufferString(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			imageReq, err := GetAndValidOpenAIImageRequest(c, relayconstant.RelayModeImagesGenerations)
			if err != nil {
				t.Fatalf("GetAndValidOpenAIImageRequest failed: %v", err)
			}

			hasFidelity := len(imageReq.InputFidelity) > 0
			if hasFidelity != tc.expectHasFidelity {
				t.Errorf("For model %s, expected hasFidelity=%v, got %v", tc.model, tc.expectHasFidelity, hasFidelity)
			}

			// Verify in rewritten body storage too
			storage, err := common.GetBodyStorage(c)
			if err == nil {
				bodyBytes, _ := storage.Bytes()
				var parsed map[string]any
				_ = json.Unmarshal(bodyBytes, &parsed)
				_, exists := parsed["input_fidelity"]
				if exists != tc.expectHasFidelity {
					t.Errorf("For model %s in rewritten body, expected input_fidelity exists=%v, got %v", tc.model, tc.expectHasFidelity, exists)
				}
			}
		})
	}
}
