package alipaybridge

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	HeaderTimestamp = "X-Alipay-Bridge-Timestamp"
	HeaderNonce     = "X-Alipay-Bridge-Nonce"
	HeaderSignature = "X-Alipay-Bridge-Signature"
)

var (
	ErrMissingSecret    = errors.New("missing alipay bridge secret")
	ErrMissingSignature = errors.New("missing alipay bridge signature headers")
	ErrExpiredSignature = errors.New("expired alipay bridge signature")
	ErrInvalidSignature = errors.New("invalid alipay bridge signature")
)

func BodySHA256(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func CanonicalString(method string, path string, timestamp string, nonce string, body []byte) string {
	return strings.ToUpper(method) + "\n" + path + "\n" + timestamp + "\n" + nonce + "\n" + BodySHA256(body)
}

func Sign(secret string, method string, path string, timestamp string, nonce string, body []byte) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", ErrMissingSecret
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(CanonicalString(method, path, timestamp, nonce, body)))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func Verify(secret string, method string, path string, timestamp string, nonce string, signature string, body []byte, now time.Time, window time.Duration) error {
	if strings.TrimSpace(secret) == "" {
		return ErrMissingSecret
	}
	if timestamp == "" || nonce == "" || signature == "" {
		return ErrMissingSignature
	}
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: bad timestamp", ErrMissingSignature)
	}
	signedAt := time.Unix(ts, 0)
	if signedAt.Before(now.Add(-window)) || signedAt.After(now.Add(window)) {
		return ErrExpiredSignature
	}
	expected, err := Sign(secret, method, path, timestamp, nonce, body)
	if err != nil {
		return err
	}
	if !hmac.Equal([]byte(expected), []byte(strings.ToLower(signature))) {
		return ErrInvalidSignature
	}
	return nil
}
