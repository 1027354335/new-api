package storage_setting

import (
	"testing"
)

func TestGetStorageSetting(t *testing.T) {
	setting := GetStorageSetting()
	if setting == nil {
		t.Fatal("expected non-nil storage setting")
	}

	// Default value should be disabled
	if setting.Enabled {
		t.Error("expected default enabled to be false")
	}
}

func TestUpdateAndSync(t *testing.T) {
	setting := GetStorageSetting()
	setting.Enabled = true
	setting.Endpoint = "localhost:9000"
	setting.AccessKey = "testkey"
	setting.SecretKey = "testsecret"
	setting.Bucket = "testbucket"
	setting.UseSSL = false

	UpdateAndSync()

	clientMu.RLock()
	isNil := minioClient == nil
	clientMu.RUnlock()

	if !isNil {
		t.Error("expected minioClient to be nil after UpdateAndSync")
	}
}

func TestGetClientDisabled(t *testing.T) {
	setting := GetStorageSetting()
	setting.Enabled = false

	UpdateAndSync()

	_, err := GetClient()
	if err == nil {
		t.Error("expected error when storage setting is disabled")
	}
}
