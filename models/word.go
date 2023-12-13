package models

// models/word.go
package models

type Word struct {
	ID    int    `json:"id"`
	WordTxt  string `json:"word_txt"`
}


type WordDetail struct {
	WordID      int    	`json:"id"`
	WordTxt     string	`json:"word"`
	Description string 	`json:"description"`
	// VoiceUrl string `json:"voice_url"` // belum dipake karena voicenya belum ada
}
