package model

import (
	"github.com/QuantumNous/new-api/common"

	"gorm.io/gorm"
)

type PlaygroundSession struct {
	Id            int            `json:"id"`
	UserId        int            `json:"user_id" gorm:"index;not null"`
	Title         string         `json:"title" gorm:"type:varchar(128);not null"`
	Messages      JSONValue      `json:"messages" gorm:"type:text"`
	SelectedImage JSONValue      `json:"selected_image" gorm:"type:text"`
	CreatedTime   int64          `json:"created_time" gorm:"bigint;index"`
	UpdatedTime   int64          `json:"updated_time" gorm:"bigint;index"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type PlaygroundSessionMySQL struct {
	Id            int            `json:"id"`
	UserId        int            `json:"user_id" gorm:"index;not null"`
	Title         string         `json:"title" gorm:"type:varchar(128);not null"`
	Messages      JSONValue      `json:"messages" gorm:"type:longtext"`
	SelectedImage JSONValue      `json:"selected_image" gorm:"type:longtext"`
	CreatedTime   int64          `json:"created_time" gorm:"bigint;index"`
	UpdatedTime   int64          `json:"updated_time" gorm:"bigint;index"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (PlaygroundSessionMySQL) TableName() string {
	return "playground_sessions"
}

func (s *PlaygroundSession) Insert() error {
	now := common.GetTimestamp()
	s.CreatedTime = now
	s.UpdatedTime = now
	if s.Title == "" {
		s.Title = "Untitled conversation"
	}
	if s.Messages == nil {
		s.Messages = JSONValue([]byte("[]"))
	}
	if s.SelectedImage == nil {
		s.SelectedImage = JSONValue([]byte("null"))
	}
	return DB.Create(s).Error
}

func (s *PlaygroundSession) Update() error {
	s.UpdatedTime = common.GetTimestamp()
	return DB.Model(&PlaygroundSession{}).
		Where("id = ? AND user_id = ?", s.Id, s.UserId).
		Updates(map[string]any{
			"title":          s.Title,
			"messages":       s.Messages,
			"selected_image": s.SelectedImage,
			"updated_time":   s.UpdatedTime,
		}).Error
}

func GetUserPlaygroundSessions(userId int, offset int, limit int) ([]*PlaygroundSession, error) {
	var sessions []*PlaygroundSession
	if limit <= 0 {
		limit = 50
	}
	if err := DB.Where("user_id = ?", userId).
		Order("updated_time DESC").
		Offset(offset).
		Limit(limit).
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func CountUserPlaygroundSessions(userId int) (int64, error) {
	var count int64
	err := DB.Model(&PlaygroundSession{}).Where("user_id = ?", userId).Count(&count).Error
	return count, err
}

func GetUserPlaygroundSessionByID(id int, userId int) (*PlaygroundSession, error) {
	var session PlaygroundSession
	if err := DB.Where("id = ? AND user_id = ?", id, userId).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func DeleteUserPlaygroundSessionByID(id int, userId int) error {
	return DB.Where("id = ? AND user_id = ?", id, userId).Delete(&PlaygroundSession{}).Error
}
