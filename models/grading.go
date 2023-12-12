package models

type Grading struct{
	UserID 	int 	`json:"id"`
	WordID 		string 	`json:"wordId"`
	FileLink 	string 	`json:"fileLink"`
}