package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"go_echo_rest/config"
	"go_echo_rest/controller"
)


func main() {
	e:= echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			 "hello": "world",
		})
	})

	config.DatabaseInit()
	gorm := config.DB()

	dbGorm, err := gorm.DB()
	if err!=nil{
		panic(err)
	}

	dbGorm.Ping()

	todoRoute := e.Group("/todo")
	todoRoute.POST("/", controller.CreateTodo)
	todoRoute.GET("/:id", controller.GetTodo)
	todoRoute.PUT("/:id", controller.UpdateTodo)
	// todoRoute.DELETE("/:id", controller.DeleteBook)
	e.Logger.Fatal(e.Start(":8080"))



}