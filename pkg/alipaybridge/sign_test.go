package alipaybridge

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSignVerify(t *testing.T) {
	body := []byte(`{"trade_no":"ref_123","amount":"1.00"}`)
	now := time.Unix(1700000000, 0)
	sig, err := Sign("secret", "POST", "/api/alipay/bridge/settle", "1700000000", "nonce-1", body)
	require.NoError(t, err)

	err = Verify("secret", "POST", "/api/alipay/bridge/settle", "1700000000", "nonce-1", sig, body, now, 5*time.Minute)
	require.NoError(t, err)
}

func TestVerifyRejectsTamperedBody(t *testing.T) {
	body := []byte(`{"trade_no":"ref_123","amount":"1.00"}`)
	now := time.Unix(1700000000, 0)
	sig, err := Sign("secret", "POST", "/api/alipay/bridge/settle", "1700000000", "nonce-1", body)
	require.NoError(t, err)

	err = Verify("secret", "POST", "/api/alipay/bridge/settle", "1700000000", "nonce-1", sig, []byte(`{"trade_no":"ref_123","amount":"9.00"}`), now, 5*time.Minute)
	require.ErrorIs(t, err, ErrInvalidSignature)
}

func TestVerifyRejectsExpiredTimestamp(t *testing.T) {
	body := []byte(`{}`)
	now := time.Unix(1700000600, 0)
	sig, err := Sign("secret", "POST", "/api/alipay/bridge/settle", "1700000000", "nonce-1", body)
	require.NoError(t, err)

	err = Verify("secret", "POST", "/api/alipay/bridge/settle", "1700000000", "nonce-1", sig, body, now, time.Minute)
	require.True(t, errors.Is(err, ErrExpiredSignature))
}
