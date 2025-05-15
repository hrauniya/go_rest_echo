package main

import (
	"go_echo_rest/config"
	"go_echo_rest/controller"
	"go_echo_rest/model"
	"net/http"
	"go_echo_rest/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Welcome to ToDo API",
		})
	})

	config.DatabaseInit()
	gorm := config.DB()

	gorm.AutoMigrate(&model.User{}, &model.ToDo{})

	dbGorm, err := gorm.DB()
	if err != nil {
		panic(err)
	}

	dbGorm.Ping()

	e.POST("/register", controller.Register)
	e.POST("/login", controller.Login)

	todoRoute := e.Group("/todo")
	todoRoute.Use(middleware.JWTMiddleware)
	todoRoute.POST("/", controller.CreateTodo)
	todoRoute.GET("/:id", controller.GetTodo)
	todoRoute.PUT("/:id", controller.UpdateTodo)
	todoRoute.DELETE("/:id", controller.DeleteTodo)
	e.Logger.Fatal(e.Start(":8080"))

}
