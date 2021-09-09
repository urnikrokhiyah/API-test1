package controllers

import (
	"cleancode/lib/databases"
	"cleancode/middlewares"
	"cleancode/models"
	"cleancode/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllUsersController(c echo.Context) error {
	users, rowAffected, err := databases.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("failed"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponse("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success", users))
}

func GetSingleUserController(c echo.Context) error {
	userId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid user id"))
	}

	loggedUserId := middlewares.ExtractToken(c)
	if loggedUserId != userId {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("unauthorized"))
	}

	user, rowAffected, err := databases.GetSingleUser(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("failed"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponse("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success", user))
}

func CreateUserControllers(c echo.Context) error {
	var user models.User
	c.Bind(&user)

	newUser, err := databases.CreateNewUser(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success", newUser))
}

func DeleteUserController(c echo.Context) error {
	userId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid user id"))
	}

	loggedUserId := middlewares.ExtractToken(c)
	if loggedUserId != userId {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("unauthorized"))
	}

	message, rowAffected, err := databases.DeleteUser(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("failed"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponse("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success", message))
}

func UpdateUserController(c echo.Context) error {
	userId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("invalid user id"))
	}

	loggedUserId := middlewares.ExtractToken(c)
	if loggedUserId != userId {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("unauthorized"))
	}

	newUser := models.User{}
	c.Bind(&newUser)

	updatedUser, rowAffected, err := databases.UpdateUser(userId, newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("failed"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponse("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success", updatedUser))
}

func LoginUserController(c echo.Context) error {
	user := models.User{}
	c.Bind(&user)

	loggedUser, err := databases.LoginUsers(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success", loggedUser))
}

func GetUserDetailControllersTesting() echo.HandlerFunc {
	return GetSingleUserController
}

func UpdatedDetailUserTesting() echo.HandlerFunc {
	return UpdateUserController
}

func DeleteDetailUserTesting() echo.HandlerFunc {
	return DeleteUserController
}
