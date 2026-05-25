package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/QuantumNous/new-api/pkg/alipaybridge"
	"github.com/smartwalle/alipay/v3"
)

type config struct {
	ListenAddr        string
	AppID             string
	PrivateKey        string
	PublicKey         string
	Sandbox           bool
	PublicBaseURL     string
	OverseasSettleURL string
	SharedSecret      string
	ReturnSuccessURL  string
	ReturnFailURL     string
}

type createRequest struct {
	TradeNo     string `json:"trade_no"`
	Subject     string `json:"subject"`
	TotalAmount string `json:"total_amount"`
	Currency    string `json:"currency"`
	ReturnURL   string `json:"return_url"`
	NotifyURL   string `json:"notify_url"`
}

type createResponse struct {
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
}

type settleRequest struct {
	TradeNo       string `json:"trade_no"`
	AlipayTradeNo string `json:"alipay_trade_no,omitempty"`
	TradeStatus   string `json:"trade_status"`
	TotalAmount   string `json:"total_amount"`
	Currency      string `json:"currency,omitempty"`
}

var nonceCache sync.Map

func main() {
	loadDotEnv()
	cfg := loadConfig()
	client, err := newAlipayClient(cfg)
	if err != nil {
		log.Fatalf("init alipay client failed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/alipay/create", func(w http.ResponseWriter, r *http.Request) {
		handleCreate(w, r, cfg, client)
	})
	mux.HandleFunc("/api/alipay/notify", func(w http.ResponseWriter, r *http.Request) {
		handleNotify(w, r, cfg, client)
	})
	mux.HandleFunc("/api/alipay/return", func(w http.ResponseWriter, r *http.Request) {
		handleReturn(w, r, cfg, client)
	})
	mux.HandleFunc("/health", handleHealth)

	log.Printf("alipay bridge listening on %s", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, _ = w.Write([]byte("ok"))
}

func loadDotEnv() {
	paths := make([]string, 0, 3)
	if custom := strings.TrimSpace(os.Getenv("ALIPAY_BRIDGE_ENV_FILE")); custom != "" {
		paths = append(paths, custom)
	} else {
		paths = append(paths, ".env")
		if exe, err := os.Executable(); err == nil {
			paths = append(paths, filepath.Join(filepath.Dir(exe), ".env"))
		}
	}
	seen := map[string]bool{}
	for _, p := range paths {
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		if err := loadEnvFile(p); err != nil && !os.IsNotExist(err) {
			log.Printf("load env file failed path=%s error=%v", p, err)
		}
	}
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key, value, ok := parseEnvLine(scanner.Text())
		if !ok {
			continue
		}
		if os.Getenv(key) != "" {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func parseEnvLine(line string) (string, string, bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	line = strings.TrimPrefix(line, "export ")
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	if key == "" || strings.ContainsAny(key, " \t") {
		return "", "", false
	}
	if len(value) >= 2 {
		quote := value[0]
		if (quote == '"' || quote == '\'') && value[len(value)-1] == quote {
			value = value[1 : len(value)-1]
		}
	} else if idx := strings.Index(value, " #"); idx >= 0 {
		value = strings.TrimSpace(value[:idx])
	}
	return key, value, true
}

func loadConfig() config {
	cfg := config{
		ListenAddr:        getenv("ALIPAY_BRIDGE_LISTEN_ADDR", ":8088"),
		AppID:             strings.TrimSpace(os.Getenv("ALIPAY_APP_ID")),
		PrivateKey:        strings.TrimSpace(os.Getenv("ALIPAY_PRIVATE_KEY")),
		PublicKey:         strings.TrimSpace(os.Getenv("ALIPAY_PUBLIC_KEY")),
		Sandbox:           getenvBool("ALIPAY_SANDBOX", false),
		PublicBaseURL:     strings.TrimRight(strings.TrimSpace(os.Getenv("ALIPAY_BRIDGE_PUBLIC_BASE_URL")), "/"),
		OverseasSettleURL: strings.TrimSpace(os.Getenv("ALIPAY_BRIDGE_OVERSEAS_SETTLE_URL")),
		SharedSecret:      strings.TrimSpace(os.Getenv("ALIPAY_BRIDGE_SECRET")),
		ReturnSuccessURL:  strings.TrimSpace(os.Getenv("ALIPAY_BRIDGE_RETURN_SUCCESS_URL")),
		ReturnFailURL:     strings.TrimSpace(os.Getenv("ALIPAY_BRIDGE_RETURN_FAIL_URL")),
	}
	missing := make([]string, 0)
	required := map[string]string{
		"ALIPAY_APP_ID":                     cfg.AppID,
		"ALIPAY_PRIVATE_KEY":                cfg.PrivateKey,
		"ALIPAY_PUBLIC_KEY":                 cfg.PublicKey,
		"ALIPAY_BRIDGE_OVERSEAS_SETTLE_URL": cfg.OverseasSettleURL,
		"ALIPAY_BRIDGE_SECRET":              cfg.SharedSecret,
	}
	for key, value := range required {
		if value == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		log.Fatalf("missing required env: %s", strings.Join(missing, ", "))
	}
	return cfg
}

func newAlipayClient(cfg config) (*alipay.Client, error) {
	client, err := alipay.New(cfg.AppID, cfg.PrivateKey, !cfg.Sandbox)
	if err != nil {
		return nil, err
	}
	if err := client.LoadAliPayPublicKey(cfg.PublicKey); err != nil {
		return nil, err
	}
	return client, nil
}

func handleCreate(w http.ResponseWriter, r *http.Request, cfg config, client *alipay.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, ok := readAndVerifyBridgeRequest(w, r, cfg.SharedSecret)
	if !ok {
		return
	}
	var req createRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, createResponse{Message: "error", Data: map[string]any{"error": "invalid request"}})
		return
	}
	if err := validateCreateRequest(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, createResponse{Message: "error", Data: map[string]any{"error": err.Error()}})
		return
	}
	if req.NotifyURL == "" && cfg.PublicBaseURL != "" {
		req.NotifyURL = cfg.PublicBaseURL + "/api/alipay/notify"
	}
	if req.ReturnURL == "" && cfg.PublicBaseURL != "" {
		req.ReturnURL = cfg.PublicBaseURL + "/api/alipay/return"
	}

	var p alipay.TradePagePay
	p.NotifyURL = req.NotifyURL
	p.ReturnURL = appendTradeNo(req.ReturnURL, req.TradeNo)
	p.Subject = req.Subject
	p.OutTradeNo = req.TradeNo
	p.TotalAmount = req.TotalAmount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	payURL, err := client.TradePagePay(p)
	if err != nil {
		log.Printf("alipay create failed trade_no=%s error=%v", req.TradeNo, err)
		writeJSON(w, http.StatusBadGateway, createResponse{Message: "error", Data: map[string]any{"error": "failed to create alipay order"}})
		return
	}
	writeJSON(w, http.StatusOK, createResponse{
		Message: "success",
		Data: map[string]any{
			"approve_link": payURL.String(),
			"order_id":     req.TradeNo,
		},
	})
}

func handleNotify(w http.ResponseWriter, r *http.Request, cfg config, client *alipay.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("alipay notify parse failed: %v", err)
		_, _ = w.Write([]byte("fail"))
		return
	}
	if err := client.VerifySign(r.Context(), r.Form); err != nil {
		log.Printf("alipay notify verify failed: %v", err)
		_, _ = w.Write([]byte("fail"))
		return
	}
	status := r.Form.Get("trade_status")
	if status != "TRADE_SUCCESS" && status != "TRADE_FINISHED" {
		_, _ = w.Write([]byte("success"))
		return
	}
	if err := notifyOverseas(r.Context(), cfg, settleFromValues(r.Form, status)); err != nil {
		log.Printf("overseas settle failed trade_no=%s error=%v", r.Form.Get("out_trade_no"), err)
		_, _ = w.Write([]byte("fail"))
		return
	}
	_, _ = w.Write([]byte("success"))
}

func handleReturn(w http.ResponseWriter, r *http.Request, cfg config, client *alipay.Client) {
	values := r.URL.Query()
	redirectURL := cfg.ReturnFailURL
	if err := client.VerifySign(r.Context(), values); err != nil {
		log.Printf("alipay return verify failed: %v", err)
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}
	status := values.Get("trade_status")
	if status == "" {
		status = "TRADE_SUCCESS"
	}
	if status == "TRADE_SUCCESS" || status == "TRADE_FINISHED" {
		if err := notifyOverseas(r.Context(), cfg, settleFromValues(values, status)); err != nil {
			log.Printf("overseas return settle failed trade_no=%s error=%v", values.Get("out_trade_no"), err)
		} else {
			redirectURL = cfg.ReturnSuccessURL
		}
	}
	if redirectURL == "" {
		redirectURL = "/"
	}
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func notifyOverseas(ctx context.Context, cfg config, payload settleRequest) error {
	if payload.TradeNo == "" {
		return fmt.Errorf("missing trade_no")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	parsed, err := url.Parse(cfg.OverseasSettleURL)
	if err != nil || parsed.Path == "" {
		return fmt.Errorf("invalid overseas settle url")
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := randomString(24)
	signature, err := alipaybridge.Sign(cfg.SharedSecret, http.MethodPost, parsed.Path, timestamp, nonce, body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.OverseasSettleURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(alipaybridge.HeaderTimestamp, timestamp)
	req.Header.Set(alipaybridge.HeaderNonce, nonce)
	req.Header.Set(alipaybridge.HeaderSignature, signature)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("settle returned status %d", resp.StatusCode)
	}
	return nil
}

func readAndVerifyBridgeRequest(w http.ResponseWriter, r *http.Request, secret string) ([]byte, bool) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return nil, false
	}
	timestamp := r.Header.Get(alipaybridge.HeaderTimestamp)
	nonce := r.Header.Get(alipaybridge.HeaderNonce)
	signature := r.Header.Get(alipaybridge.HeaderSignature)
	if err := alipaybridge.Verify(secret, r.Method, r.URL.Path, timestamp, nonce, signature, body, time.Now(), 5*time.Minute); err != nil {
		log.Printf("bridge create signature rejected: %v", err)
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return nil, false
	}
	if !rememberNonce(nonce, timestamp) {
		http.Error(w, "replayed request", http.StatusUnauthorized)
		return nil, false
	}
	return body, true
}

func rememberNonce(nonce string, timestamp string) bool {
	now := time.Now().Unix()
	nonceCache.Range(func(key any, value any) bool {
		if ts, ok := value.(int64); ok && ts < now-600 {
			nonceCache.Delete(key)
		}
		return true
	})
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if nonce == "" || err != nil {
		return false
	}
	_, loaded := nonceCache.LoadOrStore(nonce, ts)
	return !loaded
}

func validateCreateRequest(req *createRequest) error {
	if req.TradeNo == "" {
		return fmt.Errorf("missing trade_no")
	}
	if req.Subject == "" {
		return fmt.Errorf("missing subject")
	}
	if req.Currency != "" && req.Currency != "CNY" {
		return fmt.Errorf("invalid currency")
	}
	amount, err := strconv.ParseFloat(req.TotalAmount, 64)
	if err != nil || amount <= 0 {
		return fmt.Errorf("invalid total_amount")
	}
	return nil
}

func settleFromValues(values url.Values, status string) settleRequest {
	return settleRequest{
		TradeNo:       values.Get("out_trade_no"),
		AlipayTradeNo: values.Get("trade_no"),
		TradeStatus:   status,
		TotalAmount:   values.Get("total_amount"),
		Currency:      "CNY",
	}
}

func appendTradeNo(rawURL string, tradeNo string) string {
	if rawURL == "" {
		return rawURL
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	query := parsed.Query()
	if query.Get("trade_no") == "" {
		query.Set("trade_no", tradeNo)
	}
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "json error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func randomString(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(buf)
}

func getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getenvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
