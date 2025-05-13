package controller


import (
	"net/http"
	"go_echo_rest/config"
	"go_echo_rest/model"
	"github.com/labstack/echo/v4"
)

//Creating a new todo task
func CreateToDo(c echo.Context) error {
	b := new(model.ToDo)
	db := config.DB()

	if err := c.Bind(b); err!=nil{
		dats := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	todo := &model.ToDo{
		Title: b.Title
		Description: b.Description,
	}

	if err:= db.Create(&todo).Error; err!=nil {
		data := map[string]interface{}{
			"message" : err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data":b,
	}

	return c.JSON(http.StatusOK, response)
}

//Get Todo
func GetTodo(c echo.Context) error {
		id := c.Param("id")
		db := config.DB()
	
		var todos []*model.ToDo
	
		if res := db.Find(&todos, id); res.Error != nil {
			data := map[string]interface{}{
				"message": res.Error.Error(),
			}
	
			return c.JSON(http.StatusOK, data)
		}
	
		response := map[string]interface{}{
			"data": todos[0],
		}
	
		return c.JSON(http.StatusOK, response)
}