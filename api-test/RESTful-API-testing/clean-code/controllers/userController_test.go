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

func InitEchoTestAPI() *echo.Echo {
	config.InitDbTest()
	e := echo.New()
	return e
}

func InsertDataUserForGetUsers() error {
	user := models.User{
		Name:     "Alta",
		Password: "123",
		Email:    "alta@gmail.com",
	}

	err := config.Db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func TestGetUserController(t *testing.T) {
	type Expected struct {
		name         string
		expectedCode int
	}

	testCases := Expected{
		name:         "get all user",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	type UserResponse struct {
		Message string
		Data    models.User
	}

	if assert.NoError(t, GetAllUsersController(c)) {
		body := rec.Body.String()

		var user UserResponse
		err := json.Unmarshal([]byte(body), &user)

		if err != nil {
			assert.Error(t, err, "error")
		}
		assert.Equal(t, testCases.expectedCode, rec.Code)
		assert.Equal(t, "success", user.Message)
	}

}

func TestGetUserControllerNilFailed(t *testing.T) {
	type Expected struct {
		name         string
		expectedCode int
	}

	testCases := Expected{
		name:         "failed get all users",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPI()

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	type UserResponse struct {
		Message string      `json:"message"`
		Data    models.User `json:"data"`
	}

	if assert.NoError(t, GetAllUsersController(c)) {
		body := rec.Body.String()
		var user UserResponse
		err := json.Unmarshal([]byte(body), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectedCode, rec.Code)
		assert.Equal(t, "failed", user.Message)
	}

}

func TestGetUserControllerFailed(t *testing.T) {
	type Expected struct {
		name         string
		expectedCode int
	}

	testCases := Expected{
		name:         "failed get all users",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	config.Db.Migrator().DropTable(&models.User{})

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	type UserResponse struct {
		Message string
		Data    models.User
	}

	GetAllUsersController(c)

	body := rec.Body.String()
	var user UserResponse
	err := json.Unmarshal([]byte(body), (&user))
	if err != nil {
		assert.Error(t, err, "error")
	}

	assert.Equal(t, testCases.expectedCode, rec.Code)
	assert.Equal(t, "failed", user.Message)

}

func TestCreateUserController(t *testing.T) {
	type Expected struct {
		name         string
		expectedCode int
	}

	testCases := Expected{
		name:         "succes create new user",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPI()

	user := models.User{
		Name:     "urnik",
		Email:    "urnik@gmail.com",
		Password: "urnik123",
	}

	body, err := json.Marshal(user)
	if err != nil {
		t.Error(t, err, "error")
	}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	type UserResponse struct {
		Message string
		Data    models.User
	}

	if assert.NoError(t, CreateUserControllers(c)) {
		body := rec.Body.String()
		var user UserResponse

		err := json.Unmarshal([]byte(body), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectedCode, rec.Code)
		assert.Equal(t, "urnik", user.Data.Name)
		assert.Equal(t, "urnik@gmail.com", user.Data.Email)
		assert.Equal(t, "success", user.Message)
	}

}

func TestCreateUserControllerFailed(t *testing.T) {
	type Expected struct {
		name         string
		expectedCode int
	}

	testCases := Expected{
		name:         "failed create new user",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	config.Db.Migrator().DropTable(&models.User{})

	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	type UserResponse struct {
		Message string
		Data    models.User
	}

	if assert.NoError(t, CreateUserControllers(c)) {
		body := rec.Body.String()
		var user UserResponse

		err := json.Unmarshal([]byte(body), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectedCode, rec.Code)
		assert.Equal(t, "failed", user.Message)
	}

}

func TestGetSingleUserController(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCaseSuccess := Expected{
		name:         "success get single user",
		id:           "1",
		expectedCode: http.StatusOK,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/jwt/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCaseSuccess.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserDetailControllersTesting())(c)

	type UserResponse struct {
		Message string
		Data    models.User
	}

	body := rec.Body.String()
	var users UserResponse

	err1 := json.Unmarshal([]byte(body), &users)
	if err != nil {
		assert.Error(t, err1, "error")
	}
	assert.Equal(t, testCaseSuccess.expectedCode, rec.Code)
	assert.Equal(t, "success", users.Message)
	assert.Equal(t, "Alta", users.Data.Name)

}

func TestGetSingleUserControllerInvalidId(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "invalid user id",
		id:           "s",
		expectedCode: http.StatusBadRequest,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/jwt/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserDetailControllersTesting())(c)

	type UserResponse struct {
		Message string
	}

	body := rec.Body.String()
	var users UserResponse

	err1 := json.Unmarshal([]byte(body), &users)
	if err != nil {
		assert.Error(t, err1, "error")
	}
	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "invalid user id", users.Message)

}

func TestGetSingleUserControllerUnauthorized(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "failed get single user",
		id:           "22",
		expectedCode: http.StatusBadRequest,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/jwt/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserDetailControllersTesting())(c)

	type UserResponse struct {
		Message string
		Data    models.User
	}

	body := rec.Body.String()
	var users UserResponse

	err1 := json.Unmarshal([]byte(body), &users)
	if err != nil {
		assert.Error(t, err1, "error")
	}
	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "unauthorized", users.Message)

}

// func TestGetSingleUserControllerNilFailed(t *testing.T) {
// 	type Expected struct {
// 		name         string
// 		id           string
// 		expectedCode int
// 	}

// 	testCase := Expected{
// 		name:         "failed get single user",
// 		id:           "1",
// 		expectedCode: http.StatusOK,
// 	}

// 	dummyData := models.User{
// 		Email:    "alta@gmail.com",
// 		Password: "123",
// 	}

// 	user := models.User{}

// 	e := InitEchoTestAPI()
// 	InsertDataUserForGetUsers()

// 	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}

// 	token, err := middlewares.CreateToken(int(user.ID))
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	req := httptest.NewRequest(http.MethodGet, "/jwt/users/:id", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	c.SetParamNames("id")
// 	c.SetParamValues(testCase.id)

// 	middleware.JWT([]byte(constants.SECRET_JWT))(GetUserDetailControllersTesting())(c)

// 	type UserResponse struct {
// 		Message string
// 		Data    models.User
// 	}

// 	body := rec.Body.String()
// 	var users UserResponse

// 	err1 := json.Unmarshal([]byte(body), &users)
// 	if err != nil {
// 		assert.Error(t, err1, "error")
// 	}
// 	assert.Equal(t, testCase.expectedCode, rec.Code)
// 	assert.Equal(t, "failed", users.Message)
// }

func TestGetSingleUserControllerFailed(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "failed get single user",
		id:           "1",
		expectedCode: http.StatusBadRequest,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		assert.Error(t, result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	config.Db.Migrator().DropTable(&models.User{})

	req := httptest.NewRequest(http.MethodGet, "/jwt/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	type UserResponse struct {
		Message string
		Data    models.User
	}

	if assert.NoError(t, middleware.JWT([]byte(constants.SECRET_JWT))(GetUserDetailControllersTesting())(c)) {
		body := rec.Body.String()
		var users UserResponse

		err1 := json.Unmarshal([]byte(body), &users)
		if err != nil {
			assert.Error(t, err1, "error")
		}
		assert.Equal(t, testCase.expectedCode, rec.Code)
		assert.Equal(t, "failed", users.Message)
	}
}

func TestLoginUserControllerFailed(t *testing.T) {
	type Expected struct {
		name         string
		expectedCode int
	}

	testCases := Expected{
		name:         "failed to login",
		expectedCode: http.StatusBadRequest,
	}

	e := InitEchoTestAPI()
	config.Db.Migrator().DropTable(&models.User{})

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	type UserResponse struct {
		Message string
		Data    models.User
	}

	if assert.NoError(t, LoginUserController(c)) {
		body := rec.Body.String()
		var user UserResponse

		err := json.Unmarshal([]byte(body), &user)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectedCode, rec.Code)
		assert.Equal(t, "failed", user.Message)
	}

}

func TestLoginUserController(t *testing.T) {
	type Expected struct {
		name         string
		expectedCode int
	}

	testCases := Expected{
		name:         "succes to login ",
		expectedCode: http.StatusOK,
	}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	user := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	body, err := json.Marshal(user)
	if err != nil {
		t.Error(t, err, "error")
	}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	type UserResponse struct {
		Message string
		Data    models.User
	}

	if assert.NoError(t, LoginUserController(c)) {
		body := rec.Body.String()
		var users UserResponse

		err := json.Unmarshal([]byte(body), &users)
		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, testCases.expectedCode, rec.Code)
		assert.Equal(t, user.Email, users.Data.Email)
		assert.Equal(t, user.Password, users.Data.Password)
		assert.Equal(t, "success", users.Message)
	}

}

func TestUpdatedUserControllerInvalidId(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "invalid user id",
		id:           "s",
		expectedCode: http.StatusBadRequest,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/jwt/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(UpdatedDetailUserTesting())(c)

	type UserResponse struct {
		Message string
	}

	body := rec.Body.String()
	var users UserResponse

	err1 := json.Unmarshal([]byte(body), &users)
	if err != nil {
		assert.Error(t, err1, "error")
	}
	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "invalid user id", users.Message)

}

func TestUpdateUserControllerUnauthorized(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "invalid user id",
		id:           "2",
		expectedCode: http.StatusBadRequest,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPut, "/jwt/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(UpdatedDetailUserTesting())(c)

	type UserResponse struct {
		Message string
	}

	body := rec.Body.String()
	var users UserResponse

	err1 := json.Unmarshal([]byte(body), &users)
	if err != nil {
		assert.Error(t, err1, "error")
	}
	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "unauthorized", users.Message)

}

// func TestUpdateUserControllerFailed(t *testing.T) {
// 	type Expected struct {
// 		name         string
// 		id           string
// 		expectedCode int
// 	}

// 	testCase := Expected{
// 		name:         "failed to update user",
// 		id:           "1",
// 		expectedCode: http.StatusBadRequest,
// 	}

// 	dummyData := models.User{
// 		Email:    "alta@gmail.com",
// 		Password: "123",
// 	}

// 	user := models.User{}

// 	e := InitEchoTestAPI()
// 	InsertDataUserForGetUsers()

// 	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
// 	if result.Error != nil {
// 		t.Error(result.Error)
// 	}

// 	token, err := middlewares.CreateToken(int(user.ID))
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	config.Db.Migrator().DropTable(&models.User{})
// 	req := httptest.NewRequest(http.MethodPut, "/jwt/users/:id", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	c.SetParamNames("id")
// 	c.SetParamValues(testCase.id)

// 	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateBookTesting())(c)

// 	type UserResponse struct {
// 		Message string
// 	}

// 	body := rec.Body.String()
// 	var users UserResponse

// 	err1 := json.Unmarshal([]byte(body), &users)
// 	if err != nil {
// 		assert.Error(t, err1, "error")
// 	}
// 	assert.Equal(t, testCase.expectedCode, rec.Code)
// 	assert.Equal(t, "failed", users.Message)
// }

func TestUpdatedUserControllerSuccess(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "success to update user",
		id:           "1",
		expectedCode: http.StatusOK,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	NewUser := models.User{
		Name:     "urnik rokhiyah",
		Email:    "urnik456",
		Password: "669",
	}

	body, err1 := json.Marshal(NewUser)
	if err1 != nil {
		t.Error(err1, "error")
	}

	req := httptest.NewRequest(http.MethodPut, "/jwt/users/:id", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(UpdateBookTesting())(c)

	type UserResponse struct {
		Message string
		Data    models.User
	}

	recBody := rec.Body.String()
	var users UserResponse

	err2 := json.Unmarshal([]byte(recBody), &users)
	if err != nil {
		assert.Error(t, err2, "error")
	}

	fmt.Println(testCase.expectedCode, rec.Code)
	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, NewUser.Name, users.Data.Name)
	assert.Equal(t, "success", users.Message)
}

func TestDeleteUserControllerInvalidId(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "invalid user id",
		id:           "s",
		expectedCode: http.StatusBadRequest,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/jwt/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteDetailUserTesting())(c)

	type UserResponse struct {
		Message string
	}

	body := rec.Body.String()
	var users UserResponse

	err1 := json.Unmarshal([]byte(body), &users)
	if err != nil {
		assert.Error(t, err1, "error")
	}
	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "invalid user id", users.Message)

}

func TestDeleteUserControllerUnauthorized(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "invalid user id",
		id:           "2",
		expectedCode: http.StatusBadRequest,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/jwt/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteDetailUserTesting())(c)

	type UserResponse struct {
		Message string
	}

	body := rec.Body.String()
	var users UserResponse

	err1 := json.Unmarshal([]byte(body), &users)
	if err != nil {
		assert.Error(t, err1, "error")
	}
	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "unauthorized", users.Message)

}

func TestDeleteUserControllerSuccess(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "success to delete user",
		id:           "1",
		expectedCode: http.StatusOK,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/jwt/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteDetailUserTesting())(c)

	type UserResponse struct {
		Message string
	}

	body := rec.Body.String()
	var users UserResponse

	err1 := json.Unmarshal([]byte(body), &users)
	if err != nil {
		assert.Error(t, err1, "error")
	}
	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "success", users.Message)
}

func TestDeleteUserControllerFailed(t *testing.T) {
	type Expected struct {
		name         string
		id           string
		expectedCode int
	}

	testCase := Expected{
		name:         "failed to delete user",
		id:           "1",
		expectedCode: http.StatusBadRequest,
	}

	dummyData := models.User{
		Email:    "alta@gmail.com",
		Password: "123",
	}

	user := models.User{}

	e := InitEchoTestAPI()
	InsertDataUserForGetUsers()

	result := config.Db.Where("email = ? AND password = ?", dummyData.Email, dummyData.Password).First(&user)
	if result.Error != nil {
		t.Error(result.Error)
	}

	token, err := middlewares.CreateToken(int(user.ID))
	if err != nil {
		t.Error(err)
	}

	config.Db.Migrator().DropTable(&models.User{})
	req := httptest.NewRequest(http.MethodDelete, "/jwt/users/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(testCase.id)

	middleware.JWT([]byte(constants.SECRET_JWT))(DeleteDetailUserTesting())(c)

	type UserResponse struct {
		Message string
	}

	body := rec.Body.String()
	var users UserResponse

	err1 := json.Unmarshal([]byte(body), &users)
	if err != nil {
		assert.Error(t, err1, "error")
	}
	assert.Equal(t, testCase.expectedCode, rec.Code)
	assert.Equal(t, "failed", users.Message)
}
