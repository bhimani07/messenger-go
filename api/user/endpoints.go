package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitgub.com/bhimani07/messenger-go/db/models"
	"gitgub.com/bhimani07/messenger-go/utils"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/mail"
	"strings"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Username string
		Password string
		Email    string
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

	if len(strings.TrimSpace(body.Username)) != 0 || len(strings.TrimSpace(body.Password)) != 0 || !isValidEmail(body.Email) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	fmt.Println("body: ", body)

	db, ok := r.Context().Value("db").(*gorm.DB)
	if !ok {
		http.Error(w, "Internal Server error, try after some time", 500)
	}

	user := &models.User{
		Username: strings.TrimSpace(body.Username),
		Password: strings.TrimSpace(body.Username),
		Email:    strings.TrimSpace(body.Email),
	}

	db.Create(user)

}

func FindByUserName(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value("db").(*gorm.DB)
	if !ok {
		http.Error(w, "Internal Server error, try after some time", 500)
	}

	username, ok := mux.Vars(r)["username"]
	if !ok {
		http.Error(w, "Invalid/Incorrect request parameter", 400)
	}

	var user models.User
	if result := db.Where("username = ?", username).First(&user); result.Error != nil {
		log.Print(result.Error)
		http.Error(w, "No such user with username", 404)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func isValidEmail(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}
