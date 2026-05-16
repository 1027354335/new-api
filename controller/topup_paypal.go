package controller

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/operation_setting"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/thanhpk/randstr"
)

type PayPalPayRequest struct {
	Amount        int64  `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}

type paypalAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type paypalOrderResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Links  []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
}

type paypalCaptureResponse struct {
	Id            string `json:"id"`
	Status        string `json:"status"`
	PurchaseUnits []struct {
		ReferenceId string `json:"reference_id"`
		Payments    struct {
			Captures []struct {
				Status string `json:"status"`
				Amount struct {
					CurrencyCode string `json:"currency_code"`
					Value        string `json:"value"`
				} `json:"amount"`
			} `json:"captures"`
		} `json:"payments"`
	} `json:"purchase_units"`
}

func paypalBaseURL() string {
	if setting.PayPalMode == "live" {
		return "https://api-m.paypal.com"
	}
	return "https://api-m.sandbox.paypal.com"
}

func getPayPalAccessToken() (string, error) {
	if strings.TrimSpace(setting.PayPalClientId) == "" || strings.TrimSpace(setting.PayPalClientSecret) == "" {
		return "", fmt.Errorf("paypal credentials are not configured")
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	req, err := http.NewRequest(http.MethodPost, paypalBaseURL()+"/v1/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(setting.PayPalClientId + ":" + setting.PayPalClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := service.GetHttpClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("paypal token api status %d", resp.StatusCode)
	}

	var tokenResp paypalAccessTokenResponse
	if err := common.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}
	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("paypal token response missing access_token")
	}
	return tokenResp.AccessToken, nil
}

func genPayPalOrder(c *gin.Context, tradeNo string, payMoney float64) (string, error) {
	token, err := getPayPalAccessToken()
	if err != nil {
		return "", err
	}

	returnURL := getPayPalReturnURL(tradeNo)
	cancelURL := getPayPalFrontendReturnPath("/console/topup?pay=fail")
	value := strconv.FormatFloat(payMoney, 'f', 2, 64)
	payload := map[string]any{
		"intent": "CAPTURE",
		"purchase_units": []map[string]any{
			{
				"reference_id": tradeNo,
				"amount": map[string]string{
					"currency_code": "USD",
					"value":         value,
				},
			},
		},
		"application_context": map[string]string{
			"return_url": returnURL,
			"cancel_url": cancelURL,
		},
	}

	body, err := common.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, paypalBaseURL()+"/v2/checkout/orders", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PayPal-Request-Id", tradeNo)

	resp, err := service.GetHttpClient().Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode/100 != 2 {
		logger.LogWarn(c.Request.Context(), fmt.Sprintf("PayPal create order failed trade_no=%s status=%d body=%q", tradeNo, resp.StatusCode, string(respBody)))
		return "", fmt.Errorf("paypal create order status %d", resp.StatusCode)
	}

	var order paypalOrderResponse
	if err := common.Unmarshal(respBody, &order); err != nil {
		return "", err
	}
	for _, link := range order.Links {
		if link.Rel == "approve" && link.Href != "" {
			return link.Href, nil
		}
	}
	return "", fmt.Errorf("paypal order response missing approve link")
}

func getPayPalReturnURL(tradeNo string) string {
	callbackURL := strings.TrimSpace(setting.PayPalCallbackUrl)
	if callbackURL == "" {
		callbackURL = strings.TrimRight(service.GetCallbackAddress(), "/") + "/api/paypal/return"
	}

	separator := "?"
	if strings.Contains(callbackURL, "?") {
		separator = "&"
	}
	return callbackURL + separator + "trade_no=" + url.QueryEscape(tradeNo)
}

func getPayPalFrontendReturnPath(suffix string) string {
	callbackURL := strings.TrimSpace(setting.PayPalCallbackUrl)
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

func capturePayPalOrder(orderId string, expectedTradeNo string, expectedAmount float64) error {
	token, err := getPayPalAccessToken()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, paypalBaseURL()+"/v2/checkout/orders/"+url.PathEscape(orderId)+"/capture", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("PayPal-Request-Id", "capture-"+orderId)

	resp, err := service.GetHttpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("paypal capture status %d body=%s", resp.StatusCode, string(body))
	}

	var order paypalCaptureResponse
	if err := common.Unmarshal(body, &order); err != nil {
		return err
	}
	if order.Status != "COMPLETED" {
		return fmt.Errorf("paypal capture status %s", order.Status)
	}

	expectedValue := strconv.FormatFloat(expectedAmount, 'f', 2, 64)
	for _, unit := range order.PurchaseUnits {
		if unit.ReferenceId != expectedTradeNo {
			continue
		}
		for _, capture := range unit.Payments.Captures {
			if capture.Status == "COMPLETED" &&
				capture.Amount.CurrencyCode == "USD" &&
				capture.Amount.Value == expectedValue {
				return nil
			}
		}
	}

	return fmt.Errorf("paypal capture does not match local order")
}

func getPayPalPayMoney(amount int64, group string) float64 {
	dAmount := decimal.NewFromInt(amount)
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		dAmount = dAmount.Div(decimal.NewFromFloat(common.QuotaPerUnit))
	}

	topupGroupRatio := common.GetTopupGroupRatio(group)
	if topupGroupRatio == 0 {
		topupGroupRatio = 1
	}

	discount := 1.0
	if ds, ok := operation_setting.GetPaymentSetting().AmountDiscount[int(amount)]; ok && ds > 0 {
		discount = ds
	}

	return dAmount.Mul(decimal.NewFromFloat(topupGroupRatio)).Mul(decimal.NewFromFloat(discount)).InexactFloat64()
}

func getPayPalMinTopup() int64 {
	minTopup := operation_setting.MinTopUp
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		minTopup = int(decimal.NewFromInt(int64(minTopup)).Mul(decimal.NewFromFloat(common.QuotaPerUnit)).IntPart())
	}
	return int64(minTopup)
}

func RequestPayPalAmount(c *gin.Context) {
	var req PayPalPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid request"})
		return
	}
	if req.Amount < getPayPalMinTopup() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": fmt.Sprintf("topup amount cannot be less than %d", getPayPalMinTopup())})
		return
	}
	id := c.GetInt("id")
	group, err := model.GetUserGroup(id, true)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "failed to get user group"})
		return
	}
	payMoney := getPayPalPayMoney(req.Amount, group)
	if payMoney <= 0.01 {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid payment amount"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": strconv.FormatFloat(payMoney, 'f', 2, 64)})
}

func RequestPayPalPay(c *gin.Context) {
	var req PayPalPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid request"})
		return
	}
	if req.PaymentMethod != model.PaymentMethodPayPal {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid payment method"})
		return
	}
	if !isPayPalTopUpEnabled() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "paypal is not configured"})
		return
	}
	if req.Amount < getPayPalMinTopup() {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": fmt.Sprintf("topup amount cannot be less than %d", getPayPalMinTopup())})
		return
	}

	id := c.GetInt("id")
	user, _ := model.GetUserById(id, false)
	group := user.Group
	payMoney := getPayPalPayMoney(req.Amount, group)
	if payMoney <= 0.01 {
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "invalid payment amount"})
		return
	}

	reference := fmt.Sprintf("paypal-ref-%d-%d-%s", user.Id, time.Now().UnixMilli(), randstr.String(4))
	tradeNo := "ref_" + common.Sha1([]byte(reference))
	amount := req.Amount
	if operation_setting.GetQuotaDisplayType() == operation_setting.QuotaDisplayTypeTokens {
		amount = decimal.NewFromInt(amount).Div(decimal.NewFromFloat(common.QuotaPerUnit)).IntPart()
	}

	topUp := &model.TopUp{
		UserId:          id,
		Amount:          amount,
		Money:           payMoney,
		TradeNo:         tradeNo,
		PaymentMethod:   model.PaymentMethodPayPal,
		PaymentProvider: model.PaymentProviderPayPal,
		CreateTime:      time.Now().Unix(),
		Status:          common.TopUpStatusPending,
	}
	if err := topUp.Insert(); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("PayPal create topup failed user_id=%d trade_no=%s amount=%d error=%q", id, tradeNo, req.Amount, err.Error()))
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "failed to create order"})
		return
	}

	approveLink, err := genPayPalOrder(c, tradeNo, payMoney)
	if err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("PayPal create order failed user_id=%d trade_no=%s amount=%d error=%q", id, tradeNo, req.Amount, err.Error()))
		_ = model.UpdatePendingTopUpStatus(tradeNo, model.PaymentProviderPayPal, common.TopUpStatusFailed)
		c.JSON(http.StatusOK, gin.H{"message": "error", "data": "failed to create payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"approve_link": approveLink,
			"order_id":     tradeNo,
		},
	})
}

func PayPalReturn(c *gin.Context) {
	tradeNo := c.Query("trade_no")
	orderId := c.Query("token")
	if tradeNo == "" || orderId == "" {
		c.Redirect(http.StatusFound, getPayPalFrontendReturnPath("/console/topup?pay=fail"))
		return
	}

	LockOrder(tradeNo)
	defer UnlockOrder(tradeNo)

	topUp := model.GetTopUpByTradeNo(tradeNo)
	if topUp == nil || topUp.PaymentProvider != model.PaymentProviderPayPal {
		c.Redirect(http.StatusFound, getPayPalFrontendReturnPath("/console/topup?pay=fail"))
		return
	}
	if topUp.Status == common.TopUpStatusSuccess {
		c.Redirect(http.StatusFound, getPayPalFrontendReturnPath("/console/topup?pay=success"))
		return
	}

	if err := capturePayPalOrder(orderId, tradeNo, topUp.Money); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("PayPal capture failed trade_no=%s paypal_order_id=%s client_ip=%s error=%q", tradeNo, orderId, c.ClientIP(), err.Error()))
		c.Redirect(http.StatusFound, getPayPalFrontendReturnPath("/console/topup?pay=fail"))
		return
	}

	if err := model.RechargePayPal(tradeNo, c.ClientIP()); err != nil {
		logger.LogError(c.Request.Context(), fmt.Sprintf("PayPal topup failed trade_no=%s paypal_order_id=%s client_ip=%s error=%q", tradeNo, orderId, c.ClientIP(), err.Error()))
		c.Redirect(http.StatusFound, getPayPalFrontendReturnPath("/console/topup?pay=fail"))
		return
	}

	c.Redirect(http.StatusFound, getPayPalFrontendReturnPath("/console/topup?pay=success"))
}
