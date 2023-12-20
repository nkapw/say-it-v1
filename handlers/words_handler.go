package handlers

import (
	"net/http"
	"say-it/helper"
	"say-it/models"
	"github.com/gorilla/mux"
	"strconv"
)


func GetAllWordsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse page query parameter
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Number of items per page
	itemsPerPage := 10

	// Calculate offset based on page number
	offset := (page - 1) * itemsPerPage

	// Query to fetch paginated words from the database

	// Execute the query
	rows, err := db.Query("SELECT id, word FROM words ORDER BY id LIMIT $1 OFFSET $2", strconv.Itoa(itemsPerPage), strconv.Itoa(offset))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error querying database:", err)
		return
	}
	defer rows.Close()

	// Fetch words from the result set
	words := make([]models.Word, 0)
	for rows.Next() {
		var word models.Word

		err := rows.Scan(&word.ID, &word.WordTxt)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return
		}
		words = append(words, word)
	}

	// Create JSON response
	response := map[string]interface{}{
		"page":  page,
		"words": words,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	res := models.NewSuccessResponse("ok", response)
	helper.WriteToResponseBody(w, http.StatusOK, res)

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
