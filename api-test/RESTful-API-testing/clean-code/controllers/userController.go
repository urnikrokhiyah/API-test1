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
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponse("can't get all users"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponse("failed to get all user data"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success get all users", users))
}

func GetSingleUserController(c echo.Context) error {
	userId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponse("invalid user id"))
	}

	loggedUserId := middlewares.ExtractToken(c)
	if loggedUserId != userId {
		return c.JSON(http.StatusOK, response.ErrorResponse("can't access other accounts"))
	}

	user, rowAffected, err := databases.GetSingleUser(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponse("can't get single user"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponse("failed to get single user"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success get single user", user))
}

func CreateUserControllers(c echo.Context) error {
	var user models.User
	c.Bind(&user)

	newUser, err := databases.CreateNewUser(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponse("can't create new user"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success create new user", newUser))
}

func DeleteUserController(c echo.Context) error {
	userId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponse("invalid user id"))
	}

	loggedUserId := middlewares.ExtractToken(c)
	if loggedUserId != userId {
		return c.JSON(http.StatusOK, response.ErrorResponse("can't access other accounts"))
	}

	message, rowAffected, err := databases.DeleteUser(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponse("can't delete user data"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponse("failed to delete user data"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("successfully deleted user data", message))
}

func UpdateUserController(c echo.Context) error {
	userId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponse("invalid user id"))
	}

	loggedUserId := middlewares.ExtractToken(c)
	if loggedUserId != userId {
		return c.JSON(http.StatusOK, response.ErrorResponse("can't access other accounts"))
	}

	newUser := models.User{}
	c.Bind(&newUser)

	updatedUser, rowAffected, err := databases.UpdateUser(userId, newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponse("can't update user data"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponse("failed to update user data"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("successfully updated user data", updatedUser))
}

func LoginUserController(c echo.Context) error {
	user := models.User{}
	c.Bind(&user)

	loggedUser, err := databases.LoginUsers(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponse("failed to login"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponse("success to login", loggedUser))
}
