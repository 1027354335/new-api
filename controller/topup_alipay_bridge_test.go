package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/pkg/alipaybridge"
	"github.com/QuantumNous/new-api/setting"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupAlipayBridgeControllerTestDB(t *testing.T) {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	model.DB = db
	model.LOG_DB = db
	common.UsingSQLite = true
	common.UsingMySQL = false
	common.UsingPostgreSQL = false
	common.RedisEnabled = false
	common.BatchUpdateEnabled = false

	require.NoError(t, db.AutoMigrate(&model.User{}, &model.TopUp{}, &model.Log{}))
	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})
}

func configureAlipayBridgeTest(t *testing.T) {
	t.Helper()
	originalEnabled := setting.AlipayBridgeEnabled
	originalSecret := setting.AlipayBridgeSecret
	originalQuotaPerUnit := common.QuotaPerUnit
	t.Cleanup(func() {
		setting.AlipayBridgeEnabled = originalEnabled
		setting.AlipayBridgeSecret = originalSecret
		common.QuotaPerUnit = originalQuotaPerUnit
	})
	setting.AlipayBridgeEnabled = true
	setting.AlipayBridgeSecret = "bridge-secret"
	common.QuotaPerUnit = 500000
}

func insertAlipayBridgeTopUp(t *testing.T, tradeNo string, provider string, status string) {
	t.Helper()
	require.NoError(t, model.DB.Create(&model.User{
		Id:       901,
		Username: "alipay_bridge_user",
		Status:   common.UserStatusEnabled,
		Quota:    0,
	}).Error)
	require.NoError(t, model.DB.Create(&model.TopUp{
		UserId:          901,
		Amount:          1,
		Money:           7.20,
		CreditAmountUsd: 1,
		PaidAmount:      7.20,
		PaidCurrency:    "CNY",
		ExchangeRate:    7.2,
		TradeNo:         tradeNo,
		PaymentMethod:   model.PaymentMethodAlipay,
		PaymentProvider: provider,
		CreateTime:      time.Now().Unix(),
		Status:          status,
	}).Error)
}

func signedAlipayBridgeSettleRequest(t *testing.T, nonce string, body []byte) *http.Request {
	t.Helper()
	path := "/api/alipay/bridge/settle"
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature, err := alipaybridge.Sign(setting.AlipayBridgeSecret, http.MethodPost, path, timestamp, nonce, body)
	require.NoError(t, err)
	req.Header.Set(alipaybridge.HeaderTimestamp, timestamp)
	req.Header.Set(alipaybridge.HeaderNonce, nonce)
	req.Header.Set(alipaybridge.HeaderSignature, signature)
	return req
}

func performAlipayBridgeSettle(t *testing.T, nonce string, payload alipayBridgeSettleRequest) *httptest.ResponseRecorder {
	t.Helper()
	body, err := common.Marshal(payload)
	require.NoError(t, err)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = signedAlipayBridgeSettleRequest(t, nonce, body)
	AlipayBridgeSettle(c)
	return recorder
}

func getAlipayBridgeUserQuota(t *testing.T) int {
	t.Helper()
	var user model.User
	require.NoError(t, model.DB.Where("id = ?", 901).First(&user).Error)
	return user.Quota
}

func TestAlipayBridgeSettleRechargesOnce(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAlipayBridgeControllerTestDB(t)
	configureAlipayBridgeTest(t)
	insertAlipayBridgeTopUp(t, "bridge-once", model.PaymentProviderAlipay, common.TopUpStatusPending)

	payload := alipayBridgeSettleRequest{
		TradeNo:     "bridge-once",
		TradeStatus: "TRADE_SUCCESS",
		TotalAmount: "7.20",
		Currency:    "CNY",
	}
	first := performAlipayBridgeSettle(t, "nonce-once-1", payload)
	require.Equal(t, http.StatusOK, first.Code)
	require.Equal(t, 500000, getAlipayBridgeUserQuota(t))

	second := performAlipayBridgeSettle(t, "nonce-once-2", payload)
	require.Equal(t, http.StatusOK, second.Code)
	require.Equal(t, 500000, getAlipayBridgeUserQuota(t))
}

func TestAlipayBridgeSettleRejectsAmountMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAlipayBridgeControllerTestDB(t)
	configureAlipayBridgeTest(t)
	insertAlipayBridgeTopUp(t, "bridge-amount", model.PaymentProviderAlipay, common.TopUpStatusPending)

	recorder := performAlipayBridgeSettle(t, "nonce-amount", alipayBridgeSettleRequest{
		TradeNo:     "bridge-amount",
		TradeStatus: "TRADE_SUCCESS",
		TotalAmount: "0.01",
		Currency:    "CNY",
	})
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Equal(t, 0, getAlipayBridgeUserQuota(t))
}

func TestAlipayBridgeSettleRejectsPaymentProviderMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAlipayBridgeControllerTestDB(t)
	configureAlipayBridgeTest(t)
	insertAlipayBridgeTopUp(t, "bridge-provider", model.PaymentProviderStripe, common.TopUpStatusPending)

	recorder := performAlipayBridgeSettle(t, "nonce-provider", alipayBridgeSettleRequest{
		TradeNo:     "bridge-provider",
		TradeStatus: "TRADE_SUCCESS",
		TotalAmount: "7.20",
		Currency:    "CNY",
	})
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Equal(t, 0, getAlipayBridgeUserQuota(t))
}

func TestAlipayBridgeSettleRejectsInvalidStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupAlipayBridgeControllerTestDB(t)
	configureAlipayBridgeTest(t)
	insertAlipayBridgeTopUp(t, "bridge-failed", model.PaymentProviderAlipay, common.TopUpStatusFailed)

	recorder := performAlipayBridgeSettle(t, "nonce-status", alipayBridgeSettleRequest{
		TradeNo:     "bridge-failed",
		TradeStatus: "TRADE_SUCCESS",
		TotalAmount: "7.20",
		Currency:    "CNY",
	})
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Equal(t, 0, getAlipayBridgeUserQuota(t))
}
