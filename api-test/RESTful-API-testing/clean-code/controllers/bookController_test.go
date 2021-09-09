package controllers

import (
	"bytes"
	"cleancode/config"
	"cleancode/constants"
	"cleancode/middlewares"
	"cleancode/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func InitEchoTestAPIBook() *echo.Echo {
	config.InitDbTest()
	e := echo.New()
	return e
}

func InsertDataBookForGetBooks() error {
	book := models.Book{
		Title:        "chemistry",
		Author:       "urnik",
		Published_at: "2021",
	}

	err := config.Db.Save(&book).Error
	if err != nil {
		return err
	}
	return nil
}

func TestGetAllBooksController(t *testing.T) {
	type Expected struct {
		name         string
		path         string
		expectedCode int
	}

	testCases := Expected{
		name:         "success get all books",
		path:         "/books",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPIBook()
	InsertDataBookForGetBooks()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	record := httptest.NewRecorder()
	c := e.NewContext(request, record)

	c.SetPath(testCases.path)

	type BookResponse struct {
		Message string
		Data    models.Book
	}

	if assert.NoError(t, GetAllBooksController(c)) {
		body := record.Body.String()
		var book BookResponse
		err := json.Unmarshal([]byte(body), &book)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectedCode, record.Code)
		assert.Equal(t, "success", book.Message)
	}

}

func TestGetAllBooksControllerNilFailed(t *testing.T) {
	type Expected struct {
		name         string
		path         string
		expectedCode int
	}

	testCases := Expected{
		name:         "success get all books",
		path:         "/books",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPIBook()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	record := httptest.NewRecorder()
	c := e.NewContext(request, record)

	c.SetPath(testCases.path)

	type BookResponse struct {
		Message string
		Data    models.Book
	}

	if assert.NoError(t, GetAllBooksController(c)) {
		body := record.Body.String()
		var book BookResponse
		err := json.Unmarshal([]byte(body), &book)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectedCode, record.Code)
		assert.Equal(t, "failed", book.Message)
	}

}

func TestGetAllBooksControllerFailed(t *testing.T) {
	type Expected struct {
		name         string
		path         string
		expectedCode int
	}

	testCases := Expected{
		name:         "success get all books",
		path:         "/books",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPIBook()
	config.Db.Migrator().DropTable(&models.Book{})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	record := httptest.NewRecorder()
	c := e.NewContext(request, record)

	c.SetPath(testCases.path)

	type BookResponse struct {
		Message string
		Data    models.Book
	}

	if assert.NoError(t, GetAllBooksController(c)) {
		body := record.Body.String()
		var book BookResponse
		err := json.Unmarshal([]byte(body), &book)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectedCode, record.Code)
		assert.Equal(t, "failed", book.Message)
	}

}

func TestCreateBookController(t *testing.T) {
	type Expected struct {
		name         string
		expectedCode int
	}
	testCases := Expected{
		name:         "success create new book",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPIBook()

	book := models.Book{
		Title:        "math",
		Author:       "urnik",
		Published_at: "2013",
	}

	body, err := json.Marshal(book)
	if err != nil {
		t.Error(err, "error")
	}

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	type BookResponse struct {
		Message string
		Data    models.Book
	}

	if assert.NoError(t, CreateBookControllers(c)) {
		body := rec.Body.String()
		var books BookResponse

		err := json.Unmarshal([]byte(body), &books)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectedCode, rec.Code)
		assert.Equal(t, "math", books.Data.Title)
		assert.Equal(t, "urnik", books.Data.Author)
		assert.Equal(t, "success", books.Message)
	}

}

func TestCreateBookControllerFailed(t *testing.T) {
	type Expected struct {
		name         string
		expectedCode int
	}
	testCases := Expected{
		name:         "failed create new book",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPIBook()
	config.Db.Migrator().DropTable(&models.Book{})

	req := httptest.NewRequest(http.MethodPost, "/books", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	type BookResponse struct {
		Message string
		Data    models.Book
	}

	if assert.NoError(t, CreateBookControllers(c)) {
		body := rec.Body.String()
		var books BookResponse

		err := json.Unmarshal([]byte(body), &books)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectedCode, rec.Code)
		assert.Equal(t, "failed", books.Message)
	}

}

func TestGetSingleBookControllerSuccess(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "success get single book",
		id:           "1",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPIBook()
	InsertDataBookForGetBooks()

	req := httptest.NewRequest(http.MethodGet, "/users/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	type BookResponse struct {
		Message string
		Data    models.Book
	}

	if assert.NoError(t, GetSingleBookController(c)) {
		body := rec.Body.String()
		var book BookResponse
		err := json.Unmarshal([]byte(body), &book)
		if err != nil {
			assert.Error(t, err, "error")
		}
		assert.Equal(t, testCase.expectedCode, rec.Code)
		assert.Equal(t, "success", book.Message)
	}
}

func TestGetSingleBookControllerFailed(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "failed to get single book",
		id:           "1",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPIBook()
	config.Db.Migrator().DropTable(&models.Book{})

	req := httptest.NewRequest(http.MethodGet, "/users/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	type BookResponse struct {
		Message string
		Data    models.Book
	}

	if assert.NoError(t, GetSingleBookController(c)) {
		body := rec.Body.String()
		var book BookResponse
		err := json.Unmarshal([]byte(body), &book)
		if err != nil {
			assert.Error(t, err, "error")
		}
		assert.Equal(t, testCase.expectedCode, rec.Code)
		assert.Equal(t, "failed", book.Message)
	}
}

func TestGetSingleBookControllerInvalidId(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "invalid book id",
		id:           "s",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPIBook()
	InsertDataBookForGetBooks()

	req := httptest.NewRequest(http.MethodGet, "/users/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	type BookResponse struct {
		Message string
		Data    models.Book
	}

	if assert.NoError(t, GetSingleBookController(c)) {
		body := rec.Body.String()
		var book BookResponse
		err := json.Unmarshal([]byte(body), &book)
		if err != nil {
			assert.Error(t, err, "error")
		}
		assert.Equal(t, testCase.expectedCode, rec.Code)
		assert.Equal(t, "invalid book id", book.Message)
	}
}

func TestGetSingleBookControllerNilFailed(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "failed to get single book",
		id:           "2",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPIBook()
	InsertDataBookForGetBooks()

	req := httptest.NewRequest(http.MethodGet, "/users/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	type BookResponse struct {
		Message string
		Data    models.Book
	}

	if assert.NoError(t, GetSingleBookController(c)) {
		body := rec.Body.String()
		var book BookResponse
		err := json.Unmarshal([]byte(body), &book)
		if err != nil {
			assert.Error(t, err, "error")
		}
		assert.Equal(t, testCase.expectedCode, rec.Code)
		assert.Equal(t, "failed", book.Message)
	}
}

func TestDeleteBookControllerSuccess(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "success delete book",
		id:           "1",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPIBook()
	InsertDataBookForGetBooks()
	InsertDataUserForGetUsers()

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}
	user := models.User{}
	result := config.Db.Where("Email = ? AND Password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err1 := middlewares.CreateToken(int(user.ID))
	if err1 != nil {
		t.Error(err1)
	}

	req := httptest.NewRequest(http.MethodDelete, "/jwt/books/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteBookTesting())(c)

	type BookResponse struct {
		Message string
	}

	body := rec.Body.String()
	var book BookResponse
	err2 := json.Unmarshal([]byte(body), &book)
	if err2 != nil {
		assert.Error(t, err2, "error")
	}

	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "success", book.Message)
}

func TestDeleteBookControllerNilFailed(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "failed delete book",
		id:           "22",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPIBook()
	InsertDataBookForGetBooks()
	InsertDataUserForGetUsers()

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}
	user := models.User{}
	result := config.Db.Where("Email = ? AND Password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err1 := middlewares.CreateToken(int(user.ID))
	if err1 != nil {
		t.Error(err1)
	}

	req := httptest.NewRequest(http.MethodDelete, "/jwt/books/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteBookTesting())(c)

	type BookResponse struct {
		Message string
	}

	body := rec.Body.String()
	var book BookResponse
	err2 := json.Unmarshal([]byte(body), &book)
	if err2 != nil {
		assert.Error(t, err2, "error")
	}

	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "failed", book.Message)
}

func TestDeleteBookControllerInvalidId(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "failed delete book",
		id:           "s",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPIBook()
	InsertDataBookForGetBooks()
	InsertDataUserForGetUsers()

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}
	user := models.User{}
	result := config.Db.Where("Email = ? AND Password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err1 := middlewares.CreateToken(int(user.ID))
	if err1 != nil {
		t.Error(err1)
	}

	req := httptest.NewRequest(http.MethodDelete, "/jwt/books/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteBookTesting())(c)

	type BookResponse struct {
		Message string
	}

	body := rec.Body.String()
	var book BookResponse
	err2 := json.Unmarshal([]byte(body), &book)
	if err2 != nil {
		assert.Error(t, err2, "error")
	}

	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "invalid book id", book.Message)
}

func TestDeleteBookControllerFailed(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "failed delete book",
		id:           "1",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPIBook()
	config.Db.Migrator().DropTable(&models.Book{})
	InsertDataUserForGetUsers()

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}
	user := models.User{}
	result := config.Db.Where("Email = ? AND Password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err1 := middlewares.CreateToken(int(user.ID))
	if err1 != nil {
		t.Error(err1)
	}

	req := httptest.NewRequest(http.MethodDelete, "/jwt/books/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteBookTesting())(c)

	type BookResponse struct {
		Message string
	}

	body := rec.Body.String()
	var book BookResponse
	err2 := json.Unmarshal([]byte(body), &book)
	if err2 != nil {
		assert.Error(t, err2, "error")
	}

	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "failed", book.Message)
}

func TestUpdateBookControllerSuccess(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "success update book",
		id:           "1",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPIBook()
	InsertDataBookForGetBooks()
	InsertDataUserForGetUsers()

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}
	user := models.User{}
	result := config.Db.Where("Email = ? AND Password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err1 := middlewares.CreateToken(int(user.ID))
	if err1 != nil {
		t.Error(err1)
	}

	Newbook := models.Book{
		Title:        "mathematics",
		Author:       "lukman",
		Published_at: "2013",
	}

	body, err := json.Marshal(Newbook)
	if err != nil {
		t.Error(err, "error")
	}

	req := httptest.NewRequest(http.MethodPut, "/jwt/books/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateBookTesting())(c)

	type BookResponse struct {
		Message string
		Data    models.Book
	}

	recordBody := rec.Body.String()
	var book BookResponse
	err2 := json.Unmarshal([]byte(recordBody), &book)
	if err2 != nil {
		assert.Error(t, err2, "error")
	}

	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "success", book.Message)
	assert.Equal(t, Newbook.Title, book.Data.Title)

}

func TestUpdateBookControllerNilFailed(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "failed to update book",
		id:           "22",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPIBook()
	InsertDataBookForGetBooks()
	InsertDataUserForGetUsers()

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}
	user := models.User{}
	result := config.Db.Where("Email = ? AND Password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err1 := middlewares.CreateToken(int(user.ID))
	if err1 != nil {
		t.Error(err1)
	}

	req := httptest.NewRequest(http.MethodPut, "/jwt/books/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateBookTesting())(c)

	type BookResponse struct {
		Message string
	}

	recordBody := rec.Body.String()
	var book BookResponse
	err2 := json.Unmarshal([]byte(recordBody), &book)
	if err2 != nil {
		assert.Error(t, err2, "error")
	}

	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "failed", book.Message)
}

func TestUpdateBookControllerInvalidId(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "invalid book id",
		id:           "s",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPIBook()
	InsertDataBookForGetBooks()
	InsertDataUserForGetUsers()

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}
	user := models.User{}
	result := config.Db.Where("Email = ? AND Password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err1 := middlewares.CreateToken(int(user.ID))
	if err1 != nil {
		t.Error(err1)
	}

	req := httptest.NewRequest(http.MethodPut, "/jwt/books/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateBookTesting())(c)

	type BookResponse struct {
		Message string
	}

	recordBody := rec.Body.String()
	var book BookResponse
	err2 := json.Unmarshal([]byte(recordBody), &book)
	if err2 != nil {
		assert.Error(t, err2, "error")
	}

	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "invalid book id", book.Message)
}

func TestUpdateBookControllerFailed(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "invalid book id",
		id:           "1",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPIBook()
	config.Db.Migrator().DropTable(&models.Book{})
	InsertDataUserForGetUsers()

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}
	user := models.User{}
	result := config.Db.Where("Email = ? AND Password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err1 := middlewares.CreateToken(int(user.ID))
	if err1 != nil {
		t.Error(err1)
	}

	req := httptest.NewRequest(http.MethodPut, "/jwt/books/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateBookTesting())(c)

	type BookResponse struct {
		Message string
	}

	recordBody := rec.Body.String()
	var book BookResponse
	err2 := json.Unmarshal([]byte(recordBody), &book)
	if err2 != nil {
		assert.Error(t, err2, "error")
	}

	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "failed", book.Message)
}
