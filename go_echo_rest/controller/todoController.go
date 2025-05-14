package controller


import (
	"net/http"
	"go_echo_rest/config"
	"go_echo_rest/model"
	"github.com/labstack/echo/v4"
	"log"
)

//Creating a new todo task
func CreateTodo(c echo.Context) error {
	var todo model.ToDo 
	db := config.DB()

	if err := c.Bind(&todo); err!=nil{
		log.Println("An error occured", err)
		data := map[string]interface{}{
			"message": "Invalid request body",
		}
		return c.JSON(http.StatusBadRequest, data)
	}

	if err:= db.Create(&todo).Error; err!=nil {
		data := map[string]interface{}{
			"message" : err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": todo,
	}

	return c.JSON(http.StatusOK, response)
}

//Get Todo
func GetTodo(c echo.Context) error {
		id := c.Param("id")
		fmt.Println(id)
		db := config.DB()
	
		var todo model.ToDo
	
		if res := db.First(&todos, id); res.Error != nil {
		
		}
	
		return SuccessResponse(c, todo)
}

//update 
func UpdateTodo(c echo.Context) error {
	id := c.Param("id")
	b := new(model.ToDo)
	db := config.DB()

	// Binding data
	if err := c.Bind(b); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	existing_todo := new(model.ToDo)

	if err := db.First(&existing_todo, id).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}

	existing_todo.Title = b.Title
	existing_todo.Description = b.Description
	existing_todo.Status = b.Status
	if err := db.Save(&existing_todo).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	return SuccessResponse(c, existing_todo)
}

func SuccessResponse(c, data any) nil{
	response := map[string]interface{}{
		"info": data,
	}
	return c.JSON(http.StatusOK, response)
}