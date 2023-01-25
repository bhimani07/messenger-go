package main

import (
	"context"
	"gitgub.com/bhimani07/messenger-go/api/conversation"
	"gitgub.com/bhimani07/messenger-go/api/message"
	"gitgub.com/bhimani07/messenger-go/api/user"
	"gitgub.com/bhimani07/messenger-go/db"
	"gitgub.com/bhimani07/messenger-go/db/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

//type Server struct {
//	DB     *gorm.DB
//	Router *mux.Router
//}

func InitializeAndStartServer() {
	db := db.GetDbInstance()
	router := mux.NewRouter()

	models.CreateModels(db)
	bindRoutes(db, router)

	http.ListenAndServe(":8082", router)
}

func injectDB(db *gorm.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "db", db)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func bindRoutes(db *gorm.DB, router *mux.Router) {
	router.HandleFunc("/user/{username}", injectDB(db, user.FindByUserName)).Methods(http.MethodGet)
	router.HandleFunc("/messages", injectDB(db, message.SaveMessage)).Methods(http.MethodPost)
	router.HandleFunc("/messages/markSeen", injectDB(db, message.MarkSeenMessage)).Methods(http.MethodPut)
	router.HandleFunc("/conversation/{userId}", injectDB(db, conversation.GetConversations)).Methods(http.MethodGet)
	router.HandleFunc("/auth/register", injectDB(db, user.RegisterUser)).Methods(http.MethodPost)
	router.HandleFunc("/health", test).Methods(http.MethodGet)

	//router.Use(mux.CORSMethodMiddleware(router))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("okay"))
}
