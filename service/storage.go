package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/storage_setting"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// UploadPlaygroundImageFromURL downloads an image from originUrl and uploads it to MinIO, returning the relative path.
func UploadPlaygroundImageFromURL(ctx context.Context, userID int, originUrl string) (string, error) {
	cfg := storage_setting.GetStorageSetting()
	if !cfg.Enabled {
		return "", fmt.Errorf("storage setting is disabled")
	}

	client, err := storage_setting.GetClient()
	if err != nil {
		return "", err
	}

	// 1. Download image
	resp, err := DoDownloadRequest(originUrl, "upload to minio")
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read downloaded image body: %w", err)
	}

	// 2. Determine content type and file extension
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	// Sanitize content type (remove parameters like charset)
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = contentType[:idx]
	}

	ext := ".png" // default extension
	exts, err := mime.ExtensionsByType(contentType)
	if err == nil && len(exts) > 0 {
		ext = exts[0]
	}

	// 3. Ensure bucket exists
	err = ensureBucketExists(ctx, client, cfg.Bucket, cfg.Region)
	if err != nil {
		return "", err
	}

	// 4. Generate object path: {userID}/{uuid}{ext} (no bucket name or "playground/" prefix in key)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	objectName := fmt.Sprintf("%d/%s", userID, fileName)

	// 5. Upload to MinIO
	reader := bytes.NewReader(data)
	_, err = client.PutObject(ctx, cfg.Bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to minio: %w", err)
	}

	// 6. Return relative path
	return objectName, nil
}

// UploadPlaygroundImageFromBase64 uploads a base64 encoded image to MinIO, returning the relative path.
func UploadPlaygroundImageFromBase64(ctx context.Context, userID int, base64Data string) (string, error) {
	cfg := storage_setting.GetStorageSetting()
	if !cfg.Enabled {
		return "", fmt.Errorf("storage setting is disabled")
	}

	client, err := storage_setting.GetClient()
	if err != nil {
		return "", err
	}

	// Remove data URI prefix if present (e.g. "data:image/png;base64,")
	var contentType string
	if idx := strings.Index(base64Data, ","); idx != -1 {
		prefix := base64Data[:idx]
		base64Data = base64Data[idx+1:]
		if strings.HasPrefix(prefix, "data:") {
			contentType = prefix[5:strings.Index(prefix, ";")]
		}
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	ext := ".png"
	exts, err := mime.ExtensionsByType(contentType)
	if err == nil && len(exts) > 0 {
		ext = exts[0]
	}

	// Ensure bucket exists
	err = ensureBucketExists(ctx, client, cfg.Bucket, cfg.Region)
	if err != nil {
		return "", err
	}

	// Generate object path: {userID}/{uuid}{ext}
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	objectName := fmt.Sprintf("%d/%s", userID, fileName)

	reader := bytes.NewReader(data)
	_, err = client.PutObject(ctx, cfg.Bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload base64 to minio: %w", err)
	}

	return objectName, nil
}

// GetPlaygroundImageReader fetches an object from MinIO and returns a Reader, Content-Length, Content-Type, and error.
func GetPlaygroundImageReader(ctx context.Context, objectName string) (io.ReadCloser, int64, string, error) {
	cfg := storage_setting.GetStorageSetting()
	if !cfg.Enabled {
		return nil, 0, "", fmt.Errorf("storage setting is disabled")
	}

	client, err := storage_setting.GetClient()
	if err != nil {
		return nil, 0, "", err
	}

	// Get object info first for metadata
	info, err := client.StatObject(ctx, cfg.Bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to stat object: %w", err)
	}

	contentType := info.ContentType
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(objectName))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	obj, err := client.GetObject(ctx, cfg.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to get object: %w", err)
	}

	return obj, info.Size, contentType, nil
}

func ensureBucketExists(ctx context.Context, client *minio.Client, bucketName, region string) error {
	// Cache bucket check to avoid redundant API calls
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket %s exists: %w", bucketName, err)
	}
	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
		if err != nil {
			return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
		common.SysLog(fmt.Sprintf("created bucket %s in region %s", bucketName, region))
	}
	return nil
}

// FormatPlaygroundImageURL prepends the prefix and bucket name to the relative path
func FormatPlaygroundImageURL(relativePath string) string {
	if relativePath == "" {
		return ""
	}
	// If it's already an absolute or API URL, leave it alone
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") || strings.HasPrefix(relativePath, "/") {
		return relativePath
	}

	cfg := storage_setting.GetStorageSetting()
	if !cfg.Enabled {
		return relativePath
	}

	prefix := cfg.UrlPrefix
	if prefix == "" {
		prefix = "/api/playground/images"
	}
	prefix = strings.TrimSuffix(prefix, "/")
	bucket := strings.Trim(cfg.Bucket, "/")
	relativePath = strings.TrimPrefix(relativePath, "/")

	return fmt.Sprintf("%s/%s/%s", prefix, bucket, relativePath)
}

// StripPlaygroundImageURL removes the prefix and bucket name from the URL, leaving only the relative path
func StripPlaygroundImageURL(urlStr string) string {
	if urlStr == "" {
		return ""
	}

	cfg := storage_setting.GetStorageSetting()
	if !cfg.Enabled {
		return urlStr
	}

	bucket := strings.Trim(cfg.Bucket, "/")

	// 1. Try stripping custom prefix + bucket
	if cfg.UrlPrefix != "" {
		prefix := strings.TrimSuffix(cfg.UrlPrefix, "/")
		prefixWithBucket := fmt.Sprintf("%s/%s/", prefix, bucket)
		if strings.HasPrefix(urlStr, prefixWithBucket) {
			return strings.TrimPrefix(urlStr, prefixWithBucket)
		}
	}

	// 2. Try stripping default prefix + bucket
	defaultPrefixWithBucket := fmt.Sprintf("/api/playground/images/%s/", bucket)
	if strings.HasPrefix(urlStr, defaultPrefixWithBucket) {
		return strings.TrimPrefix(urlStr, defaultPrefixWithBucket)
	}

	// 3. Try stripping legacy prefix + bucket
	legacyPrefixWithBucket := "/api/playground/images/playground/"
	if strings.HasPrefix(urlStr, legacyPrefixWithBucket) {
		return strings.TrimPrefix(urlStr, legacyPrefixWithBucket)
	}

	return urlStr
}

// ProcessPlaygroundJSONUrls walks a JSONValue and applies a string transformer to all image URLs
func ProcessPlaygroundJSONUrls(jsonVal model.JSONValue, transform func(string) string) model.JSONValue {
	if len(jsonVal) == 0 {
		return jsonVal
	}

	var data any
	err := common.Unmarshal(jsonVal, &data)
	if err != nil {
		return jsonVal
	}

	data = transformJSONUrls(data, transform)

	bytes, err := common.Marshal(data)
	if err != nil {
		return jsonVal
	}

	return model.JSONValue(bytes)
}

func transformJSONUrls(val any, transform func(string) string) any {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case map[string]any:
		// If it's a map containing "url", we transform the value of "url"
		if urlVal, ok := v["url"]; ok {
			if urlStr, ok := urlVal.(string); ok {
				v["url"] = transform(urlStr)
			}
		}
		// Recursively process map values
		for k, item := range v {
			v[k] = transformJSONUrls(item, transform)
		}
		return v
	case []any:
		for i, item := range v {
			v[i] = transformJSONUrls(item, transform)
		}
		return v
	}
	return val
}

// ConvertPlaygroundURLToBase64 checks if the URL is a playground image and converts it to a base64 data URI by fetching from MinIO.
func ConvertPlaygroundURLToBase64(ctx context.Context, urlStr string) (string, error) {
	if urlStr == "" {
		return "", nil
	}

	// If it's already a base64 data URI, return as-is
	if strings.HasPrefix(urlStr, "data:") {
		return urlStr, nil
	}

	// Extract the relative path / object name from the URL
	stripped := StripPlaygroundImageURL(urlStr)

	// If it was not stripped (i.e. not a playground image URL)
	if stripped == urlStr && !strings.Contains(urlStr, "/api/playground/images/") {
		return urlStr, nil
	}

	// Clean up any double slashes or leading/trailing slashes
	stripped = strings.TrimPrefix(stripped, "/")

	// Fetch from MinIO
	reader, _, contentType, err := GetPlaygroundImageReader(ctx, stripped)
	if err != nil {
		return "", fmt.Errorf("failed to get image reader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read image bytes: %w", err)
	}

	base64Str := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Str), nil
}

// NormalizeImagesField parses Images raw JSON, converts any playground URLs to base64, and formats to object array for upstream
func NormalizeImagesField(ctx context.Context, raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return raw, nil
	}
	var arr []any
	// Try unmarshalling as array
	err := common.Unmarshal(raw, &arr)
	if err != nil {
		// Not an array, maybe a single item?
		var singleItem any
		if err2 := common.Unmarshal(raw, &singleItem); err2 == nil {
			arr = []any{singleItem}
		} else {
			return raw, err
		}
	}

	resultArr := make([]map[string]any, 0, len(arr))
	for _, item := range arr {
		if item == nil {
			continue
		}
		switch v := item.(type) {
		case string:
			resolved, err := ConvertPlaygroundURLToBase64(ctx, v)
			if err != nil {
				common.SysError("failed to resolve playground URL: " + err.Error())
				resolved = v // fallback
			}
			resultArr = append(resultArr, map[string]any{
				"image_url": resolved,
			})
		case map[string]any:
			if imgUrlVal, ok := v["image_url"]; ok {
				switch u := imgUrlVal.(type) {
				case string:
					resolved, err := ConvertPlaygroundURLToBase64(ctx, u)
					if err != nil {
						common.SysError("failed to resolve playground URL: " + err.Error())
						resolved = u
					}
					v["image_url"] = resolved
				case map[string]any:
					if urlVal, ok := u["url"]; ok {
						if urlStr, ok := urlVal.(string); ok {
							resolved, err := ConvertPlaygroundURLToBase64(ctx, urlStr)
							if err != nil {
								common.SysError("failed to resolve playground URL: " + err.Error())
								resolved = urlStr
							}
							u["url"] = resolved
						}
					}
				}
			} else if urlVal, ok := v["url"]; ok {
				if urlStr, ok := urlVal.(string); ok {
					resolved, err := ConvertPlaygroundURLToBase64(ctx, urlStr)
					if err != nil {
						common.SysError("failed to resolve playground URL: " + err.Error())
						resolved = urlStr
					}
					v["url"] = resolved
				}
			}
			resultArr = append(resultArr, v)
		}
	}

	bytes, err := common.Marshal(resultArr)
	if err != nil {
		return raw, err
	}
	return json.RawMessage(bytes), nil
}

// NormalizeMaskField parses Mask raw JSON, converts any playground URLs to base64, and formats to object for upstream
func NormalizeMaskField(ctx context.Context, raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return raw, nil
	}
	var item any
	err := common.Unmarshal(raw, &item)
	if err != nil {
		return raw, err
	}

	switch v := item.(type) {
	case string:
		resolved, err := ConvertPlaygroundURLToBase64(ctx, v)
		if err != nil {
			common.SysError("failed to resolve playground URL for mask: " + err.Error())
			resolved = v
		}
		result := map[string]any{
			"image_url": resolved,
		}
		bytes, err := common.Marshal(result)
		if err != nil {
			return raw, err
		}
		return json.RawMessage(bytes), nil
	case map[string]any:
		if imgUrlVal, ok := v["image_url"]; ok {
			switch u := imgUrlVal.(type) {
			case string:
				resolved, err := ConvertPlaygroundURLToBase64(ctx, u)
				if err != nil {
					common.SysError("failed to resolve playground URL for mask: " + err.Error())
					resolved = u
				}
				v["image_url"] = resolved
			case map[string]any:
				if urlVal, ok := u["url"]; ok {
					if urlStr, ok := urlVal.(string); ok {
						resolved, err := ConvertPlaygroundURLToBase64(ctx, urlStr)
						if err != nil {
							common.SysError("failed to resolve playground URL for mask: " + err.Error())
							resolved = urlStr
						}
						u["url"] = resolved
					}
				}
			}
		} else if urlVal, ok := v["url"]; ok {
			if urlStr, ok := urlVal.(string); ok {
				resolved, err := ConvertPlaygroundURLToBase64(ctx, urlStr)
				if err != nil {
					common.SysError("failed to resolve playground URL for mask: " + err.Error())
					resolved = urlStr
				}
				v["url"] = resolved
			}
		}
		bytes, err := common.Marshal(v)
		if err != nil {
			return raw, err
		}
		return json.RawMessage(bytes), nil
	}

	return raw, nil
}
