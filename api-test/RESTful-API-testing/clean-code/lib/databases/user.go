package databases

import (
	"cleancode/config"
	"cleancode/middlewares"
	"cleancode/models"
)

func GetAllUsers() (interface{}, int, error) {
	users := []models.User{}
	userOutput := []models.OutputUser{}

	result := config.Db.Model(&users).Find(&userOutput)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if result.RowsAffected > 0 {
		return userOutput, 1, nil
	}

	return "user data not found", 0, nil
}

func GetSingleUser(userId int) (interface{}, int, error) {
	user := models.User{}
	userOutput := models.OutputUser{}

	result := config.Db.Model(&user).Find(&userOutput, userId)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if result.RowsAffected > 0 {
		return userOutput, 1, nil
	}

	return "user data not found", 0, nil
}

func CreateNewUser(user *models.User) (interface{}, error) {
	result := config.Db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	userOutput := models.OutputUser{}
	userOutput.Name = user.Name
	userOutput.Email = user.Email

	return userOutput, nil
}

func DeleteUser(userId int) (interface{}, int, error) {
	user := models.User{}
	result := config.Db.Delete(&user, userId)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if result.RowsAffected > 0 {
		return "deleted", 1, nil
	}
	return "user data not found", 0, nil
}

func UpdateUser(userId int, newUser models.User) (interface{}, int, error) {
	user := models.User{}
	findResult := config.Db.Find(&user, userId)
	if findResult.Error != nil {
		return nil, 0, findResult.Error
	}

	if findResult.RowsAffected > 0 {
		updatedResult := config.Db.Model(&user).Updates(newUser)
		if updatedResult.Error != nil {
			return nil, 0, updatedResult.Error
		}

		userOutput := models.OutputUser{}
		userOutput.Name = user.Name
		userOutput.Email = user.Email

		return userOutput, 1, nil
	}

	return "user data not found", 0, nil
}

func LoginUsers(user *models.User) (interface{}, error) {
	result := config.Db.Where("email = ? AND password = ?", user.Email, user.Password).First(user)
	if result.Error != nil {
		return nil, result.Error
	}

	var err error
	user.Token, err = middlewares.CreateToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	saveToken := config.Db.Save(&user)
	if saveToken.Error != nil {
		return nil, saveToken.Error
	}

	return user, nil
}
