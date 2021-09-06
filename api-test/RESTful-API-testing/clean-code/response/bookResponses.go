package response

func SuccessResponseBook(message string, book interface{}) map[string]interface{} {
	var response = map[string]interface{}{
		"message": message,
		"books":   book,
	}
	return response
}

func ErrorResponseBook(message string) map[string]interface{} {
	var response = map[string]interface{}{
		"message": message,
	}
	return response
}
