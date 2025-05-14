package controller

import (
	"go_echo_rest/config"
	"go_echo_rest/model"
	"go_echo_rest/dto"
	"log"
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
)

//Create new todo task
func CreateTodo(c echo.Context) error {
	 
	db := config.DB()
	todoDTO := new(dto.CreateTodoDTO)

	if err := c.Bind(&todoDTO); err!=nil{
		log.Println("Invalid request body passed in by client", err)
		data := map[string]interface{}{
			"message": "Invalid request body",
		}
		return c.JSON(http.StatusBadRequest, data)
	}
	
	todo := model.ToDo{
		Title : todoDTO.Title,
		Description : todoDTO.Description,
	}

	if err:= db.Create(&todo).Error; err!=nil {
		log.Printf("An error occured while creating record", err)
		data := map[string]interface{}{
			"message" : "An error occured while trying to create your record",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}
	return SuccessResponse(c,todo)
}

//Get specific todo task
func GetTodo(c echo.Context) error {

			//sanitize id parameter, check for validity
			requestId := c.Param("id")
			id, err := strconv.Atoi(requestId)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": "Invalid ID parameter",
				})
			}
			db := config.DB()

			var todo model.ToDo
		
			if res := db.First(&todo, id); res.Error != nil {
				data := map[string]interface{}{
					"message" : "Record not found!",
				}
				return c.JSON(http.StatusNotFound, data)
			}
		
			return SuccessResponse(c, todo)
}

//Update specific todo task
func UpdateTodo(c echo.Context) error {

	requestId := c.Param("id")
	id, err := strconv.Atoi(requestId)
	if err != nil {
		log.Printf("Invalid ID parameter was passed", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid ID parameter",
		})
	}
	updateTodoDTO:= new(dto.UpdateTodoDTO)
	db := config.DB()

	if err := c.Bind(&updateTodoDTO); err != nil {
		log.Printf("Invalid request body", err)
		data := map[string]interface{}{
			"message": "Invalid request body",
		}
		return c.JSON(http.StatusBadRequest, data)
	}

	var existing_todo model.ToDo

	if err := db.First(&existing_todo, id).Error; err != nil {
		log.Printf("Record not found", err)
		data := map[string]interface{}{
			"message": "Request record was not found in table",
		}
		return c.JSON(http.StatusNotFound, data)
	}

	existing_todo.Title = updateTodoDTO.Title
	existing_todo.Description = updateTodoDTO.Description
	existing_todo.Status = updateTodoDTO.Status
	if err := db.Save(&existing_todo).Error; err != nil {
		log.Printf("Updated record could not be saved", err)
		data := map[string]interface{}{
			"message": "Could not save updated record in table",
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