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
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponseBook("can't get all books"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponseBook("failed to get all books"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseBook("success get all books", books))
}

func GetSingleBookController(c echo.Context) error {
	bookId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponseBook("invalid book id"))
	}

	book, rowAffected, err := databases.GetSingleBook(bookId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponseBook("can't get single book"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponseBook("failed to get single book"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseBook("success get single book", book))
}

func CreateBookControllers(c echo.Context) error {
	var book models.Book
	c.Bind(&book)

	newBook, err := databases.CreateNewBook(&book)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponseBook("can't create new book"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseBook("success create new book", newBook))
}

func DeleteBookController(c echo.Context) error {
	bookId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponseBook("invalid book id"))
	}

	message, rowAffected, err := databases.DeleteBook(bookId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponseBook("can't delete book data"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponseBook("failed to delete book data"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseBook("successfully deleted book data", message))
}

func UpdateBookController(c echo.Context) error {
	bookId, errorId := strconv.Atoi(c.Param("id"))
	if errorId != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponseBook("invalid book id"))
	}

	newBook := models.Book{}
	c.Bind(&newBook)

	updatedBook, rowAffected, err := databases.UpdateBook(bookId, newBook)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, response.ErrorResponseBook("can't update book data"))
	}

	if rowAffected == 0 {
		return c.JSON(http.StatusOK, response.ErrorResponseBook("failed to update book data"))
	}

	return c.JSON(http.StatusOK, response.SuccessResponseBook("successfully updated book data", updatedBook))
}
