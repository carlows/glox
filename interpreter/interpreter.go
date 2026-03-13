package interpreter

import (
	"fmt"
	"glox/expr"
	"glox/scanner"
	"strings"
)

type Interpreter struct{
	errFn func(token scanner.Token, message string)
}

func NewInterpreter(errFn func(token scanner.Token, message string)) *Interpreter {
	return &Interpreter{
		errFn: errFn,
	}
}

func (i *Interpreter) Interpret(expr expr.Expr) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*RuntimeError); !ok {
				panic(r) // re-panic if it's not a runtime error
			}
			// runtime error already reported via callback, just return
		}
	}()

	value := i.evaluate(expr)
	fmt.Println(i.stringify(value))
}

func (i *Interpreter) stringify(value any) string {
	if value == nil {
		return "nil"
	}

	obj, ok := value.(float64)
	if ok {
		formatted := fmt.Sprintf("%.2f", obj)
		if strings.HasSuffix(formatted, ".00") {
			return formatted[:len(formatted)-3]
		}
		return formatted
	}

	return fmt.Sprint(value)
}

func (i *Interpreter) VisitLiteral(expr *expr.Literal) any {
	return expr.Value
}

func (i *Interpreter) VisitGrouping(expr *expr.Grouping) any {
	return i.evaluate(expr.Expr)
}

func (i *Interpreter) VisitUnary(expr *expr.Unary) any {
	right := i.evaluate(expr.Expr)

	switch expr.Op.Type {
	case scanner.Bang:
		return !i.isTruthy(right)
	case scanner.Minus:
		i.checkNumberOperand(expr.Op, right)
		return -right.(float64)
	}

	return right
}

func (i *Interpreter) VisitBinary(expr *expr.Binary) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Op.Type {
	case scanner.Minus:
		i.checkNumberOperands(expr.Op, left, right)
		return left.(float64) - right.(float64)
	case scanner.Slash:
		i.checkNumberOperands(expr.Op, left, right)
		return left.(float64) / right.(float64)
	case scanner.Star:
		i.checkNumberOperands(expr.Op, left, right)
		return left.(float64) * right.(float64)
	case scanner.Plus:
		l, lok := left.(float64)
		r, rok := right.(float64)
		if lok && rok {
			return l + r
		}
		ls, slok := left.(string)
		rs, srok := right.(string)
		if slok && srok {
			return ls + rs
		}

		i.error(&RuntimeError{
			Token:   expr.Op,
			Message: fmt.Sprintf("Cannot add %v and %v", left, right),
		})
	case scanner.Greater:
		i.checkNumberOperands(expr.Op, left, right)
		return left.(float64) > right.(float64)
	case scanner.GreaterEqual:
		i.checkNumberOperands(expr.Op, left, right)
		return left.(float64) >= right.(float64)
	case scanner.Less:
		i.checkNumberOperands(expr.Op, left, right)
		return left.(float64) < right.(float64)
	case scanner.LessEqual:
		i.checkNumberOperands(expr.Op, left, right)
		return left.(float64) <= right.(float64)
	case scanner.EqualEqual:
		return left == right
	case scanner.BangEqual:
		return left != right
	}

	return nil
}

func (i *Interpreter) checkNumberOperand(operator scanner.Token, operand any) {
	if _, ok := operand.(float64); !ok {
		i.error(&RuntimeError{
			Token:   operator,
			Message: fmt.Sprintf("Expected number operand, got %v", operand),
		})
	}
}

func (i *Interpreter) checkNumberOperands(operator scanner.Token, left any, right any) {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)

	if !leftOk || !rightOk {
		i.error(&RuntimeError{
			Token:   operator,
			Message: fmt.Sprintf("Expected number operands, got %v and %v", left, right),
		})
	}
}

func (i *Interpreter) evaluate(expr expr.Expr) any {
	return expr.Accept(i)
}

// Ruby-ish truthy check, nil and false are false, everything else is true
func (i *Interpreter) isTruthy(value any) bool {
	if value == nil {
		return false
	}
	boolean, ok := value.(bool)
	if ok {
		return boolean
	}

	return true
}

func (i *Interpreter) error(runtimeError *RuntimeError) {
	i.errFn(runtimeError.Token, runtimeError.Message)
	panic(runtimeError)
}

type RuntimeError struct {
	Token   scanner.Token
	Message string
}

func (r *RuntimeError) Error() string {
	return fmt.Sprintf("Error at line %d: %s", r.Token.Line, r.Message)
}
