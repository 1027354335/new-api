package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/pkg/alipaybridge"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/operation_setting"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/smartwalle/alipay/v3"
	"github.com/thanhpk/randstr"
)

type AlipayPayRequest struct {
	Amount        int64  `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}

type alipayBridgeCreateRequest struct {
	TradeNo     string `json:"trade_no"`
	Subject     string `json:"subject"`
	TotalAmount string `json:"total_amount"`
	Currency    string `json:"currency"`
	ReturnURL   string `json:"return_url"`
	NotifyURL   string `json:"notify_url"`
}

type alipayBridgeCreateResponse struct {
	Message string `json:"message"`
	Data    struct {
		ApproveLink string `json:"approve_link"`
		OrderID     string `json:"order_id"`
	} `json:"data"`
}

type alipayBridgeSettleRequest struct {
	TradeNo       string `json:"trade_no"`
	AlipayTradeNo string `json:"alipay_trade_no,omitempty"`
	TradeStatus   string `json:"trade_status"`
	TotalAmount   string `json:"total_amount"`
	Currency      string `json:"currency,omitempty"`
}

var alipayBridgeNonceCache sync.Map

func getAlipayClient() (*alipay.Client, error) {
	if strings.TrimSpace(setting.AlipayAppId) == "" || strings.TrimSpace(setting.AlipayPrivateKey) == "" || strings.TrimSpace(setting.AlipayPublicKey) == "" {
		return nil, fmt.Errorf("alipay credentials are not configured")
	}

	client, err := alipay.New(setting.AlipayAppId, setting.AlipayPrivateKey, !setting.AlipaySandbox)
	if err != nil {
		return nil, err
	}

	err = client.LoadAliPayPublicKey(setting.AlipayPublicKey)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getAlipayQuote(amount int64, group string) (float64, float64, float64) {
	dAmount := decimal.NewFromInt(amount)
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dQuotaPerUnit := decimal.NewFromFloat(common.QuotaPerUnit)
		dAmount = dAmount.Div(dQuotaPerUnit)
	}

	topupGroupRatio := common.GetTopupGroupRatio(group)
	if topupGroupRatio == 0 {
		topupGroupRatio = 1
	}

	discount := 1.0
	if ds, ok := operation_setting.GetPaymentSetting().AmountDiscount[int(amount)]; ok {
		if ds > 0 {
			discount = ds
		}
	}

	exchangeRate := setting.AlipayUsdToCnyRate
	if exchangeRate <= 0 {
		exchangeRate = 7.2
	}
	payMoney := dAmount.
		Mul(decimal.NewFromFloat(exchangeRate)).
		Mul(decimal.NewFromFloat(topupGroupRatio)).
		Mul(decimal.NewFromFloat(discount)).
		InexactFloat64()

	return dAmount.InexactFloat64(), payMoney, exchangeRate
}

func getAlipayMinTopup() int64 {
	minTopup := operation_setting.MinTopUp
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		minTopup = int(decimal.NewFromInt(int64(minTopup)).Mul(decimal.NewFromFloat(common.QuotaPerUnit)).IntPart())
	}
	return int64(minTopup)
}

func getAlipayCallbackURL(tradeNo string) string {
	callbackURL := strings.TrimSpace(setting.AlipayCallbackUrl)
	if callbackURL == "" {
		callbackURL = strings.TrimRight(service.GetCallbackAddress(), "/") + "/api/alipay/return"
	}

	separator := "?"
	if strings.Contains(callbackURL, "?") {
		separator = "&"
	}
	return callbackURL + separator + "trade_no=" + url.QueryEscape(tradeNo)
}

func getAlipayNotifyURL() string {
	notifyURL := strings.TrimSpace(setting.AlipayNotifyUrl)
	if notifyURL == "" {
		notifyURL = strings.TrimRight(service.GetCallbackAddress(), "/") + "/api/alipay/notify"
	}
	return notifyURL
}

func getAlipayFrontendReturnPath(suffix string) string {
	callbackURL := strings.TrimSpace(setting.AlipayCallbackUrl)
	if callbackURL == "" {
		return paymentReturnPath(suffix)
	}

	parsed, err := url.Parse(callbackURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return paymentReturnPath(suffix)
	}

	base := parsed.Scheme + "://" + parsed.Host
	return strings.TrimRight(base, "/") + common.ThemeAwarePath(suffix)
}

func validateAlipayPaidAmount(values url.Values, topUp *model.TopUp) error {
	actual := strings.TrimSpace(values.Get("total_amount"))
	if actual == "" {
		return nil
	}
	expected := strconv.FormatFloat(topUp.GetPaidAmount(), 'f', 2, 64)
	if actual != expected {
		return fmt.Errorf("alipay paid amount mismatch, expected %s, got %s", expected, actual)
	}
	return nil
}

func validateAlipayPaidAmountString(actual string, topUp *model.TopUp) error {
	actual = strings.TrimSpace(actual)
	if actual == "" {
		return fmt.Errorf("alipay paid amount is empty")
	}
	expected := strconv.FormatFloat(topUp.GetPaidAmount(), 'f', 2, 64)
	if actual != expected {
		return fmt.Errorf("alipay paid amount mismatch, expected %s, got %s", expected, actual)
	}
	return nil
}

func callAlipayBridgeCreate(c *gin.Context, req alipayBridgeCreateRequest) (string, error) {
	createURL := strings.TrimSpace(setting.AlipayBridgeCreateUrl)
	secret := strings.TrimSpace(setting.AlipayBridgeSecret)
	if createURL == "" || secret == "" {
		return "", fmt.Errorf("alipay bridge is not configured")
	}
	parsed, err := url.Parse(createURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || parsed.Path == "" {
		return "", fmt.Errorf("invalid alipay bridge create url")
	}
	body, err := common.Marshal(req)
	if err != nil {
		return "", err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := common.GetRandomString(24)
	signature, err := alipaybridge.Sign(secret, http.MethodPost, parsed.Path, timestamp, nonce, body)
	if err != nil {
		return "", err
	}
	httpReq, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, createURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set(alipaybridge.HeaderTimestamp, timestamp)
	httpReq.Header.Set(alipaybridge.HeaderNonce, nonce)
	httpReq.Header.Set(alipaybridge.HeaderSignature, signature)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("alipay bridge create failed with status %d", resp.StatusCode)
	}
	var payload alipayBridgeCreateResponse
	if err := common.DecodeJson(resp.Body, &payload); err != nil {
		return "", err
	}
	if payload.Message != "success" || payload.Data.ApproveLink == "" {
		return "", fmt.Errorf("alipay bridge returned unsuccessful response")
	}
	if payload.Data.OrderID != "" && payload.Data.OrderID != req.TradeNo {
		return "", fmt.Errorf("alipay bridge returned mismatched order id")
	}
	return payload.Data.ApproveLink, nil
}

func rememberAlipayBridgeNonce(nonce string, timestamp string) bool {
	if nonce == "" {
		return false
	}
	now := time.Now().Unix()
	alipayBridgeNonceCache.Range(func(key any, value any) bool {
		if ts, ok := value.(int64); ok && ts < now-600 {
			alipayBridgeNonceCache.Delete(key)
		}
		return true
	})
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	_, loaded := alipayBridgeNonceCache.LoadOrStore(nonce, ts)
	return !loaded
}

func RequestAlipayAmount(c *gin.Context) {
	var req AlipayPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid request"})
		return
	}
	if req.Amount < getAlipayMinTopup() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": fmt.Sprintf("topup amount cannot be less than %d", getAlipayMinTopup())})
		return
	}
	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "failed to get user group"})
		return
	}
	creditAmountUSD, payMoney, exchangeRate := getAlipayQuote(req.Amount, group)
	if payMoney <= 0.01 {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid payment amount"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": paymentAmountQuote{
			Amount:          strconv.FormatFloat(payMoney, 'f', 2, 64),
			Currency:        "CNY",
			ExchangeRate:    exchangeRate,
			CreditAmountUSD: creditAmountUSD,
		},
	})
}

func RequestAlipayPay(c *gin.Context) {
	var req AlipayPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid request"})
		return
	}
	if req.PaymentMethod != model.PaymentMethodAlipay {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid payment method"})
		return
	}
	if !isAlipayTopUpEnabled() {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "该支付方式未配置完善"})
		return
	}
	if req.Amount < getAlipayMinTopup() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": fmt.Sprintf("topup amount cannot be less than %d", getAlipayMinTopup())})
		return
	}

	id := c.GetInt("id")
	user, _ := model.GetUserById(id, false)
	group := user.Group
	creditAmountUSD, payMoney, exchangeRate := getAlipayQuote(req.Amount, group)
	if payMoney <= 0.01 {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid payment amount"})
		return
	}

	reference := fmt.Sprintf("alipay-ref-%d-%d-%s", user.Id, time.Now().UnixMilli(), randstr.String(4))
	tradeNo := "ref_" + common.Sha1([]byte(reference))
	amount := req.Amount
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		amount = decimal.NewFromInt(amount).Div(decimal.NewFromFloat(common.QuotaPerUnit)).IntPart()
	}

	topUp := &model.TopUp{
		UserId:          id,
		Amount:          amount,
		Money:           payMoney,
		CreditAmountUsd: creditAmountUSD,
		PaidAmount:      payMoney,
		PaidCurrency:    "CNY",
		ExchangeRate:    exchangeRate,
		TradeNo:         tradeNo,
		PaymentMethod:   model.PaymentMethodAlipay,
		PaymentProvider: model.PaymentProviderAlipay,
		CreateTime:      time.Now().Unix(),
		Status:          common.TopUpStatusPending,
	}
	if err := topUp.Insert(); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("Alipay create topup failed user_id=%d trade_no=%s amount=%d error=%q", id, tradeNo, req.Amount, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "failed to create order"})
		return
	}

	if setting.AlipayBridgeEnabled {
		payURL, err := callAlipayBridgeCreate(c, alipayBridgeCreateRequest{
			TradeNo:     tradeNo,
			Subject:     "API Quota Topup",
			TotalAmount: strconv.FormatFloat(payMoney, 'f', 2, 64),
			Currency:    "CNY",
			ReturnURL:   getAlipayCallbackURL(tradeNo),
			NotifyURL:   getAlipayNotifyURL(),
		})
		if err != nil {
			logger.LogError(c.Request.Context(), fmt.Sprintf("Alipay bridge create failed trade_no=%s error=%q", tradeNo, err.Error()))
			_ = model.UpdatePendingTopUpStatus(tradeNo, model.PaymentProviderAlipay, common.TopUpStatusFailed)
			c.JSON(http.StatusOK, gin.H{"message": "error", "data": "failed to start payment"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"data": gin.H{
				"approve_link": payURL,
				"order_id":     tradeNo,
			},
		})
		return
	}

	client, err := getAlipayClient()
	if err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("Alipay init client failed trade_no=%s error=%q", tradeNo, err.Error()))
		_ = model.UpdatePendingTopUpStatus(tradeNo, model.PaymentProviderAlipay, common.TopUpStatusFailed)
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "failed to initialize payment"})
		return
	}

	var p alipay.TradePagePay
	p.NotifyURL = getAlipayNotifyURL()
	p.ReturnURL = getAlipayCallbackURL(tradeNo)
	p.Subject = "API Quota Topup"
	p.OutTradeNo = tradeNo
	p.TotalAmount = strconv.FormatFloat(payMoney, 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	payURL, err := client.TradePagePay(p)
	if err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("Alipay trade page pay failed trade_no=%s error=%q", tradeNo, err.Error()))
		_ = model.UpdatePendingTopUpStatus(tradeNo, model.PaymentProviderAlipay, common.TopUpStatusFailed)
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "failed to generate payment URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"approve_link": payURL.String(),
			"order_id":     tradeNo,
		},
	})
}

func AlipayReturn(c *gin.Context) {
	tradeNo := c.Query("trade_no")
	outTradeNo := c.Query("out_trade_no")
	if tradeNo == "" && outTradeNo != "" {
		tradeNo = outTradeNo
	}
	if tradeNo == "" {
		c.Redirect(http.StatusFound, getAlipayFrontendReturnPath("/console/topup?pay=fail"))
		return
	}

	client, err := getAlipayClient()
	if err != nil {
		logger.LogError(c.Request.Context(), "Alipay client init failed: "+err.Error())
		c.Redirect(http.StatusFound, getAlipayFrontendReturnPath("/console/topup?pay=fail"))
		return
	}

	err = client.VerifySign(c.Request.Context(), c.Request.URL.Query())
	if err != nil {
		logger.LogError(c.Request.Context(), "Alipay return verify signature failed: "+err.Error())
		c.Redirect(http.StatusFound, getAlipayFrontendReturnPath("/console/topup?pay=fail"))
		return
	}

	LockOrder(tradeNo)
	defer UnlockOrder(tradeNo)

	topUp := model.GetTopUpByTradeNo(tradeNo)
	if topUp == nil || topUp.PaymentProvider != model.PaymentProviderAlipay {
		c.Redirect(http.StatusFound, getAlipayFrontendReturnPath("/console/topup?pay=fail"))
		return
	}
	if topUp.Status == common.TopUpStatusSuccess {
		c.Redirect(http.StatusFound, getAlipayFrontendReturnPath("/console/topup?pay=success"))
		return
	}

	if err := validateAlipayPaidAmount(c.Request.URL.Query(), topUp); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("Alipay return amount validation failed trade_no=%s client_ip=%s error=%q", tradeNo, c.ClientIP(), err.Error()))
		c.Redirect(http.StatusFound, getAlipayFrontendReturnPath("/console/topup?pay=fail"))
		return
	}

	if err := model.RechargeAlipay(tradeNo, c.ClientIP()); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("Alipay topup failed trade_no=%s client_ip=%s error=%q", tradeNo, c.ClientIP(), err.Error()))
		c.Redirect(http.StatusFound, getAlipayFrontendReturnPath("/console/topup?pay=fail"))
		return
	}

	c.Redirect(http.StatusFound, getAlipayFrontendReturnPath("/console/topup?pay=success"))
}

func AlipayNotify(c *gin.Context) {
	client, err := getAlipayClient()
	if err != nil {
		logger.LogError(c.Request.Context(), "Alipay client init failed: "+err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	err = c.Request.ParseForm()
	if err != nil {
		logger.LogError(c.Request.Context(), "Alipay notify parse form failed: "+err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	reqValues := c.Request.Form
	err = client.VerifySign(c.Request.Context(), reqValues)
	if err != nil {
		logger.LogError(c.Request.Context(), "Alipay notify verify signature failed: "+err.Error())
		c.String(http.StatusBadRequest, "fail")
		return
	}

	tradeNo := reqValues.Get("out_trade_no")
	tradeStatus := reqValues.Get("trade_status")

	if tradeNo == "" {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if tradeStatus != "TRADE_SUCCESS" && tradeStatus != "TRADE_FINISHED" {
		c.String(http.StatusOK, "success")
		return
	}

	LockOrder(tradeNo)
	defer UnlockOrder(tradeNo)

	topUp := model.GetTopUpByTradeNo(tradeNo)
	if topUp == nil || topUp.PaymentProvider != model.PaymentProviderAlipay {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if topUp.Status == common.TopUpStatusSuccess {
		c.String(http.StatusOK, "success")
		return
	}

	if err := validateAlipayPaidAmount(reqValues, topUp); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("Alipay notify amount validation failed trade_no=%s client_ip=%s error=%q", tradeNo, c.ClientIP(), err.Error()))
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if err := model.RechargeAlipay(tradeNo, c.ClientIP()); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("Alipay notify recharge failed trade_no=%s client_ip=%s error=%q", tradeNo, c.ClientIP(), err.Error()))
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

func AlipayBridgeSettle(c *gin.Context) {
	if !setting.AlipayBridgeEnabled {
		c.JSON(http.StatusForbidden, gin.H{"message": "error", "data": "alipay bridge is disabled"})
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "data": "invalid request body"})
		return
	}
	timestamp := c.GetHeader(alipaybridge.HeaderTimestamp)
	nonce := c.GetHeader(alipaybridge.HeaderNonce)
	signature := c.GetHeader(alipaybridge.HeaderSignature)
	if err := alipaybridge.Verify(setting.AlipayBridgeSecret, c.Request.Method, c.Request.URL.Path, timestamp, nonce, signature, body, time.Now(), 5*time.Minute); err != nil {
		logger.LogWarn(c.Request.Context(), fmt.Sprintf("Alipay bridge settle signature rejected client_ip=%s error=%q", c.ClientIP(), err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "error", "data": "invalid signature"})
		return
	}
	if !rememberAlipayBridgeNonce(nonce, timestamp) {
		logger.LogWarn(c.Request.Context(), fmt.Sprintf("Alipay bridge settle replay rejected client_ip=%s nonce=%q", c.ClientIP(), nonce))
		c.JSON(http.StatusUnauthorized, gin.H{"message": "error", "data": "replayed request"})
		return
	}

	var req alipayBridgeSettleRequest
	if err := common.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "data": "invalid request"})
		return
	}
	if req.TradeNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "data": "missing trade_no"})
		return
	}
	if req.TradeStatus != "TRADE_SUCCESS" && req.TradeStatus != "TRADE_FINISHED" {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		return
	}
	if req.Currency != "" && req.Currency != "CNY" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "data": "invalid currency"})
		return
	}

	LockOrder(req.TradeNo)
	defer UnlockOrder(req.TradeNo)

	topUp := model.GetTopUpByTradeNo(req.TradeNo)
	if topUp == nil || topUp.PaymentProvider != model.PaymentProviderAlipay {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "data": "order not found"})
		return
	}
	if err := validateAlipayPaidAmountString(req.TotalAmount, topUp); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("Alipay bridge settle amount validation failed trade_no=%s client_ip=%s error=%q", req.TradeNo, c.ClientIP(), err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "data": "amount mismatch"})
		return
	}
	if topUp.Status == common.TopUpStatusSuccess {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		return
	}
	if topUp.Status != common.TopUpStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "data": "order status invalid"})
		return
	}
	if err := model.RechargeAlipay(req.TradeNo, c.ClientIP()); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("Alipay bridge settle recharge failed trade_no=%s client_ip=%s error=%q", req.TradeNo, c.ClientIP(), err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "data": "topup failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
