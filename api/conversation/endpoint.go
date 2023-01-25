package conversation

import (
	"encoding/json"
	"gitgub.com/bhimani07/messenger-go/db/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Conversation struct {
	ID                uint32           `json:"id"`
	LatestMessageText string           `json:"latestMessageText"`
	UnseenCount       uint32           `json:"unseenCount"`
	Messages          []models.Message `json:"messages"`
	OtherUser         models.User      `json:"otherUser"`
}

type queryConversation struct {
	MessageID      uint32
	ConversationID uint32
	User1Id        uint32
	User2Id        uint32
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Text           string
	SenderId       uint32
	IsSeen         bool
}

func GetConversations(w http.ResponseWriter, r *http.Request) {
	userId, ok := mux.Vars(r)["userId"]
	if !ok {
		http.Error(w, "Invalid/Incorrect request parameter", 400)
	}

	db, ok := r.Context().Value("db").(*gorm.DB)
	if !ok {
		http.Error(w, "Internal Server error, try after some time", 500)
		return
	}

	var dbQueryConversations []queryConversation
	if result := db.Model(&models.Conversation{}).Preload("User1", "id != (?)", userId).Preload("User2", "id != (?)", userId).Select("conversations.user1_id, conversations.user2_id, messages.id AS message_id, messages.text, messages.sender_id, messages.is_seen, messages.created_at, messages.updated_at, messages.conversation_id").Where("conversations.user1_id = ?", userId).Or("conversations.user2_id = ?", userId).Joins("left join messages on messages.conversation_id = conversations.id").Order("messages.created_at desc").Group("conversations.id, messages.id").Scan(&dbQueryConversations); result.Error != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	conversationGroup := make(map[uint32][]queryConversation, len(dbQueryConversations))
	for _, convo := range dbQueryConversations {
		if _, found := conversationGroup[convo.ConversationID]; !found {
			conversationGroup[convo.ConversationID] = append(conversationGroup[convo.ConversationID], convo)
		} else {
			conversationGroup[convo.ConversationID] = append(conversationGroup[convo.ConversationID], convo)
		}
	}

	var response = make([]Conversation, len(conversationGroup))
	var i = 0
	for key := range conversationGroup {

		var latestText string
		var unseenCount uint32
		var otherUser models.User
		var messages []models.Message
		for i, message := range conversationGroup[key] {
			if i == 0 {
				latestText = message.Text
			}
			if !message.IsSeen {
				unseenCount++
			}

			if otherUser.ID == 0 {
				var otherUserId uint32
				if message.User1Id != message.SenderId {
					otherUserId = message.SenderId
				} else {
					otherUserId = message.User2Id
				}
				db.First(&otherUser, otherUserId)
			}

			messages = append(messages, models.Message{
				ID:             message.MessageID,
				Text:           message.Text,
				SenderId:       message.SenderId,
				IsSeen:         message.IsSeen,
				ConversationID: key,
				CreatedAt:      message.CreatedAt,
				UpdatedAt:      message.UpdatedAt,
			})
		}

		response[i] = Conversation{
			ID:                key,
			LatestMessageText: latestText,
			UnseenCount:       unseenCount,
			Messages:          messages,
			OtherUser:         otherUser,
		}
		i++
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
