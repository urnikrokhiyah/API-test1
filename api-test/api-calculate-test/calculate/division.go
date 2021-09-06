package calculate

func DivisionNumber(number1, number2 int) interface{} {
	if number2 == 0 {
		return "division by zero"
	}
	return number1 / number2
}
