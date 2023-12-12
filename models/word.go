package models

// models/word.go
package models

type Word struct {
	ID    int    `json:"id"`
	WordTxt  string `json:"word_txt"`
	VoiceUrl string `json:"voice_url"`
}


  type WordDetail struct {
	ID          int    `json:"id"`
	WordID      int    `json:"word_id"`
	Description string `json:"description"`
}
