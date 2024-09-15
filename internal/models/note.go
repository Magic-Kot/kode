package models

type NotesResponse struct {
	Id    int    `json:"id" db:"id"`
	Texts string `json:"texts" db:"texts"`
}

type TextsInput struct {
	Texts []string `json:"texts"`
}
