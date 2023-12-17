/**
	Step grading API
	1. Receive inputs with the following data 
	- Account ID : The account that is trying to grade
	- Word ID : The ID for the word that is being graded
	- FIle : The sound file from the used being graded
	2. Upload the file into google cloud storage
	3. Get the link for that file object in the google cloud storage
	4. Store the inputed data into SQL with the following format
	- Account ID
	- Word ID 
	- File Link : The link for the sound file in the google cloud storage
	5. Access the phonetic alphabet of the word being graded
	6. Pass the file into the model and get the phonetic alphabet as a result
	7. Compare the two and get an error number
	8. Grade user's performance based on the error
**/

package handlers

import (
	"fmt"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"say-it/helper"
	"say-it/models"
	"io"
)

func GradingHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan Word ID dari URL endpoint
	wordID := mux.Vars(r)["WordID"]

	// Mendapatkan ID pengguna dari token
	userID, err := helper.GetUserIDFromToken(r)
	if err != nil {
		response := models.NewErrorResponse("Failed to Process Grade ", "Unauthorized", "Invalid Token")
		helper.WriteToResponseBody(w, http.StatusUnauthorized, &response)
		return
	}

	// Parse data yang didapat
	err = r.ParseMultipartForm(10 << 20) // hanya membolehkan file size kurang dari 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mendapatkan file suara dari user
	file, header, err := r.FormFile("user_sound")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Persiapkan Google Cloud Storage
	defer file.Close()

	// Inisialisasi client GCS
	gcsClient, err := helper.CreateGCSClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer gcsClient.Close()	

	// simpan file di gcs
	bucketName := "say-it-grading-bucket"
	objectName := fmt.Sprintf("grading_%d_%s", userID, header.Filename)

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

	// Dapatkan URL file di GCS
	fileURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	// Buat entry baru ke dalam tabel grading
	_, err = db.Exec("INSERT INTO grading (userID, wordID, sound_link) VALUES ($1, $2, $3);", userID, wordID, fileURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var gradingResponse models.Grading
	gradingResponse.UserID = userID
	gradingResponse.WordID = wordID
	gradingResponse.FileLink = fileURL

	response := models.NewSuccessResponse("User Grading Requirements Have Been Saved", gradingResponse)
	helper.WriteToResponseBody(w, http.StatusOK, response)

}

