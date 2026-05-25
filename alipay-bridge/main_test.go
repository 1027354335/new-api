package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/pkg/alipaybridge"
	"github.com/stretchr/testify/require"
)

func TestValidateCreateRequest(t *testing.T) {
	require.NoError(t, validateCreateRequest(&createRequest{
		TradeNo:     "ref_123",
		Subject:     "Topup",
		TotalAmount: "7.20",
		Currency:    "CNY",
	}))
	require.Error(t, validateCreateRequest(&createRequest{
		TradeNo:     "ref_123",
		Subject:     "Topup",
		TotalAmount: "0",
		Currency:    "CNY",
	}))
}

func TestReadAndVerifyBridgeRequestRejectsMissingSignature(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/alipay/create", bytes.NewReader([]byte(`{}`)))
	recorder := httptest.NewRecorder()
	_, ok := readAndVerifyBridgeRequest(recorder, req, "secret")
	require.False(t, ok)
	require.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestReadAndVerifyBridgeRequestAcceptsSignedRequest(t *testing.T) {
	body, err := json.Marshal(createRequest{
		TradeNo:     "ref_signed",
		Subject:     "Topup",
		TotalAmount: "7.20",
		Currency:    "CNY",
	})
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/api/alipay/create", bytes.NewReader(body))
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature, err := alipaybridge.Sign("secret", http.MethodPost, "/api/alipay/create", timestamp, "bridge-test-nonce", body)
	require.NoError(t, err)
	req.Header.Set(alipaybridge.HeaderTimestamp, timestamp)
	req.Header.Set(alipaybridge.HeaderNonce, "bridge-test-nonce")
	req.Header.Set(alipaybridge.HeaderSignature, signature)

	recorder := httptest.NewRecorder()
	got, ok := readAndVerifyBridgeRequest(recorder, req, "secret")
	require.True(t, ok)
	require.Equal(t, body, got)
}

func TestLoadEnvFileSetsMissingValues(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte(`
# comment
ALIPAY_APP_ID=app_from_file
ALIPAY_BRIDGE_SECRET="secret from file"
export ALIPAY_SANDBOX=true
`), 0600))
	t.Setenv("ALIPAY_APP_ID", "")
	t.Setenv("ALIPAY_BRIDGE_SECRET", "")
	t.Setenv("ALIPAY_SANDBOX", "")

	require.NoError(t, loadEnvFile(envPath))
	require.Equal(t, "app_from_file", os.Getenv("ALIPAY_APP_ID"))
	require.Equal(t, "secret from file", os.Getenv("ALIPAY_BRIDGE_SECRET"))
	require.Equal(t, "true", os.Getenv("ALIPAY_SANDBOX"))
}

func TestLoadEnvFileKeepsExistingEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte(`ALIPAY_APP_ID=app_from_file`), 0600))
	t.Setenv("ALIPAY_APP_ID", "app_from_env")

	require.NoError(t, loadEnvFile(envPath))
	require.Equal(t, "app_from_env", os.Getenv("ALIPAY_APP_ID"))
}

func TestParseEnvLine(t *testing.T) {
	key, value, ok := parseEnvLine(`export ALIPAY_BRIDGE_SECRET='abc123'`)
	require.True(t, ok)
	require.Equal(t, "ALIPAY_BRIDGE_SECRET", key)
	require.Equal(t, "abc123", value)

	_, _, ok = parseEnvLine("# ignored")
	require.False(t, ok)
}
