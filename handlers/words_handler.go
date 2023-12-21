package handlers

import (
	"net/http"
	"say-it/helper"
	"say-it/models"
	"github.com/gorilla/mux"
	"strconv"
)


func GetAllWordsHandler(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["page"]
	pageNum, err := strconv.Atoi(param)
	if err != nil {
		response := models.NewErrorResponse("Failed to get words list", "The parameter provided is invalid", "Invalid parameter")
		helper.WriteToResponseBody(w, http.StatusBadRequest, &response)
		return
	}

	minID := 1 + (pageNum - 1) * 16
	maxID := pageNum * 16

	rows, err := db.Query("SELECT id, word FROM words where id >= $1 AND id <= $2", minID, maxID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var id int
		var word string
		err := rows.Scan(&id, &word)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		words = append(words, models.Word{ID: id, WordTxt: word})
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
	ParamID := mux.Vars(r)["WordID"]
	var wordTXT string
	var wordDetail string

	wordID, err := strconv.Atoi(ParamID)
	if err != nil {
		response := models.NewErrorResponse("Failed to get word detail", "The ID provided is invalid", "Invalid ID")
		helper.WriteToResponseBody(w, http.StatusBadRequest, &response)
		return
	}

	// Cari Word ID di DB
	// Sound link nanti lagi karena untuk sekarang sumber data suara belum ada
	err = db.QueryRow("SELECT id, word, description FROM words WHERE id=$1", wordID).
		Scan(&wordID, &wordTXT, &wordDetail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var wordResponse models.WordDetail

	wordResponse.WordID = wordID
	wordResponse.WordTxt = wordTXT
	wordResponse.Description = wordDetail

	response := models.NewSuccessResponse("OK", wordResponse)
	helper.WriteToResponseBody(w, http.StatusOK, response)
}
