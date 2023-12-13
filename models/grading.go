package models

type Grading struct{
	UserID 	int 	`json:"id"`
	WordID 		string 	`json:"word_id"`
	FileLink 	string 	`json:"file_link"`
}