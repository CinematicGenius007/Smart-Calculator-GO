package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Stack []string
type StackInt []int

func (st *Stack) IsEmpty() bool {
	return len(*st) == 0
}

func (st *StackInt) IsEmpty() bool {
	return len(*st) == 0
}

func (st *Stack) Push(str string) {
	*st = append(*st, str)
}

func (st *StackInt) Push(number int) {
	*st = append(*st, number)
}

func (st *Stack) Pop() bool {
	if st.IsEmpty() {
		return false
	} else {
		index := len(*st) - 1
		*st = (*st)[:index]
		return true
	}
}

func (st *StackInt) Pop() bool {
	if st.IsEmpty() {
		return false
	} else {
		index := len(*st) - 1
		*st = (*st)[:index]
		return true
	}
}

func (st *Stack) Top() string {
	if st.IsEmpty() {
		return ""
	} else {
		index := len(*st) - 1
		element := (*st)[index]
		return element
	}
}

func (st *StackInt) Top() int {
	if st.IsEmpty() {
		return -1
	} else {
		index := len(*st) - 1
		element := (*st)[index]
		return element
	}
}

func IsLetter(letter uint8) bool {
	return (letter >= 'a' && letter <= 'z') || (letter >= 'A' && letter <= 'Z') || (letter >= '0' && letter <= '9')
}

func IsOperator(letter uint8) bool {
	return letter == '^' || letter == '*' || letter == '/' || letter == '+' || letter == '-'
}

func ProcessOperator(operator string) (string, bool) {
	if strings.ContainsAny(operator, "*") {
		return "*", len(operator) == 1
	} else if strings.ContainsAny(operator, "^") {
		return "^", len(operator) == 1
	} else if strings.ContainsAny(operator, "/") {
		return "/", len(operator) == 1
	} else if strings.ContainsAny(operator, "+") || strings.ContainsAny(operator, "-") {
		count := 0
		currentChar := operator[0]
		ans := 1
		for i := range operator {
			if currentChar != operator[i] {
				if currentChar == '-' && count%2 != 0 {
					ans *= -1
				}
				count = 0
				currentChar = operator[i]
			}
			count++
		}
		if currentChar == '-' && count%2 != 0 {
			ans *= -1
		}
		if ans == -1 {
			return "-", true
		}
		return "+", true
	}

	return "", false
}

func precedence(operand string) int {
	if operand == "^" {
		return 3
	} else if (operand == "/") || (operand == "*") {
		return 2
	} else if (operand == "+") || (operand == "-") {
		return 1
	} else {
		return -1
	}
}

func processInput(input string) (string, bool) {
	postfix := ""
	letters := 0
	operators := 0
	var stack Stack
	for i := 0; i < len(input); {
		word := ""
		if IsLetter(input[i]) {
			word += string(input[i])
			for i+1 < len(input) && IsLetter(input[i+1]) {
				i++
				word += string(input[i])
			}
			postfix += word + " "
			letters++
		} else if input[i] == '(' {
			stack.Push(string(input[i]))
		} else if input[i] == ')' {
			for !stack.IsEmpty() && stack.Top() != "(" {
				postfix += stack.Top() + " "
				stack.Pop()
			}
			if stack.IsEmpty() {
				return "", true
			}
			stack.Pop()
		} else if IsOperator(input[i]) {
			operator := string(input[i])
			for i+1 < len(input) && IsOperator(input[i+1]) {
				i++
				operator += string(input[i])
			}
			if temp, ok := ProcessOperator(operator); ok {
				operator = temp
				if (operator == "-" || operator == "+") &&
					i+1 < len(input) && IsLetter(input[i+1]) &&
					((!stack.IsEmpty() &&
						len(strings.Split(strings.TrimSpace(postfix), " "))%2 != 0 &&
						operators == letters) ||
						(stack.IsEmpty() && len(postfix) == 0)) {
					for i+1 < len(input) && IsLetter(input[i+1]) {
						i++
						operator += string(input[i])
					}
					postfix += operator + " "
					i++
					continue
				}
			} else {
				return "", true
			}
			for !stack.IsEmpty() && precedence(operator) <= precedence(stack.Top()) {
				postfix = postfix + stack.Top() + " "
				stack.Pop()
			}
			stack.Push(operator)
			operators++
		}

		i++
	}

	for !stack.IsEmpty() {
		if stack.Top() == "(" || stack.Top() == ")" {
			return "", true
		}
		postfix = postfix + stack.Top() + " "
		stack.Pop()
	}

	return postfix, false
}

func checkVariableName(name string) bool {
	for i := 0; i < len(name); i++ {
		if !(name[i] >= 'a' && name[i] <= 'z' || name[i] >= 'A' && name[i] <= 'Z') {
			return false
		}
	}
	return true
}

func EvaluatePostfixExpression(expression []string, variables map[string]int) (int, bool) {
	var stack StackInt

	numbers := regexp.MustCompile("\\d+")
	letters := regexp.MustCompile("[a-zA-Z]+")

	for _, i := range expression {
		if numbers.MatchString(i) {
			number, _ := strconv.Atoi(i)
			stack.Push(number)
		} else if letters.MatchString(i) {
			if number, ok := variables[i]; ok {
				stack.Push(number)
			} else {
				return 0, false
			}
		} else {
			num1 := stack.Top()
			stack.Pop()
			num2 := stack.Top()
			stack.Pop()
			switch {
			case i == "^":
				stack.Push(int(math.Pow(float64(num2), float64(num1))))
			case i == "*":
				stack.Push(num2 * num1)
			case i == "/":
				stack.Push(num2 / num1)
			case i == "+":
				stack.Push(num2 + num1)
			case i == "-":
				stack.Push(num2 - num1)
			}
		}
	}

	return stack.Top(), true
}

func main() {
	space := regexp.MustCompile("\\s+")
	spaceExpression := regexp.MustCompile("\\s*=\\s*")
	command := regexp.MustCompile("/\\w+")

	variables := make(map[string]int)

	var numString string
	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		numString = strings.TrimSpace(scanner.Text())
	}

	for numString != "/exit" {
		if command.MatchString(numString) {
			if numString == "/help" {
				fmt.Println("The program calculates expressions containing addition and subtraction of numbers")
			} else {
				fmt.Println("Unknown command")
			}
		} else if strings.ContainsAny(numString, "=") {
			values := spaceExpression.Split(numString, -1)
			if checkVariableName(values[0]) {
				num, err := strconv.Atoi(values[1])
				if err != nil {
					if checkVariableName(values[1]) {
						if val, ok := variables[values[1]]; ok {
							variables[values[0]] = val
						} else {
							fmt.Println("Unknown variable")
						}
					} else {
						fmt.Println("Invalid assignment")
					}
				} else {
					variables[values[0]] = num
				}
			} else {
				fmt.Println("Invalid identifier")
			}
		} else {
			if expression, err := processInput(numString); err {
				fmt.Println("Invalid expression")
			} else {
				if len(strings.TrimSpace(expression)) != 0 {
					if ans, ok := EvaluatePostfixExpression(space.Split(strings.TrimSpace(expression), -1), variables); ok {
						fmt.Println(ans)
					} else {
						fmt.Println("Unknown variable")
					}
				}
			}
		}

		if scanner.Scan() {
			numString = strings.TrimSpace(scanner.Text())
		}
	}
	fmt.Println("Bye!")

}
