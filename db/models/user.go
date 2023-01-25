package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Username  string `gorm:"size:255;not null;unique" json:"username"`
	Email     string `gorm:"size:100;not null;unique" json:"email"`
	PhotoUrl  string `gorm:"size:255" json:"photoUrl"`
	Password  string `gorm:"size:255" json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func validatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (user *User) BeforeCreate(db *gorm.DB) error {
	hashedPassword, err := Hash(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return nil
}
