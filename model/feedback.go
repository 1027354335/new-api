package model

import (
	"errors"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
)

type Feedback struct {
	Id           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId       int    `json:"user_id" gorm:"index"`
	Username     string `json:"username" gorm:"type:varchar(100)"`
	Email        string `json:"email" gorm:"type:varchar(255)"`
	Title        string `json:"title" gorm:"type:varchar(255)"`
	Content      string `json:"content" gorm:"type:text"`
	Category     string `json:"category" gorm:"type:varchar(32);index"`
	Priority     string `json:"priority" gorm:"type:varchar(32);index"`
	Status       string `json:"status" gorm:"type:varchar(32);index;default:'open'"`
	AdminId      int    `json:"admin_id" gorm:"index"`
	AdminName    string `json:"admin_name" gorm:"type:varchar(100)"`
	Reply        string `json:"reply" gorm:"type:text"`
	CreateTime   int64  `json:"create_time" gorm:"index"`
	UpdateTime   int64  `json:"update_time" gorm:"index"`
	ResolvedTime int64  `json:"resolved_time"`
}

const (
	FeedbackStatusOpen       = "open"
	FeedbackStatusInProgress = "in_progress"
	FeedbackStatusResolved   = "resolved"
	FeedbackStatusRejected   = "rejected"
	FeedbackStatusClosed     = "closed"

	FeedbackCategoryBug     = "bug"
	FeedbackCategoryAccount = "account"
	FeedbackCategoryBilling = "billing"
	FeedbackCategoryModel   = "model"
	FeedbackCategoryOther   = "other"

	FeedbackPriorityLow    = "low"
	FeedbackPriorityNormal = "normal"
	FeedbackPriorityHigh   = "high"
	FeedbackPriorityUrgent = "urgent"
)

var feedbackStatuses = map[string]bool{
	FeedbackStatusOpen:       true,
	FeedbackStatusInProgress: true,
	FeedbackStatusResolved:   true,
	FeedbackStatusRejected:   true,
	FeedbackStatusClosed:     true,
}

var feedbackReplyStatuses = map[string]bool{
	FeedbackStatusResolved: true,
	FeedbackStatusRejected: true,
}

var feedbackCategories = map[string]bool{
	FeedbackCategoryBug:     true,
	FeedbackCategoryAccount: true,
	FeedbackCategoryBilling: true,
	FeedbackCategoryModel:   true,
	FeedbackCategoryOther:   true,
}

var feedbackPriorities = map[string]bool{
	FeedbackPriorityLow:    true,
	FeedbackPriorityNormal: true,
	FeedbackPriorityHigh:   true,
	FeedbackPriorityUrgent: true,
}

func IsValidFeedbackStatus(status string) bool {
	return feedbackStatuses[status]
}

func IsValidFeedbackReplyStatus(status string) bool {
	return feedbackReplyStatuses[status]
}

func IsValidFeedbackCategory(category string) bool {
	return feedbackCategories[category]
}

func IsValidFeedbackPriority(priority string) bool {
	return feedbackPriorities[priority]
}

func (feedback *Feedback) Insert() error {
	return DB.Create(feedback).Error
}

func (feedback *Feedback) Update() error {
	return DB.Save(feedback).Error
}

func GetFeedbackById(id int) (*Feedback, error) {
	var feedback Feedback
	err := DB.Where("id = ?", id).First(&feedback).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func GetUserFeedbackById(id int, userId int) (*Feedback, error) {
	var feedback Feedback
	err := DB.Where("id = ? AND user_id = ?", id, userId).First(&feedback).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func GetUserFeedbacks(userId int, status string, pageInfo *common.PageInfo) ([]*Feedback, int64, error) {
	query := DB.Model(&Feedback{}).Where("user_id = ?", userId)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	return findFeedbacks(query, pageInfo)
}

func GetAllFeedbacks(status string, category string, priority string, keyword string, pageInfo *common.PageInfo) ([]*Feedback, int64, error) {
	query := DB.Model(&Feedback{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if strings.TrimSpace(keyword) != "" {
		pattern, err := sanitizeLikePattern("%" + keyword + "%")
		if err != nil {
			return nil, 0, err
		}
		query = query.Where("title LIKE ? ESCAPE '!' OR content LIKE ? ESCAPE '!' OR username LIKE ? ESCAPE '!' OR email LIKE ? ESCAPE '!'", pattern, pattern, pattern, pattern)
	}
	return findFeedbacks(query, pageInfo)
}

func findFeedbacks(query *gorm.DB, pageInfo *common.PageInfo) ([]*Feedback, int64, error) {
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var feedbacks []*Feedback
	if err := query.Order("update_time desc, create_time desc").
		Limit(pageInfo.GetPageSize()).
		Offset(pageInfo.GetStartIdx()).
		Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

func IsFeedbackNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
