package handlers

//
//import (
//	"encoding/json"
//	"github.com/gorilla/mux"
//	"net/http"
//	"say-it/models"
//)
//
//func GetUserHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Add("Content-Type", "application/json")
//	vars := mux.Vars(r)
//	userID := vars["id"]
//
//	var user models.User
//	err := db.QueryRow("SELECT id, name, email, password, username FROM users WHERE id=$1", userID).
//		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Username)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusNotFound)
//		w.WriteHeader(http.StatusInternalServerError)
//
//		res := map[string]string{
//			"status":  "not found",
//			"message": "user not found",
//		}
//		json.NewEncoder(w).Encode(&res)
//		return
//	}
//
//	json.NewEncoder(w).Encode(user)
//}
//
//func EditUserHandler(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	userID := vars["id"]
//
//	var updatedUser models.User
//	err := json.NewDecoder(r.Body).Decode(&updatedUser)
//	if err != nil {
//		//http.Error(w, err.Error(), http.StatusBadRequest)
//		w.WriteHeader(http.StatusBadRequest)
//
//		res := map[string]string{
//			"status": "bad request",
//			//"message": "email or username already exist",
//		}
//		json.NewEncoder(w).Encode(&res)
//
//		return
//	}
//
//	_, err = db.Exec("UPDATE users SET name=$1, email=$2, password=$3, username=$4 WHERE id=$5",
//		updatedUser.Name, updatedUser.Email, updatedUser.Password, updatedUser.Username, userID)
//	if err != nil {
//		//http.Error(w, err.Error(), http.StatusInternalServerError)
//		w.WriteHeader(http.StatusInternalServerError)
//
//		res := map[string]string{
//			"status": "internal server error",
//			//"message": "email or username already exist",
//		}
//		json.NewEncoder(w).Encode(&res)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//}
