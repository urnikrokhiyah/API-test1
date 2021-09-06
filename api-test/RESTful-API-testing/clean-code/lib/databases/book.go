package databases

import (
	"cleancode/config"
	"cleancode/models"
)

func GetAllBooks() (interface{}, int, error) {
	books := []models.Book{}
	outputBook := []models.OutputBook{}
	result := config.Db.Model(&books).Find(&outputBook)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if result.RowsAffected > 0 {
		return outputBook, 1, nil
	}

	return "Book not found", 0, nil
}

func GetSingleBook(bookId int) (interface{}, int, error) {
	book := models.Book{}
	bookOutput := models.OutputBook{}
	result := config.Db.Model(&book).Find(&bookOutput, bookId)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if result.RowsAffected > 0 {
		return bookOutput, 1, nil
	}

	return "Book not found", 0, nil
}

func CreateNewBook(book *models.Book) (interface{}, error) {
	result := config.Db.Create(&book)
	if result.Error != nil {
		return nil, result.Error
	}

	bookOutput := models.OutputBook{}
	bookOutput.Author = book.Author
	bookOutput.Title = book.Title
	bookOutput.Published_at = book.Published_at

	return bookOutput, nil
}

func DeleteBook(bookId int) (interface{}, int, error) {
	book := models.Book{}
	result := config.Db.Delete(&book, bookId)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if result.RowsAffected > 0 {
		return "deleted", 1, nil
	}
	return "Book not found", 0, nil
}

func UpdateBook(bookId int, newBook models.Book) (interface{}, int, error) {
	book := models.Book{}
	findResult := config.Db.Find(&book, bookId)
	if findResult.Error != nil {
		return nil, 0, findResult.Error
	}

	if findResult.RowsAffected > 0 {
		updatedResult := config.Db.Model(&book).Updates(newBook)
		if updatedResult.Error != nil {
			return nil, 0, updatedResult.Error
		}

		bookOutput := models.OutputBook{}
		bookOutput.Author = book.Author
		bookOutput.Title = book.Title
		bookOutput.Published_at = book.Published_at

		return bookOutput, 1, nil
	}

	return "Book not found", 0, nil
}
