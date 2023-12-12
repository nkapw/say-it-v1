package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"say-it/connection"
	"say-it/helper"
	"say-it/models"
	"strconv"
)

func GetAllWordsHandler(w http.ResponseWriter, r *http.Request) {
  
}

func GetWordDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan Word ID dari URL  Endpoint
	wordID := mux.Vars(r)["WordID"]
	var wordTXT string
	var wordDetail string

	// Cari Word ID di DB
	// Sound link nanti lagi karena untuk sekarang sumber data suara belum ada
	err = db.QueryRow("SELECT id, word, detail FROM words WHERE id=$1", wordID).
		Scan(&wordID, &wordTXT, &wordDetail)
	if err := nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var wordResponse models.WordDetail

	wordResponse.WordID = wordID
	wordResponse.WordTXT = wordTXT
	wordResponse.Description = wordDetail

	response := models.NewSuccessResponse("OK", wordResponse)
	helper.WriteToResponseBody(w, http.StatusOK, response)
}
