package message

import (
	"encoding/json"
	"errors"
	"gitgub.com/bhimani07/messenger-go/db/models"
	"gitgub.com/bhimani07/messenger-go/utils"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Message struct {
	SenderId       uint32
	ReceiverId     uint32
	Text           string
	ConversationId uint32
}

func SaveMessage(w http.ResponseWriter, r *http.Request) {
	var body Message

	if err := utils.DecodeJSONBody(w, r, &body); err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	db, ok := r.Context().Value("db").(*gorm.DB)
	if !ok {
		http.Error(w, "Internal Server error, try after some time", 500)
		return
	}

	var conversation models.Conversation
	if body.ConversationId != 0 {
		// if conversation specified in request
		db.First(&conversation, body.ConversationId)

		isUserInConversation := conversation.User1Id == body.SenderId || conversation.User2Id == body.SenderId
		if !isUserInConversation {
			http.Error(w, "unauthorized access", 403)
			return
		}
	}

	if conversation.ID == 0 {
		// if no conversation specified in request

		// checking to see if these users already have a conversation
		db.Where(models.Conversation{
			User1Id: body.SenderId,
			User2Id: body.ReceiverId,
		}).Or(models.Conversation{
			User1Id: body.ReceiverId,
			User2Id: body.SenderId,
		}).First(&conversation)

		if conversation.ID == 0 {
			// if they never talked before create a new conversation
			conversation = models.Conversation{
				User1Id: body.SenderId,
				User2Id: body.ReceiverId,
			}
			res := db.Create(&conversation)

			if res.Error != nil {
				http.Error(w, "error while saving new conversation", 500)
				return
			}
		}
	}

	message := models.Message{
		SenderId:       body.SenderId,
		Text:           body.Text,
		IsSeen:         false,
		ConversationID: conversation.ID,
	}

	if res := db.Create(&message); res.Error != nil {
		http.Error(w, "error while saving new conversation/message", 500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(message)
}

func MarkSeenMessage(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ConversationId uint32
		SeenByUserId   uint32
	}{}

	if err := utils.DecodeJSONBody(w, r, &body); err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	db, ok := r.Context().Value("db").(*gorm.DB)
	if !ok {
		http.Error(w, "Internal Server error, try after some time", 500)
		return
	}

	var dbConversation models.Conversation
	if result := db.First(&dbConversation, body.ConversationId); result.Error != nil {
		http.Error(w, "no such conversation exist", http.StatusNotFound)
		return
	}

	if !(dbConversation.User1Id == body.SeenByUserId || dbConversation.User2Id == body.SeenByUserId) {
		http.Error(w, "unauthorized operation", http.StatusUnauthorized)
		return
	}

	if res := db.Model(&Message{}).Where("conversation_id = ? AND is_seen = ? AND sender_id != ?", body.ConversationId, false, body.SeenByUserId).Update("is_seen", true); res.Error != nil {
		http.Error(w, "internal server error while updating message seen", http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Success bool `json:"success"`
	}{Success: true})
}
