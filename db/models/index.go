package models

import (
	"gorm.io/gorm"
)

func CreateModels(dbInstance *gorm.DB) {
	dropTables(dbInstance)
	createTables(dbInstance)
	seedData(dbInstance)

}

func createTables(conn *gorm.DB) {
	err := conn.AutoMigrate(&User{})
	if err != nil {
		panic("Error while creating User table")
	}
	err = conn.AutoMigrate(&Message{})
	if err != nil {
		panic("Error while creating Message table")
	}
	err = conn.AutoMigrate(&Conversation{})
	if err != nil {
		panic("Error while creating User table")
	}
}

func dropTables(conn *gorm.DB) {
	err := conn.Migrator().DropTable(&User{})
	if err != nil {
		panic("Error while dropping User table")
	}
	err = conn.Migrator().DropTable(&Message{})
	if err != nil {
		panic("Error while dropping Message table")
	}
	err = conn.Migrator().DropTable(&Conversation{})
	if err != nil {
		panic("Error while dropping conversation table")
	}
}

func seedData(conn *gorm.DB) {
	thomas := User{
		Username: "thomas",
		Email:    "thomas@email.com",
		Password: "123456",
		PhotoUrl: "",
	}
	conn.Create(&thomas)

	santiago := User{
		Username: "santiago",
		Email:    "santiago@email.com",
		Password: "123456",
		PhotoUrl: "https://res.cloudinary.com/dmlvthmqr/image/upload/v1607914466/messenger/775db5e79c5294846949f1f55059b53317f51e30_s3back.png",
	}
	conn.Create(&santiago)

	santiagoConvo := Conversation{
		User1Id: thomas.ID,
		User2Id: santiago.ID,
	}
	conn.Create(&santiagoConvo)

	conn.Create(&Message{
		ConversationID: santiagoConvo.ID,
		SenderId:       santiago.ID,
		IsSeen:         false,
		Text:           "Where are you from?",
	})

	conn.Create(&Message{
		ConversationID: santiagoConvo.ID,
		SenderId:       thomas.ID,
		IsSeen:         false,
		Text:           "I'm from New York",
	})

	conn.Create(&Message{
		ConversationID: santiagoConvo.ID,
		SenderId:       santiago.ID,
		IsSeen:         false,
		Text:           "Share photo of your city, please",
	})

	chiumbo := User{
		Username: "chiumbo",
		Email:    "chiumbo@email.com",
		Password: "123456",
		PhotoUrl: "",
	}
	conn.Create(&chiumbo)

	chiumboConvo := Conversation{
		User1Id: thomas.ID,
		User2Id: chiumbo.ID,
	}
	conn.Create(&chiumboConvo)

	conn.Create(&Message{
		ConversationID: chiumboConvo.ID,
		SenderId:       chiumbo.ID,
		IsSeen:         false,
		Text:           "Sure! What time?",
	})

	conn.Create(&Message{
		ConversationID: chiumboConvo.ID,
		SenderId:       chiumbo.ID,
		IsSeen:         false,
		Text:           "a test message",
	})

	conn.Create(&Message{
		ConversationID: chiumboConvo.ID,
		SenderId:       chiumbo.ID,
		IsSeen:         false,
		Text:           "😂 😂 😂",
	})

}
