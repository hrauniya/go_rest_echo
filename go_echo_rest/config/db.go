package config

import (
	"fmt"
	"os"
	"log"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"go_echo_rest/model"
	
)

var database *gorm.DB 
var e error

func DatabaseInit(){
	err:= godotenv.Load()
	if err!=nil {
		log.Println("No env file found!", err)
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn:= fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", host, user, password, dbName, port)
	database, e = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if e!=nil{
		panic(e)
	}

	database.AutoMigrate(&model.ToDo{})
}

func DB() *gorm.DB{
	return database
}