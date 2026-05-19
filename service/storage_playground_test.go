package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/QuantumNous/new-api/setting/storage_setting"
)

func TestNormalizeImagesField_StringArray(t *testing.T) {
	// Setup storage setting to disabled so it falls back to original URLs
	cfg := storage_setting.GetStorageSetting()
	cfg.Enabled = false

	input := json.RawMessage(`["/api/playground/images/playground/1/uuid.png", "https://example.com/other.png"]`)
	normalized, err := NormalizeImagesField(context.Background(), input)
	if err != nil {
		t.Fatalf("NormalizeImagesField failed: %v", err)
	}

	var result []map[string]any
	err = json.Unmarshal(normalized, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(result))
	}

	if result[0]["image_url"] != "/api/playground/images/playground/1/uuid.png" {
		t.Errorf("Expected first image_url to be original URL, got %v", result[0]["image_url"])
	}

	if result[1]["image_url"] != "https://example.com/other.png" {
		t.Errorf("Expected second image_url to be original URL, got %v", result[1]["image_url"])
	}
}

func TestNormalizeImagesField_ObjectArray(t *testing.T) {
	cfg := storage_setting.GetStorageSetting()
	cfg.Enabled = false

	input := json.RawMessage(`[{"image_url": "/api/playground/images/playground/1/uuid.png"}, {"image_url": {"url": "https://example.com/other.png"}}]`)
	normalized, err := NormalizeImagesField(context.Background(), input)
	if err != nil {
		t.Fatalf("NormalizeImagesField failed: %v", err)
	}

	var result []map[string]any
	err = json.Unmarshal(normalized, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(result))
	}

	img1 := result[0]["image_url"]
	if img1 != "/api/playground/images/playground/1/uuid.png" {
		t.Errorf("Expected first image_url to be original, got %v", img1)
	}

	img2Map, ok := result[1]["image_url"].(map[string]any)
	if !ok {
		t.Fatalf("Expected second image_url to be map, got %T", result[1]["image_url"])
	}
	if img2Map["url"] != "https://example.com/other.png" {
		t.Errorf("Expected nested url to be original, got %v", img2Map["url"])
	}
}

func TestNormalizeMaskField(t *testing.T) {
	cfg := storage_setting.GetStorageSetting()
	cfg.Enabled = false

	// Test string input
	inputStr := json.RawMessage(`"/api/playground/images/playground/1/uuid.png"`)
	normalized, err := NormalizeMaskField(context.Background(), inputStr)
	if err != nil {
		t.Fatalf("NormalizeMaskField string failed: %v", err)
	}

	var result map[string]any
	err = json.Unmarshal(normalized, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal mask result: %v", err)
	}

	if result["image_url"] != "/api/playground/images/playground/1/uuid.png" {
		t.Errorf("Expected mask image_url to be original, got %v", result["image_url"])
	}

	// Test map input
	inputMap := json.RawMessage(`{"image_url": "/api/playground/images/playground/1/uuid.png"}`)
	normalizedMap, err := NormalizeMaskField(context.Background(), inputMap)
	if err != nil {
		t.Fatalf("NormalizeMaskField map failed: %v", err)
	}

	var resultMap map[string]any
	err = json.Unmarshal(normalizedMap, &resultMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal mask map result: %v", err)
	}

	if resultMap["image_url"] != "/api/playground/images/playground/1/uuid.png" {
		t.Errorf("Expected mask map image_url to be original, got %v", resultMap["image_url"])
	}
}
