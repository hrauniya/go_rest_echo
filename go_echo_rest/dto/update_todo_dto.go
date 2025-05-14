package dto

type UpdateTodoDTO struct {
	Title string `json:"name"`
	Description string `json:"description"`
	Status bool `json:"status"`
}

