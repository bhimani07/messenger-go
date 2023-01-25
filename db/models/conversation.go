package models

import "time"

type Conversation struct {
	ID        uint32 `gorm:"primary_key;auto_increment" json:"id"`
	User1Id   uint32
	User2Id   uint32
	Messages  []*Message
	CreatedAt time.Time
	UpdatedAt time.Time
	User1     User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;NOT NULL;foreignKey:User1Id"`
	User2     User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;NOT NULL;foreignKey:User2Id"`
}
