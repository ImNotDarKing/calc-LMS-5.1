package calculator

import (
	"errors"
	"strings"
)

var (
	ErrInvalidExpression  = errors.New("invalid expression")
	ErrDivisionByZero     = errors.New("division by zero")
	ErrMismatchedBrackets = errors.New("mismatched brackets")
	ErrMultiDigitNumber   = errors.New("multi-digit numbers are not allowed")
	ErrEmptyExpression    = errors.New("empty expression")
	ErrUnexpectedOperator = errors.New("unexpected operator without operand")
	ErrOperatorAtEnd      = errors.New("expression cannot end with an operator")
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")

	if len(expression) == 0 {
		return 0, ErrEmptyExpression
	}

	var num []float64
	var operator []rune
	negativeFlag := false

	priority := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	applyOperator := func(a, b float64, op rune) (float64, error) {
		switch op {
		case '+':
			return a + b, nil
		case '-':
			return a - b, nil
		case '*':
			return a * b, nil
		case '/':
			if b == 0 {
				return 0, ErrDivisionByZero
			}
			return a / b, nil
		default:
			return 0, ErrInvalidExpression
		}
	}

	calculate := func() error {
		if len(operator) == 0 || len(num) < 2 {
			return nil
		}
		b := num[len(num)-1]
		a := num[len(num)-2]
		op := operator[len(operator)-1]
		num = num[:len(num)-2]
		operator = operator[:len(operator)-1]
		result, err := applyOperator(a, b, op)
		if err != nil {
			return err
		}
		num = append(num, result)
		return nil
	}

	isOperator := func(char rune) bool {
		return char == '+' || char == '-' || char == '*' || char == '/'
	}

	for i, char := range expression {
		switch char {
		case '(':
			operator = append(operator, char)
		case ')':
			if i > 0 && (expression[i-1] == '(' || isOperator(rune(expression[i-1]))) {
				return 0, ErrInvalidExpression
			}
			for len(operator) > 0 && operator[len(operator)-1] != '(' {
				if err := calculate(); err != nil {
					return 0, err
				}
			}
			if len(operator) == 0 {
				return 0, ErrMismatchedBrackets
			}
			operator = operator[:len(operator)-1]
		case '+', '*', '/':
			if i == 0 || isOperator(rune(expression[i-1])) || expression[i-1] == '(' {
				return 0, ErrUnexpectedOperator
			}
			for len(operator) > 0 && operator[len(operator)-1] != '(' && priority[operator[len(operator)-1]] >= priority[char] {
				if err := calculate(); err != nil {
					return 0, err
				}
			}
			operator = append(operator, char)
			negativeFlag = false
		case '-':
			if i == 0 || expression[i-1] == '(' {
				negativeFlag = true
			} else {
				if isOperator(rune(expression[i-1])) || i == len(expression)-1 {
					return 0, ErrUnexpectedOperator
				}
				for len(operator) > 0 && operator[len(operator)-1] != '(' && priority[operator[len(operator)-1]] >= priority[char] {
					if err := calculate(); err != nil {
						return 0, err
					}
				}
				operator = append(operator, char)
				negativeFlag = false
			}
		default:
			if char >= '0' && char <= '9' {
				if i > 0 && expression[i-1] >= '0' && expression[i-1] <= '9' {
					return 0, ErrMultiDigitNumber
				}
				value := float64(char - '0')
				if negativeFlag {
					value = -value
					negativeFlag = false
				}
				num = append(num, value)
			} else {
				return 0, ErrInvalidExpression
			}
		}
	}

	if len(expression) > 0 && isOperator(rune(expression[len(expression)-1])) {
		return 0, ErrOperatorAtEnd
	}

	for len(operator) > 0 {
		if operator[len(operator)-1] == '(' {
			return 0, ErrMismatchedBrackets
		}
		if err := calculate(); err != nil {
			return 0, err
		}
	}

	if len(num) != 1 {
		return 0, ErrInvalidExpression
	}

	return num[0], nil
}
