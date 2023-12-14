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
	rows, err := db.Query("SELECT id, word FROM words;")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer rows.Close()

	var words []Word
	for rows.Next() {
		err != rows.Scan(&id, &word)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		words = append(words, Word{ID: id, WordTXT: word})
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := models.NewSuccessResponse("OK", words)
	helper.WriteToResponseBody(w, http.StatusOK, response)
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
	if err != nil {
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
