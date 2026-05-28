package model

import (
	"errors"

	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
)

type Invoice struct {
	Id              int     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId          int     `json:"user_id" gorm:"index"`
	Username        string  `json:"username" gorm:"type:varchar(100)"`
	TopUpId         int     `json:"topup_id" gorm:"index"`
	TradeNo         string  `json:"trade_no" gorm:"type:varchar(255);uniqueIndex"`
	Money           float64 `json:"money"`
	CreditAmountUsd float64 `json:"credit_amount_usd" gorm:"-"`
	PaidAmount      float64 `json:"paid_amount" gorm:"-"`
	PaidCurrency    string  `json:"paid_currency" gorm:"-"`
	ExchangeRate    float64 `json:"exchange_rate" gorm:"-"`
	PaymentMethod   string  `json:"payment_method" gorm:"type:varchar(50)"`
	BillingType     string  `json:"billing_type" gorm:"type:varchar(20)"`
	Title           string  `json:"title" gorm:"type:varchar(255)"`
	TaxId           string  `json:"tax_id" gorm:"type:varchar(100)"`
	Email           string  `json:"email" gorm:"type:varchar(255)"`
	Street          string  `json:"street" gorm:"type:varchar(255)"`
	AddressDetail   string  `json:"address_detail" gorm:"type:varchar(255)"`
	City            string  `json:"city" gorm:"type:varchar(100)"`
	ZipCode         string  `json:"zip_code" gorm:"type:varchar(20)"`
	Country         string  `json:"country" gorm:"type:varchar(100)"`
	Status          string  `json:"status" gorm:"type:varchar(20);index;default:'pending'"`
	DownloadUrl     string  `json:"download_url" gorm:"type:varchar(512)"`
	Message         string  `json:"message" gorm:"type:text"`
	CreateTime      int64   `json:"create_time" gorm:"index"`
	CompleteTime    int64   `json:"complete_time"`
}

const (
	InvoiceStatusPending   = "pending"
	InvoiceStatusCompleted = "completed"
	InvoiceStatusRejected  = "rejected"
)

func (inv *Invoice) Insert() error {
	return DB.Create(inv).Error
}

func (inv *Invoice) Update() error {
	return DB.Save(inv).Error
}

func GetInvoiceById(id int) (*Invoice, error) {
	var invoice Invoice
	err := DB.Where("id = ?", id).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func GetInvoiceByTradeNo(tradeNo string) (*Invoice, error) {
	var invoice Invoice
	err := DB.Where("trade_no = ?", tradeNo).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func GetUserInvoices(userId int, pageInfo *common.PageInfo) ([]*Invoice, int64, error) {
	tx := DB.Begin()
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var total int64
	err := tx.Model(&Invoice{}).Where("user_id = ?", userId).Count(&total).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	var invoices []*Invoice
	err = tx.Where("user_id = ?", userId).Order("create_time desc").Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Find(&invoices).Error
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	if err = tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func GetAllInvoices(status string, paymentMethod string, keyword string, pageInfo *common.PageInfo) ([]*Invoice, int64, error) {
	tx := DB.Begin()
	if tx.Error != nil {
		return nil, 0, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := tx.Model(&Invoice{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if paymentMethod != "" {
		query = query.Where("payment_method = ?", paymentMethod)
	}
	if keyword != "" {
		pattern, perr := sanitizeLikePattern(keyword)
		if perr != nil {
			tx.Rollback()
			return nil, 0, perr
		}
		query = query.Where("trade_no LIKE ? ESCAPE '!' OR username LIKE ? ESCAPE '!'", pattern, pattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	var invoices []*Invoice
	if err := query.Order("create_time desc").Limit(pageInfo.GetPageSize()).Offset(pageInfo.GetStartIdx()).Find(&invoices).Error; err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

func GetInvoiceStatusByTopUpId(topUpId int) (string, error) {
	var invoice Invoice
	err := DB.Where("top_up_id = ?", topUpId).First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return invoice.Status, nil
}
