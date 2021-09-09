package controllers

import (
	"cleancode/lib/databases"
	"cleancode/models"
	"cleancode/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllBooksController(c echo.Context) error {
	books, rowAffected, err := databases.GetAllBooks()
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBook("failed"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponseBook("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseBook("success", books))
}

func GetSingleBookController(c echo.Context) error {
	bookId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBook("invalid book id"))
	}

	book, rowAffected, err := databases.GetSingleBook(bookId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBook("failed"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponseBook("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseBook("success", book))
}

func CreateBookControllers(c echo.Context) error {
	var book models.Book
	c.Bind(&book)

	newBook, err := databases.CreateNewBook(&book)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBook("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseBook("success", newBook))
}

func DeleteBookController(c echo.Context) error {
	bookId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBook("invalid book id"))
	}

	message, rowAffected, err := databases.DeleteBook(bookId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBook("failed"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponseBook("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseBook("success", message))
}

func UpdateBookController(c echo.Context) error {
	bookId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBook("invalid book id"))
	}

	newBook := models.Book{}
	c.Bind(&newBook)

	updatedBook, rowAffected, err := databases.UpdateBook(bookId, newBook)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponseBook("failed"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponseBook("failed"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseBook("success", updatedBook))
}

func DeleteBookTesting() echo.HandlerFunc {
	return DeleteBookController
}

func UpdateBookTesting() echo.HandlerFunc {
	return UpdateBookController
}
