package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting/storage_setting"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/signintech/gopdf"
)

func init() {
	model.PostRechargeHook = GenerateAndUploadAgreementPDF
}

// AgreementTemplate defines the structure for agreements in different languages
type AgreementTemplate struct {
	Title    string
	MetaKeys map[string]string
	Sections []string
}

// agreementTemplates is defined in e:/allCode/new-oss-ai/service/agreement_templates.go

// helper check if file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Find a suitable CJK TrueType (.ttf) font on the current operating system.
// Note: gopdf does not support TTC (TrueType Collection) or OTF formats reliably.
func findFontPath() string {
	paths := []string{
		// Windows CJK TTF
		"C:\\Windows\\Fonts\\simkai.ttf",  // KaiTi (楷体)
		"C:\\Windows\\Fonts\\simfang.ttf", // FangSong (仿宋)
		// Linux CJK TTF (commonly installed via fonts-droid-fallback)
		"/usr/share/fonts/truetype/droid/DroidSansFallback.ttf",
		// macOS CJK TTF
		"/Library/Fonts/Arial Unicode.ttf",

		// Fallback Standard English TTF (if CJK is not found)
		"C:\\Windows\\Fonts\\arial.ttf",
		"C:\\Windows\\Fonts\\times.ttf",
		"C:\\Windows\\Fonts\\tahoma.ttf",
		"C:\\Windows\\Fonts\\calibri.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/truetype/freefont/FreeSans.ttf",
		"/Library/Fonts/Arial.ttf",
		"/System/Library/Fonts/Keyboard.ttf",
	}
	for _, p := range paths {
		if fileExists(p) {
			return p
		}
	}
	return ""
}

func resolveAgreementLanguage(topUp *model.TopUp) string {
	lang := ""
	if topUp != nil {
		lang = topUp.AgreementLanguage
	}
	if strings.TrimSpace(lang) == "" && topUp != nil {
		lang = model.GetUserLanguage(topUp.UserId)
	}

	normalized := strings.ToLower(strings.TrimSpace(lang))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	if strings.HasPrefix(normalized, "zh") {
		return "zh"
	}
	if _, ok := agreementTemplates[normalized]; ok {
		return normalized
	}
	return "en"
}

// findAnyTTF searches OS font directories recursively for the first available .ttf file.
func findAnyTTF() string {
	dirs := []string{
		"C:\\Windows\\Fonts",
		"/usr/share/fonts",
		"/Library/Fonts",
		"/System/Library/Fonts",
	}
	for _, dir := range dirs {
		if !fileExists(dir) {
			continue
		}
		var foundPath string
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".ttf" {
				foundPath = path
				return filepath.SkipDir // Stop walking
			}
			return nil
		})
		if err == nil && foundPath != "" {
			return foundPath
		}
	}
	return ""
}

// GenerateAndUploadAgreementPDF creates a PDF agreement asynchronously, uploads it to MinIO and updates TopUp record.
func GenerateAndUploadAgreementPDF(topUpId int) {
	// Execute in a background goroutine context
	ctx := context.Background()

	// 1. Fetch TopUp order details
	topUp := model.GetTopUpById(topUpId)
	if topUp == nil {
		common.SysLog(fmt.Sprintf("[Agreement PDF] Topup ID %d not found", topUpId))
		return
	}

	// 2. Fetch User info
	user, err := model.GetUserById(topUp.UserId, false)
	if err != nil {
		common.SysLog(fmt.Sprintf("[Agreement PDF] Failed to fetch user info for topup ID %d: %v", topUpId, err))
		return
	}

	// 3. Determine the agreement language captured at checkout.
	lang := resolveAgreementLanguage(topUp)
	tmpl := agreementTemplates[lang]

	// 4. Set up metadata values
	formattedDate := ""
	if topUp.CompleteTime > 0 {
		t := time.Unix(topUp.CompleteTime, 0)
		if lang == "zh" || lang == "ja" {
			formattedDate = t.Format("2006年01月02日")
		} else {
			formattedDate = t.Format("2006-01-02")
		}
	} else {
		t := time.Now()
		if lang == "zh" || lang == "ja" {
			formattedDate = t.Format("2006年01月02日")
		} else {
			formattedDate = t.Format("2006-01-02")
		}
	}

	userText := fmt.Sprintf("%s (%s)", user.Username, user.Email)
	if user.Email == "" {
		userText = user.Username
	}

	amountText := ""
	if topUp.PaidAmount > 0 {
		amountText = fmt.Sprintf("%.2f %s", topUp.PaidAmount, topUp.PaidCurrency)
	} else {
		amountText = fmt.Sprintf("%.2f USD", topUp.Money)
	}

	// 5. Initialize GoPDF document
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	fontPath := findFontPath()
	if fontPath == "" {
		common.SysLog("[Agreement PDF] Standard CJK TTF font not found, scanning OS for any TTF fallback...")
		fontPath = findAnyTTF()
	}

	fontName := "sans-serif"
	fontLoaded := false

	if fontPath != "" {
		common.SysLog(fmt.Sprintf("[Agreement PDF] Loading TTF font from: %s", fontPath))
		err = pdf.AddTTFFont(fontName, fontPath)
		if err == nil {
			fontLoaded = true
		} else {
			common.SysLog(fmt.Sprintf("[Agreement PDF] Failed to load TTF font from %s: %v", fontPath, err))
		}
	}

	if !fontLoaded {
		// Abort PDF generation to prevent nil-pointer dereference panic in gopdf when Cell is called without active subset font
		common.SysLog("[Agreement PDF] CRITICAL ERROR: No valid TrueType (.ttf) font could be loaded on this system. Aborting agreement PDF generation to prevent server crash.")
		return
	}

	// Set initial default font and size
	_ = pdf.SetFont(fontName, "", 12)

	// Setup watermark draw helper
	drawWatermark := func() {
		pdf.SetTextColor(230, 230, 230)
		_ = pdf.SetFont(fontName, "", 10)
		watermarkText := fmt.Sprintf("%s <%s>", user.Username, user.Email)
		if user.Email == "" {
			watermarkText = user.Username
		}

		// Tile watermark across the background page grid
		for y := 40.0; y < 800.0; y += 140.0 {
			for x := 30.0; x < 550.0; x += 190.0 {
				pdf.SetXY(x, y)
				_ = pdf.Cell(nil, watermarkText)
			}
		}
		// Reset styles back to normal text
		pdf.SetTextColor(51, 51, 51)
	}

	// Margin Setup
	leftMargin := 40.0
	rightMargin := 40.0
	topMargin := 50.0
	bottomMargin := 50.0
	pageWidth := 595.28
	pageHeight := 841.89
	contentWidth := pageWidth - leftMargin - rightMargin
	currentY := topMargin

	// Check page overflow and auto append new page
	ensureSpace := func(neededHeight float64) {
		if currentY+neededHeight > pageHeight-bottomMargin {
			pdf.AddPage()
			drawWatermark()
			currentY = topMargin
		}
	}

	// Add First Page
	pdf.AddPage()
	drawWatermark()

	// 6. Draw Title
	_ = pdf.SetFont(fontName, "", 18)
	ensureSpace(30)
	pdf.SetXY(leftMargin, currentY)
	_ = pdf.Cell(nil, tmpl.Title)
	currentY += 40

	// 7. Draw Meta Details
	_ = pdf.SetFont(fontName, "", 9)
	pdf.SetTextColor(100, 100, 100)

	metaKeysOrder := []string{"version", "effective", "provider", "address", "email", "user", "order", "signMethod", "applicable", "declaration"}
	for _, key := range metaKeysOrder {
		rawVal, exists := tmpl.MetaKeys[key]
		if !exists {
			continue
		}
		val := rawVal
		val = strings.ReplaceAll(val, "{{date}}", formattedDate)
		val = strings.ReplaceAll(val, "{{user}}", userText)
		val = strings.ReplaceAll(val, "{{order}}", topUp.TradeNo)
		val = strings.ReplaceAll(val, "{{amount}}", amountText)

		lines, err := pdf.SplitText(val, contentWidth)
		if err != nil {
			lines = []string{val}
		}

		for _, line := range lines {
			ensureSpace(15)
			pdf.SetXY(leftMargin, currentY)
			_ = pdf.Cell(nil, line)
			currentY += 15
		}
	}
	currentY += 15

	// 8. Draw Body Sections
	_ = pdf.SetFont(fontName, "", 10.5)
	pdf.SetTextColor(51, 51, 51)

	for _, section := range tmpl.Sections {
		lines, err := pdf.SplitText(section, contentWidth)
		if err != nil {
			lines = []string{section}
		}

		// Treat section block together
		ensureSpace(float64(len(lines)*16 + 10))

		for _, line := range lines {
			pdf.SetXY(leftMargin, currentY)
			_ = pdf.Cell(nil, line)
			currentY += 16
		}
		currentY += 8 // Spacing between articles
	}

	// 9. Export PDF buffer bytes
	var buf bytes.Buffer
	_, err = pdf.WriteTo(&buf)
	if err != nil {
		common.SysLog(fmt.Sprintf("[Agreement PDF] Failed to write PDF bytes for topup %d: %v", topUpId, err))
		return
	}

	// 10. Upload to MinIO
	cfg := storage_setting.GetStorageSetting()
	if !cfg.Enabled {
		common.SysLog(fmt.Sprintf("[Agreement PDF] Storage configuration is disabled, cannot upload agreement for topup %d", topUpId))
		return
	}

	client, err := storage_setting.GetClient()
	if err != nil {
		common.SysLog(fmt.Sprintf("[Agreement PDF] Failed to get storage client: %v", err))
		return
	}

	// Ensure bucket exists in target region
	err = ensureBucketExists(ctx, client, cfg.Bucket, cfg.Region)
	if err != nil {
		common.SysLog(fmt.Sprintf("[Agreement PDF] Failed to ensure bucket exists: %v", err))
		return
	}

	// Agreements folder path: agreements/{topUpId}/{uuid}.pdf
	fileName := fmt.Sprintf("%s.pdf", uuid.New().String())
	objectName := fmt.Sprintf("agreements/%d/%s", topUpId, fileName)

	pdfData := buf.Bytes()
	reader := bytes.NewReader(pdfData)
	_, err = client.PutObject(ctx, cfg.Bucket, objectName, reader, int64(len(pdfData)), minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		common.SysLog(fmt.Sprintf("[Agreement PDF] Failed to upload agreement PDF to MinIO: %v", err))
		return
	}

	// 11. Associate PDF with TopUp order in DB
	topUp.AgreementPdf = "minio:" + objectName
	err = topUp.Update()
	if err != nil {
		common.SysLog(fmt.Sprintf("[Agreement PDF] Failed to update TopUp record with agreement path: %v", err))
		return
	}

	common.SysLog(fmt.Sprintf("[Agreement PDF] Successfully generated and uploaded agreement PDF for topup %d to minio:%s", topUpId, objectName))
}

// DownloadAgreementReader fetches the agreement PDF from MinIO for downloading.
func DownloadAgreementReader(ctx context.Context, objectName string) (io.ReadCloser, int64, string, error) {
	cfg := storage_setting.GetStorageSetting()
	if !cfg.Enabled {
		return nil, 0, "", fmt.Errorf("storage setting is disabled")
	}

	client, err := storage_setting.GetClient()
	if err != nil {
		return nil, 0, "", err
	}

	info, err := client.StatObject(ctx, cfg.Bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to stat agreement: %w", err)
	}

	obj, err := client.GetObject(ctx, cfg.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to get agreement object: %w", err)
	}

	return obj, info.Size, "application/pdf", nil
}
