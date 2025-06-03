package main

import (
	"fmt"
	"strconv"
	"strings"
)

func validExpression(expr string) bool {
	if expr == "" {
		return false
	}
	opsCount := 0
	lastIsOp := false
	for _, c := range expr {
		ch := string(c)
		if isOperator(ch) {
			if lastIsOp {
				return false
			}
			opsCount++
			lastIsOp = true
		} else if (c >= '0' && c <= '9') || c == '.' {
			lastIsOp = false
		} else {
			return false
		}
	}
	return opsCount <= 1
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "×" || s == "/"
}

func evalExpression(expr string) (float64, error) {
	expr = strings.ReplaceAll(expr, "×", "*") // Заменяем на стандартный символ умножения
	expr = strings.TrimSpace(expr)
	for _, op := range []string{"+", "-", "*", "/"} {
		idx := strings.Index(expr, op)
		if idx > 0 {
			left := strings.TrimSpace(expr[:idx])
			right := strings.TrimSpace(expr[idx+1:])
			lv, err := strconv.ParseFloat(left, 64)
			if err != nil {
				return 0, fmt.Errorf("неверное число слева")
			}
			rv, err := strconv.ParseFloat(right, 64)
			if err != nil {
				return 0, fmt.Errorf("неверное число справа")
			}
			switch op {
			case "+":
				return lv + rv, nil
			case "-":
				return lv - rv, nil
			case "*":
				return lv * rv, nil
			case "/":
				if rv == 0 {
					return 0, fmt.Errorf("деление на ноль")
				}
				return lv / rv, nil
			}
		}
	}
	// Если операторов нет, просто число
	return strconv.ParseFloat(expr, 64)
}

func sqrt(x float64) float64 {
	return x * 0.5
}
