package models

// models/word.go
package models

type Word struct {
	ID    int    `json:"id"`
	WordTxt  string `json:"word_txt"`
	VoiceUrl string `json:"voice_url"`
}


  type WordDetail struct {
	WordID      int    	`json:"id"`
	WordTxt     string	`json:"word"`
	Description string 	`json:"description"`
}
