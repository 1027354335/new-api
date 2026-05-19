package storage_setting

import (
	"context"
	"fmt"
	"sync"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageSetting struct {
	Enabled   bool   `json:"enabled"`
	Endpoint  string `json:"endpoint"`   // MinIO endpoint (e.g. play.minio.io:9000)
	Bucket    string `json:"bucket"`     // Bucket name
	AccessKey string `json:"access_key"` // Access Key
	SecretKey string `json:"secret_key"` // Secret Key
	UseSSL    bool   `json:"use_ssl"`    // Use HTTPS/SSL
	Region    string `json:"region"`     // Bucket Region (optional)
	UrlPrefix string `json:"url_prefix"` // Custom public access URL prefix (optional)
}

var storageSetting = StorageSetting{
	Enabled:   false,
	Endpoint:  "",
	Bucket:    "",
	AccessKey: "",
	SecretKey: "",
	UseSSL:    false,
	Region:    "",
	UrlPrefix: "",
}

var (
	minioClient *minio.Client
	clientMu    sync.RWMutex
)

func init() {
	// Register to the global config manager
	config.GlobalConfig.Register("storage_setting", &storageSetting)
}

func GetStorageSetting() *StorageSetting {
	return &storageSetting
}

func UpdateAndSync() {
	clientMu.Lock()
	defer clientMu.Unlock()
	minioClient = nil // Reset client to force reinitialization on next use
}

// GetClient returns a cached MinIO client instance or initializes a new one.
func GetClient() (*minio.Client, error) {
	clientMu.RLock()
	if minioClient != nil {
		defer clientMu.RUnlock()
		return minioClient, nil
	}
	clientMu.RUnlock()

	clientMu.Lock()
	defer clientMu.Unlock()

	// Double-check
	if minioClient != nil {
		return minioClient, nil
	}

	if !storageSetting.Enabled {
		return nil, fmt.Errorf("storage setting is disabled")
	}

	if storageSetting.Endpoint == "" || storageSetting.AccessKey == "" || storageSetting.SecretKey == "" {
		return nil, fmt.Errorf("storage setting is incomplete")
	}

	client, err := minio.New(storageSetting.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(storageSetting.AccessKey, storageSetting.SecretKey, ""),
		Secure: storageSetting.UseSSL,
		Region: storageSetting.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	minioClient = client
	return minioClient, nil
}

// TestConnection tests if the storage setting can connect to MinIO.
func TestConnection() error {
	if storageSetting.Endpoint == "" || storageSetting.AccessKey == "" || storageSetting.SecretKey == "" {
		return fmt.Errorf("storage config is incomplete")
	}

	client, err := minio.New(storageSetting.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(storageSetting.AccessKey, storageSetting.SecretKey, ""),
		Secure: storageSetting.UseSSL,
		Region: storageSetting.Region,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize test minio client: %w", err)
	}

	// Try listing buckets or checking bucket existence to verify access
	ctx := context.Background()
	if storageSetting.Bucket != "" {
		exists, err := client.BucketExists(ctx, storageSetting.Bucket)
		if err != nil {
			return fmt.Errorf("failed to check bucket existence: %w", err)
		}
		if !exists {
			// Try to create the bucket
			err = client.MakeBucket(ctx, storageSetting.Bucket, minio.MakeBucketOptions{Region: storageSetting.Region})
			if err != nil {
				return fmt.Errorf("bucket does not exist and failed to create: %w", err)
			}
			common.SysLog(fmt.Sprintf("successfully created bucket %s", storageSetting.Bucket))
		}
	} else {
		_, err = client.ListBuckets(ctx)
		if err != nil {
			return fmt.Errorf("failed to list buckets: %w", err)
		}
	}

	return nil
}
