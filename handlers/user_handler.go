package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"say-it/helper"
	"say-it/models"
)

func UpdateCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan ID pengguna dari token
	userID, err := helper.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Mendapatkan data pengguna dari database
	var currentUser models.User
	var currentUserPfp sql.NullString

	if currentUserPfp.Valid {
		currentUser.ProfilePicture = currentUserPfp.String
	}
	err = db.QueryRow("SELECT id, name, email, password, username, profile_picture FROM users WHERE id=$1", userID).
		Scan(&currentUser.ID, &currentUser.Name, &currentUser.Email, &currentUser.Password, &currentUser.Username, &currentUserPfp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Mendapatkan data yang diperbarui dari formulir multipart
	err = r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mengambil file gambar dari formulir
	file, header, err := r.FormFile("profile_picture")
	if err == nil {
		// Jika ada file gambar, simpan di Google Cloud Storage
		defer file.Close()

		// Inisialisasi klien GCS
		gcsClient, err := helper.CreateGCSClient()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gcsClient.Close()

		// Simpan gambar di GCS
		bucketName := "profile_picture_bucket"
		objectName := fmt.Sprintf("profile_%d_%s", userID, header.Filename)

		ctx := context.Background()
		wc := gcsClient.Bucket(bucketName).Object(objectName).NewWriter(ctx)
		_, err = io.Copy(wc, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := wc.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Dapatkan URL gambar GCS
		imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

		// Simpan URL gambar di database
		_, err = db.Exec("UPDATE users SET profile_picture=$1 WHERE id=$2", imageURL, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Update informasi pengguna
		currentUser.ProfilePicture = imageURL
		response := models.NewSuccessResponse("User information updated successfully", currentUser)
		helper.WriteToResponseBody(w, http.StatusOK, response)
	}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := helper.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var user models.User
	err = db.QueryRow("SELECT id, name, email, password, username, profile_picture FROM users WHERE id=$1", userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Username)
	if err != nil {
		response := models.NewErrorResponse("user not found", "not found")
		helper.WriteToResponseBody(w, http.StatusNotFound, &response)
		return
	}

	response := models.NewSuccessResponse("ok", user)
	helper.WriteToResponseBody(w, http.StatusOK, response)
}
