package controller

import (
	"go_echo_rest/config"
	"go_echo_rest/dto"
	"go_echo_rest/model"
	"log"
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
)

// Create new todo task
func CreateTodo(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	db := config.DB()
	todoDTO := new(dto.CreateTodoDTO)

	if err := c.Bind(&todoDTO); err != nil {
		log.Println("Invalid request body passed in by client", err)
		return FailResponse(c, http.StatusBadRequest,"Invalid Request Body")
	}

	todo := model.ToDo{
		Title:       todoDTO.Title,
		Description: todoDTO.Description,
		UserID:      userID,
	}

	if err := db.Create(&todo).Error; err != nil {
		log.Println("An error occured while creating record", err)
		return FailResponse(c, http.StatusInternalServerError, "An error occured while trying to create your record")
	}
	return SuccessResponse(c, todo)
}

// Get specific todo task
func GetTodo(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	//sanitize id parameter, check for validity
	requestId := c.Param("id")
	id, err := strconv.Atoi(requestId)
	if err != nil {
		return FailResponse(c, http.StatusBadRequest, "Invalid ID parameter")
	}
	db := config.DB()

	var todo model.ToDo

	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		return FailResponse(c, http.StatusNotFound, "Record Not found")
	}

	return SuccessResponse(c, todo)
}

// Update specific todo task
func UpdateTodo(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	requestId := c.Param("id")
	id, err := strconv.Atoi(requestId)
	if err != nil {
		log.Println("Invalid ID parameter was passed", err)
		return FailResponse(c, http.StatusBadRequest, "Invalid ID parameter")
	}
	updateTodoDTO := new(dto.UpdateTodoDTO)
	db := config.DB()

	if err := c.Bind(&updateTodoDTO); err != nil {
		log.Println("Invalid request body", err)
		return FailResponse(c, http.StatusBadRequest, "Invalid Request Body")
	}

	var existingTodo model.ToDo

	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&existingTodo).Error; err != nil {
		log.Println("Record not found", err)
		return FailResponse(c, http.StatusNotFound,  "Request record was not found in table")
	}

	existingTodo.Title = updateTodoDTO.Title
	existingTodo.Description = updateTodoDTO.Description
	existingTodo.Status = updateTodoDTO.Status

	if err := db.Save(&existingTodo).Error; err != nil {
		log.Println("Updated record could not be saved", err)
		return FailResponse(c,http.StatusInternalServerError, "Could not save updated record in table")
	}

	return SuccessResponse(c, existingTodo)
}

func DeleteTodo(c echo.Context) error {
	
	db := config.DB()
	userID := c.Get("user_id").(uint)
	requestId := c.Param("id")
	id, err := strconv.Atoi(requestId)
	if err != nil {
		return FailResponse(c, http.StatusBadRequest, "Invalid ID parameter")
	}

	result := db.Where("id = ? AND user_id = ?", id , userID).Delete(&model.ToDo{})
	if result.RowsAffected == 0 {
		return FailResponse(c, http.StatusNotFound, "Record was not found")
	}
	if result.Error != nil {
		return FailResponse(c, http.StatusInternalServerError, "An unexpected error occured")
	}
	return SuccessResponse(c, "Successfully deleted Todo")
}

func SuccessResponse(c echo.Context, data any) error {
	response := map[string]interface{}{
		"info": data,
	}
	return c.JSON(http.StatusOK, response)
}

func FailResponse(c echo.Context,status int, message any) error {
	response := map[string]interface{}{
		"message": message,
	}
	return c.JSON(status,response)
}


