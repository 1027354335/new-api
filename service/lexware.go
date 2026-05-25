package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"
)

const lexwareAPIBaseURL = "https://api.lexware.io"
const defaultLexwareTaxRate = 19

// IsLexwareConfigured 检查 Lexware API Key 是否已配置
func IsLexwareConfigured() bool {
	return strings.TrimSpace(setting.LexwareApiKey) != ""
}

// CreateLexwareInvoice 通过 Lexware Office API 创建发票
func CreateLexwareInvoice(ctx context.Context, inv *model.Invoice) (string, error) {
	if !IsLexwareConfigured() {
		return "", fmt.Errorf("Lexware API Key is not configured, please upload invoice manually")
	}
	if strings.TrimSpace(inv.Country) == "" || len(strings.TrimSpace(inv.Country)) != 2 {
		return "", fmt.Errorf("Lexware requires a two-letter ISO country code, e.g. DE")
	}

	lexwareInvoiceId, err := createLexwareInvoice(inv)
	if err != nil {
		return "", err
	}

	// Lexware limits clients to roughly 2 requests per second. Keep this path gentle.
	time.Sleep(600 * time.Millisecond)

	objectName, err := downloadLexwareInvoiceFile(ctx, lexwareInvoiceId, inv.Id)
	if err != nil {
		return "", err
	}

	common.SysLog(fmt.Sprintf("Lexware invoice created for trade_no: %s, lexware_id: %s", inv.TradeNo, lexwareInvoiceId))
	return objectName, nil
}

func createLexwareInvoice(inv *model.Invoice) (string, error) {
	paidAmount, remark := buildLexwarePaymentDetails(inv)
	unitPrice := buildLexwareGrossUnitPrice(paidAmount)
	taxConditions := map[string]string{
		"taxType": "gross",
	}

	payload := map[string]any{
		"archived":    false,
		"voucherDate": time.Now().Format("2006-01-02T15:04:05.000Z07:00"),
		"address": map[string]string{
			"name":        inv.Title,
			"street":      inv.Street,
			"city":        inv.City,
			"zip":         inv.ZipCode,
			"countryCode": strings.ToUpper(strings.TrimSpace(inv.Country)),
		},
		"lineItems": []map[string]any{
			{
				"type":               "custom",
				"name":               "Wallet top-up " + inv.TradeNo,
				"quantity":           1,
				"unitName":           "piece",
				"unitPrice":          unitPrice,
				"discountPercentage": 0,
			},
		},
		"totalPrice": map[string]string{
			"currency": "EUR",
		},
		"taxConditions": taxConditions,
		"shippingConditions": map[string]string{
			"shippingType": "none",
		},
		"title":        "Invoice",
		"introduction": "Your wallet top-up is invoiced below.",
		"remark":       remark,
	}

	body, err := common.Marshal(payload)
	if err != nil {
		return "", err
	}

	respBody, err := doLexwareRequest(http.MethodPost, "/v1/invoices?finalize=true", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	var result struct {
		Id string `json:"id"`
	}
	if err := common.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if result.Id == "" {
		return "", fmt.Errorf("Lexware invoice response missing id")
	}

	return result.Id, nil
}

func buildLexwarePaymentDetails(inv *model.Invoice) (float64, string) {
	topUp := model.GetTopUpById(inv.TopUpId)
	if topUp == nil || topUp.PaymentProvider != model.PaymentProviderPayPal {
		return inv.Money, "Thank you for your purchase."
	}

	paidAmount := topUp.GetPaidAmount()
	if paidAmount <= 0 {
		paidAmount = inv.Money
	}
	creditAmountUSD := topUp.GetCreditAmountUSD()
	paidCurrency := topUp.GetPaidCurrency("EUR")
	if paidCurrency != "EUR" {
		paidCurrency = "EUR"
	}
	exchangeRate := topUp.ExchangeRate
	if exchangeRate <= 0 && creditAmountUSD > 0 {
		exchangeRate = paidAmount / creditAmountUSD
	}

	return paidAmount, fmt.Sprintf(
		"Credited account amount: USD %.2f; paid via PayPal: %s %.2f; exchange rate: 1 USD = %.6g %s",
		roundMoney(creditAmountUSD),
		paidCurrency,
		roundMoney(paidAmount),
		exchangeRate,
		paidCurrency,
	)
}

func downloadLexwareInvoiceFile(ctx context.Context, lexwareInvoiceId string, localInvoiceId int) (string, error) {
	respBody, err := doLexwareRequest(http.MethodGet, "/v1/invoices/"+lexwareInvoiceId+"/file", "application/pdf", nil)
	if err != nil {
		return "", err
	}

	return UploadInvoicePDF(ctx, localInvoiceId, respBody)
}

func buildLexwareGrossUnitPrice(grossAmount float64) map[string]any {
	return map[string]any{
		"currency":          "EUR",
		"grossAmount":       roundUnitPrice(grossAmount),
		"taxRatePercentage": defaultLexwareTaxRate,
	}
}

func doLexwareRequest(method string, path string, accept string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, lexwareAPIBaseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(setting.LexwareApiKey))
	req.Header.Set("Accept", accept)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := GetHttpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 20*1024*1024))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("Lexware API %s %s failed with status %d: %s", method, path, resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func roundMoney(value float64) float64 {
	return float64(int(value*100+0.5)) / 100
}

func roundUnitPrice(value float64) float64 {
	return float64(int(value*10000+0.5)) / 10000
}
