package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/service"

	"github.com/gin-gonic/gin"
)

type ApplyInvoiceRequest struct {
	TopUpId       int    `json:"topup_id"`
	BillingType   string `json:"billing_type"`
	Title         string `json:"title"`
	TaxId         string `json:"tax_id"`
	Email         string `json:"email"`
	Street        string `json:"street"`
	AddressDetail string `json:"address_detail"`
	City          string `json:"city"`
	ZipCode       string `json:"zip_code"`
	Country       string `json:"country"`
}

type InvoiceTitleRequest struct {
	Name          string `json:"name"`
	BillingType   string `json:"billing_type"`
	Title         string `json:"title"`
	TaxId         string `json:"tax_id"`
	Email         string `json:"email"`
	Street        string `json:"street"`
	AddressDetail string `json:"address_detail"`
	City          string `json:"city"`
	ZipCode       string `json:"zip_code"`
	Country       string `json:"country"`
	IsDefault     bool   `json:"is_default"`
}

// ApplyInvoice 用户申请开票
func ApplyInvoice(c *gin.Context) {
	var req ApplyInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "参数错误")
		return
	}

	userId := c.GetInt("id")
	username := c.GetString("username")

	// 查询充值订单
	topUp := model.GetTopUpById(req.TopUpId)
	if topUp == nil {
		common.ApiErrorMsg(c, "充值订单不存在")
		return
	}

	// 校验订单归属
	if topUp.UserId != userId {
		common.ApiErrorMsg(c, "无权操作此订单")
		return
	}

	// 校验订单状态
	if topUp.Status != common.TopUpStatusSuccess {
		common.ApiErrorMsg(c, "只有已完成的充值订单才能申请开票")
		return
	}

	// 检查是否已申请过发票
	existingStatus, err := model.GetInvoiceStatusByTopUpId(req.TopUpId)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if existingStatus != "" {
		common.ApiErrorMsg(c, "该订单已申请过发票，状态："+existingStatus)
		return
	}

	// 校验必填字段
	req.BillingType = strings.TrimSpace(req.BillingType)
	req.Title = strings.TrimSpace(req.Title)
	req.TaxId = strings.TrimSpace(req.TaxId)
	req.Email = strings.TrimSpace(req.Email)
	req.Street = strings.TrimSpace(req.Street)
	req.AddressDetail = strings.TrimSpace(req.AddressDetail)
	req.City = strings.TrimSpace(req.City)
	req.ZipCode = strings.TrimSpace(req.ZipCode)
	req.Country = strings.ToUpper(strings.TrimSpace(req.Country))

	if req.BillingType != "personal" && req.BillingType != "enterprise" {
		common.ApiErrorMsg(c, "billing_type 必须为 personal 或 enterprise")
		return
	}
	if req.Title == "" {
		common.ApiErrorMsg(c, "发票抬头不能为空")
		return
	}
	if req.Email == "" {
		common.ApiErrorMsg(c, "邮箱不能为空")
		return
	}

	// 企业开票需要税号
	if req.BillingType == "enterprise" && req.TaxId == "" {
		common.ApiErrorMsg(c, "企业开票必须提供税号")
		return
	}

	// PayPal 支付方式需要地址信息
	if topUp.PaymentMethod == "paypal" {
		if req.Street == "" || req.City == "" || req.ZipCode == "" || req.Country == "" {
			common.ApiErrorMsg(c, "PayPal 支付订单需要提供完整地址信息")
			return
		}
		if !isInvoiceCountryCode(req.Country) {
			common.ApiErrorMsg(c, "country 必须为两位 ISO 国家码，例如 DE、CN、US")
			return
		}
	}

	invoice := &model.Invoice{
		UserId:        userId,
		Username:      username,
		TopUpId:       req.TopUpId,
		TradeNo:       topUp.TradeNo,
		Money:         topUp.GetPaidAmount(),
		PaymentMethod: topUp.PaymentMethod,
		BillingType:   req.BillingType,
		Title:         req.Title,
		TaxId:         req.TaxId,
		Email:         req.Email,
		Street:        req.Street,
		AddressDetail: req.AddressDetail,
		City:          req.City,
		ZipCode:       req.ZipCode,
		Country:       req.Country,
		Status:        model.InvoiceStatusPending,
		CreateTime:    time.Now().Unix(),
	}

	if err := invoice.Insert(); err != nil {
		common.ApiError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": ""})
}

func isInvoiceCountryCode(country string) bool {
	if len(country) != 2 {
		return false
	}
	for _, char := range country {
		if char < 'A' || char > 'Z' {
			return false
		}
	}
	return true
}

func normalizeInvoiceTitleRequest(req *InvoiceTitleRequest) {
	req.Name = strings.TrimSpace(req.Name)
	req.BillingType = strings.TrimSpace(req.BillingType)
	req.Title = strings.TrimSpace(req.Title)
	req.TaxId = strings.TrimSpace(req.TaxId)
	req.Email = strings.TrimSpace(req.Email)
	req.Street = strings.TrimSpace(req.Street)
	req.AddressDetail = strings.TrimSpace(req.AddressDetail)
	req.City = strings.TrimSpace(req.City)
	req.ZipCode = strings.TrimSpace(req.ZipCode)
	req.Country = strings.ToUpper(strings.TrimSpace(req.Country))
}

func validateInvoiceTitleRequest(req InvoiceTitleRequest) string {
	if req.Name == "" {
		return "卡片名称不能为空"
	}
	if req.BillingType != "personal" && req.BillingType != "enterprise" {
		return "billing_type 必须为 personal 或 enterprise"
	}
	if req.Title == "" {
		return "发票抬头不能为空"
	}
	if req.Email == "" {
		return "邮箱不能为空"
	}
	if req.BillingType == "enterprise" && req.TaxId == "" {
		return "企业开票必须提供税号"
	}
	if req.Country != "" && !isInvoiceCountryCode(req.Country) {
		return "country 必须为两位 ISO 国家码，例如 DE、CN、US"
	}
	return ""
}

// ListInvoiceTitles 获取当前用户的发票抬头卡片
func ListInvoiceTitles(c *gin.Context) {
	userId := c.GetInt("id")
	titles, err := model.GetUserInvoiceTitles(userId)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, titles)
}

// CreateInvoiceTitle 创建发票抬头卡片
func CreateInvoiceTitle(c *gin.Context) {
	var req InvoiceTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "参数错误")
		return
	}
	normalizeInvoiceTitleRequest(&req)
	if msg := validateInvoiceTitleRequest(req); msg != "" {
		common.ApiErrorMsg(c, msg)
		return
	}

	title := &model.InvoiceTitle{
		UserId:        c.GetInt("id"),
		Name:          req.Name,
		BillingType:   req.BillingType,
		Title:         req.Title,
		TaxId:         req.TaxId,
		Email:         req.Email,
		Street:        req.Street,
		AddressDetail: req.AddressDetail,
		City:          req.City,
		ZipCode:       req.ZipCode,
		Country:       req.Country,
		IsDefault:     req.IsDefault,
	}
	if err := title.Insert(); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, title)
}

// UpdateInvoiceTitle 更新发票抬头卡片
func UpdateInvoiceTitle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ApiErrorMsg(c, "id 参数无效")
		return
	}
	userId := c.GetInt("id")
	title, err := model.GetUserInvoiceTitleById(id, userId)
	if err != nil {
		common.ApiErrorMsg(c, "抬头卡片不存在")
		return
	}

	var req InvoiceTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "参数错误")
		return
	}
	normalizeInvoiceTitleRequest(&req)
	if msg := validateInvoiceTitleRequest(req); msg != "" {
		common.ApiErrorMsg(c, msg)
		return
	}

	title.Name = req.Name
	title.BillingType = req.BillingType
	title.Title = req.Title
	title.TaxId = req.TaxId
	title.Email = req.Email
	title.Street = req.Street
	title.AddressDetail = req.AddressDetail
	title.City = req.City
	title.ZipCode = req.ZipCode
	title.Country = req.Country
	title.IsDefault = req.IsDefault
	if err := title.Update(); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, title)
}

// DeleteInvoiceTitle 删除发票抬头卡片
func DeleteInvoiceTitle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ApiErrorMsg(c, "id 参数无效")
		return
	}
	if err := model.DeleteUserInvoiceTitleById(id, c.GetInt("id")); err != nil {
		common.ApiError(c, err)
		return
	}
	common.ApiSuccess(c, nil)
}

// GetMyInvoices 用户获取自己的发票列表
func GetMyInvoices(c *gin.Context) {
	userId := c.GetInt("id")
	pageInfo := common.GetPageQuery(c)

	invoices, total, err := model.GetUserInvoices(userId, pageInfo)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	enrichInvoicePaymentSnapshots(invoices)
	normalizeInvoiceDownloadUrls(invoices)
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(invoices)
	common.ApiSuccess(c, pageInfo)
}

// AdminListInvoices 管理员获取发票列表
func AdminListInvoices(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)
	status := c.Query("status")
	paymentMethod := c.Query("payment_method")
	keyword := c.Query("keyword")

	invoices, total, err := model.GetAllInvoices(status, paymentMethod, keyword, pageInfo)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	enrichInvoicePaymentSnapshots(invoices)
	normalizeInvoiceDownloadUrls(invoices)
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(invoices)
	common.ApiSuccess(c, pageInfo)
}

func normalizeInvoiceDownloadUrls(invoices []*model.Invoice) {
	for _, invoice := range invoices {
		if invoice == nil || invoice.Status != model.InvoiceStatusCompleted || invoice.DownloadUrl == "" {
			continue
		}
		if strings.HasPrefix(invoice.DownloadUrl, "minio:") || !strings.HasPrefix(invoice.DownloadUrl, "http") {
			invoice.DownloadUrl = fmt.Sprintf("/api/invoice/download?id=%d", invoice.Id)
		}
	}
}

func enrichInvoicePaymentSnapshots(invoices []*model.Invoice) {
	for _, invoice := range invoices {
		if invoice == nil {
			continue
		}

		topUp := model.GetTopUpById(invoice.TopUpId)
		if topUp == nil {
			invoice.PaidAmount = invoice.Money
			invoice.PaidCurrency = fallbackInvoiceCurrency(invoice.PaymentMethod)
			continue
		}

		invoice.CreditAmountUsd = topUp.GetCreditAmountUSD()
		invoice.PaidAmount = topUp.GetPaidAmount()
		if invoice.PaidAmount <= 0 {
			invoice.PaidAmount = invoice.Money
		}
		invoice.PaidCurrency = topUp.GetPaidCurrency(fallbackInvoiceCurrency(invoice.PaymentMethod))
		invoice.ExchangeRate = topUp.ExchangeRate
	}
}

func fallbackInvoiceCurrency(paymentMethod string) string {
	switch paymentMethod {
	case model.PaymentMethodPayPal:
		return "EUR"
	case model.PaymentMethodAlipay:
		return "CNY"
	default:
		return "USD"
	}
}

// allowedInvoiceExts 允许的发票文件扩展名
var allowedInvoiceExts = map[string]bool{
	".pdf":  true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

// AdminUploadInvoiceFile 管理员上传发票文件
func AdminUploadInvoiceFile(c *gin.Context) {
	invoiceIdStr := c.PostForm("invoice_id")
	if invoiceIdStr == "" {
		common.ApiErrorMsg(c, "invoice_id 不能为空")
		return
	}
	invoiceId, err := strconv.Atoi(invoiceIdStr)
	if err != nil {
		common.ApiErrorMsg(c, "invoice_id 参数无效")
		return
	}

	// 校验 invoice 存在
	invoice, err := model.GetInvoiceById(invoiceId)
	if err != nil {
		common.ApiErrorMsg(c, "发票记录不存在")
		return
	}
	_ = invoice

	file, err := c.FormFile("file")
	if err != nil {
		common.ApiErrorMsg(c, "请上传文件")
		return
	}

	// 限制文件大小 10MB
	if file.Size > 10*1024*1024 {
		common.ApiErrorMsg(c, "文件大小不能超过 10MB")
		return
	}

	// 校验文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedInvoiceExts[ext] {
		common.ApiErrorMsg(c, "只允许上传 PDF、PNG、JPG、JPEG 格式的文件")
		return
	}

	// 确保上传目录存在
	uploadDir := "./uploads/invoices"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		common.ApiErrorMsg(c, "创建上传目录失败")
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("invoice_%d_%d%s", invoiceId, time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		common.ApiErrorMsg(c, "保存文件失败")
		return
	}

	common.ApiSuccess(c, gin.H{"file_path": filePath})
}

// DownloadInvoiceFile 下载发票文件
func DownloadInvoiceFile(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		common.ApiErrorMsg(c, "id 参数不能为空")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ApiErrorMsg(c, "id 参数无效")
		return
	}

	invoice, err := model.GetInvoiceById(id)
	if err != nil {
		common.ApiErrorMsg(c, "发票记录不存在")
		return
	}

	userId := c.GetInt("id")
	role := c.GetInt("role")

	// 权限检查：非管理员只能下载自己的发票
	if role < common.RoleAdminUser && invoice.UserId != userId {
		common.ApiErrorMsg(c, "无权下载此发票")
		return
	}

	if invoice.DownloadUrl == "" {
		common.ApiErrorMsg(c, "发票文件尚未上传")
		return
	}

	if strings.HasPrefix(invoice.DownloadUrl, "minio:") {
		objectName := strings.TrimPrefix(invoice.DownloadUrl, "minio:")
		reader, contentLength, contentType, err := service.GetInvoicePDFReader(c.Request.Context(), objectName)
		if err != nil {
			common.ApiErrorMsg(c, "发票文件不存在")
			return
		}
		defer reader.Close()

		c.Header("Content-Type", contentType)
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="invoice_%d.pdf"`, invoice.Id))
		if contentLength >= 0 {
			c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
		}
		_, _ = io.Copy(c.Writer, reader)
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(invoice.DownloadUrl); os.IsNotExist(err) {
		common.ApiErrorMsg(c, "发票文件不存在")
		return
	}

	c.File(invoice.DownloadUrl)
}

type AdminCompleteInvoiceRequest struct {
	InvoiceId   int    `json:"invoice_id"`
	DownloadUrl string `json:"download_url"`
}

// AdminCompleteInvoice 管理员确认开票
func AdminCompleteInvoice(c *gin.Context) {
	var req AdminCompleteInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "参数错误")
		return
	}

	invoice, err := model.GetInvoiceById(req.InvoiceId)
	if err != nil {
		common.ApiErrorMsg(c, "发票记录不存在")
		return
	}

	if invoice.Status != model.InvoiceStatusPending {
		common.ApiErrorMsg(c, "只有待处理的发票才能确认")
		return
	}

	// 如果是 PayPal 支付，尝试调用 Lexware 自动开票
	if invoice.PaymentMethod == "paypal" {
		lexwareObjectName, lexErr := service.CreateLexwareInvoice(c.Request.Context(), invoice)
		if lexErr != nil {
			if req.DownloadUrl == "" {
				common.ApiErrorMsg(c, "Lexware 自动开票失败，请手动上传发票文件："+lexErr.Error())
				return
			}
		} else if lexwareObjectName != "" {
			req.DownloadUrl = "minio:" + lexwareObjectName
		}
	}

	if req.DownloadUrl == "" {
		common.ApiErrorMsg(c, "请提供发票文件路径")
		return
	}

	invoice.Status = model.InvoiceStatusCompleted
	invoice.DownloadUrl = req.DownloadUrl
	invoice.CompleteTime = time.Now().Unix()

	if err := invoice.Update(); err != nil {
		common.ApiError(c, err)
		return
	}

	// 异步发送邮件通知用户
	go func() {
		user, userErr := model.GetUserById(invoice.UserId, false)
		if userErr != nil || user.Email == "" {
			return
		}
		subject := "您的发票已开具"
		content := fmt.Sprintf("您好 %s，\n\n您申请的发票（订单号：%s）已开具完成，请登录系统下载。\n\n感谢您的使用！", user.Username, invoice.TradeNo)
		_ = common.SendEmail(subject, user.Email, content)
	}()

	common.ApiSuccess(c, nil)
}

type AdminRejectInvoiceRequest struct {
	InvoiceId int    `json:"invoice_id"`
	Message   string `json:"message"`
}

// AdminRejectInvoice 管理员拒绝开票
func AdminRejectInvoice(c *gin.Context) {
	var req AdminRejectInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ApiErrorMsg(c, "参数错误")
		return
	}

	invoice, err := model.GetInvoiceById(req.InvoiceId)
	if err != nil {
		common.ApiErrorMsg(c, "发票记录不存在")
		return
	}

	if invoice.Status != model.InvoiceStatusPending {
		common.ApiErrorMsg(c, "只有待处理的发票才能拒绝")
		return
	}

	invoice.Status = model.InvoiceStatusRejected
	invoice.Message = req.Message
	invoice.CompleteTime = time.Now().Unix()

	if err := invoice.Update(); err != nil {
		common.ApiError(c, err)
		return
	}

	// 异步发送邮件通知用户
	go func() {
		user, userErr := model.GetUserById(invoice.UserId, false)
		if userErr != nil || user.Email == "" {
			return
		}
		subject := "您的发票申请已被拒绝"
		content := fmt.Sprintf("您好 %s，\n\n您申请的发票（订单号：%s）已被拒绝。\n原因：%s\n\n如有疑问，请联系管理员。", user.Username, invoice.TradeNo, req.Message)
		_ = common.SendEmail(subject, user.Email, content)
	}()

	common.ApiSuccess(c, nil)
}
