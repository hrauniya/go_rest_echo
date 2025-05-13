package main

type ToDo struct {
	ID int `gorm:"primaryKey" json:"id"`
	Title      string `json:"name"`
	Description string `json:"description"`
	Status bool `json:"status"`
}