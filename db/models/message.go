package models

import "time"

type Message struct {
	ID             uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Text           string `gorm:"size:255" json:"text"`
	SenderId       uint32 `json:"senderId"`
	IsSeen         bool   `gorm:"default:false" json:"isSeen"`
	ConversationID uint32 `json:"conversationId"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	conversation   Conversation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;NOT NULL;references:ID"`
	user           User         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;NOT NULL;foreignKey:SenderId;references:ID"`
}
