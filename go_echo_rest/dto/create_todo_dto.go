package dto

type CreateTodoDTO struct {
	Title       string `json:"name"`
	Description string `json:"description"`
}
