package model

import (
	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
)

type InvoiceTitle struct {
	Id            int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId        int    `json:"user_id" gorm:"index;not null"`
	Name          string `json:"name" gorm:"type:varchar(100);not null"`
	BillingType   string `json:"billing_type" gorm:"type:varchar(20);not null"`
	Title         string `json:"title" gorm:"type:varchar(255);not null"`
	TaxId         string `json:"tax_id" gorm:"type:varchar(100)"`
	Email         string `json:"email" gorm:"type:varchar(255);not null"`
	Street        string `json:"street" gorm:"type:varchar(255)"`
	AddressDetail string `json:"address_detail" gorm:"type:varchar(255)"`
	City          string `json:"city" gorm:"type:varchar(100)"`
	ZipCode       string `json:"zip_code" gorm:"type:varchar(20)"`
	Country       string `json:"country" gorm:"type:varchar(100)"`
	IsDefault     bool   `json:"is_default" gorm:"index;default:false"`
	CreateTime    int64  `json:"create_time" gorm:"index"`
	UpdateTime    int64  `json:"update_time"`
}

func (t *InvoiceTitle) Insert() error {
	now := common.GetTimestamp()
	t.CreateTime = now
	t.UpdateTime = now
	return DB.Transaction(func(tx *gorm.DB) error {
		if t.IsDefault {
			if err := tx.Model(&InvoiceTitle{}).Where("user_id = ?", t.UserId).Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(t).Error
	})
}

func (t *InvoiceTitle) Update() error {
	t.UpdateTime = common.GetTimestamp()
	return DB.Transaction(func(tx *gorm.DB) error {
		if t.IsDefault {
			if err := tx.Model(&InvoiceTitle{}).Where("user_id = ? AND id <> ?", t.UserId, t.Id).Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Save(t).Error
	})
}

func GetUserInvoiceTitles(userId int) ([]*InvoiceTitle, error) {
	var titles []*InvoiceTitle
	err := DB.Where("user_id = ?", userId).Order("is_default desc, update_time desc, id desc").Find(&titles).Error
	return titles, err
}

func GetUserInvoiceTitleById(id int, userId int) (*InvoiceTitle, error) {
	var title InvoiceTitle
	err := DB.Where("id = ? AND user_id = ?", id, userId).First(&title).Error
	if err != nil {
		return nil, err
	}
	return &title, nil
}

func DeleteUserInvoiceTitleById(id int, userId int) error {
	return DB.Where("id = ? AND user_id = ?", id, userId).Delete(&InvoiceTitle{}).Error
}
