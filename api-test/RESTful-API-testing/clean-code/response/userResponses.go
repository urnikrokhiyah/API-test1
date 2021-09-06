package response

func SuccessResponse(message string, user interface{}) map[string]interface{} {
	var response = map[string]interface{}{
		"message": message,
		"users":   user,
	}
	return response
}

func ErrorResponse(message string) map[string]interface{} {
	var response = map[string]interface{}{
		"message": message,
	}
	return response
}
