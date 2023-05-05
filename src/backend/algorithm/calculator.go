package algorithm

import (
	"errors"
	"github.com/niemeyer/golang/src/pkg/container/vector"
	"math"
	"strconv"
)

func calculator(expression string) string {
	ret, err := calculatePostfix(infixToPostFix(expression))
	if err != nil {
		return err.Error()
	}
	return strconv.FormatFloat(ret, 'f', 3, 32)
}

func infixToPostFix(expression string) string {
	operator := vector.StringVector{}
	ret := ""
	i := 0
	for i < len(expression) {
		s := string(expression[i])
	L1:
		if s == "+" || s == "-" || s == "*" || s == "/" || s == "^" || s == "(" || s == ")" {
			if operator.Len() == 0 {
				operator.Push(s)
			} else {
				if s == "(" {
					operator.Push(s)
				} else if s == ")" {
					for operator.Len() > 0 && operator.At(operator.Len()-1) != "(" {
						ret += operator.At(operator.Len() - 1)
						operator.Pop()
					}
					if operator.Len() > 0 && operator.At(operator.Len()-1) == "(" {
						operator.Pop()
					}
				} else if getPrecedence(operator.At(operator.Len()-1)) < getPrecedence(s) {
					operator.Push(s)
				} else if getPrecedence(operator.At(operator.Len()-1)) == getPrecedence(s) {
					ret += operator.At(operator.Len() - 1)
					operator.Pop()
					operator.Push(s)
				} else if getPrecedence(operator.At(operator.Len()-1)) > getPrecedence(s) {
					ret += operator.At(operator.Len() - 1)
					operator.Pop()
					goto L1
				}
			}
		} else {
			if s != "(" && s != ")" {
				ret += s
			}
		}
		i++
	}
	for operator.Len() > 0 {
		ret += operator.At(operator.Len() - 1)
		operator.Pop()
	}
	return ret
}

func getPrecedence(ops string) int {
	if ops == "-" || ops == "+" {
		return 1
	} else if ops == "*" || ops == "/" {
		return 2
	} else if ops == "^" {
		return 3
	} else {
		return 0
	}
}

func calculatePostfix(expression string) (float64, error) {
	ret := vector.Vector{}
	temp1 := 0.
	temp2 := 0.
	for i := 0; i < len(expression); i++ {
		s := string(expression[i])
		if s == "+" || s == "-" || s == "*" || s == "/" || s == "^" {
			if ret.Len() < 2 {
				return 0, errors.New("unexpected expression")
			}
			temp2 = ret.At(ret.Len() - 1).(float64)
			ret.Pop()
			temp1 = ret.At(ret.Len() - 1).(float64)
			ret.Pop()
			if s == "+" {
				res := temp1 + temp2
				ret.Push(res)
			} else if s == "-" {
				res := temp1 - temp2
				ret.Push(res)
			} else if s == "*" {
				res := temp1 * temp2
				ret.Push(res)
			} else if s == "/" {
				if temp2 == 0 {
					//panic(errors.New("zero division error"))
					return 0, errors.New("zero division error")
				}
				res := temp1 / temp2
				ret.Push(res)
			} else if s == "^" {
				res := math.Pow(temp1, temp2)
				ret.Push(res)
			} else {
				//panic(errors.New("unexpected operand"))
				return 0, errors.New("unexpected operand")
			}
		} else {
			res, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic(err)
			}
			ret.Push(res)
		}
	}
	return ret.At(0).(float64), nil
}
