package interpreter

import (
	"glox/expr"
	"glox/scanner"
	"testing"
)

func tok(t scanner.TokenType) scanner.Token {
	return scanner.Token{Type: t, Lexeme: t.String(), Line: 1}
}

func newInterp(t *testing.T) (*Interpreter, *[]string) {
	t.Helper()
	var errors []string
	interp := NewInterpreter(func(tk scanner.Token, message string) {
		errors = append(errors, message)
	})
	return interp, &errors
}

// evalSafe evaluates an expression, catching RuntimeError panics.
func evalSafe(i *Interpreter, e expr.Expr) (result any, runtimeErr *RuntimeError) {
	defer func() {
		if r := recover(); r != nil {
			if re, ok := r.(*RuntimeError); ok {
				runtimeErr = re
			} else {
				panic(r)
			}
		}
	}()
	result = i.evaluate(e)
	return
}

// --- Literals ---

func TestLiteralNumber(t *testing.T) {
	interp, _ := newInterp(t)
	result, err := evalSafe(interp, &expr.Literal{Value: 3.14})
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != 3.14 {
		t.Errorf("expected 3.14, got %v", result)
	}
}

func TestLiteralString(t *testing.T) {
	interp, _ := newInterp(t)
	result, err := evalSafe(interp, &expr.Literal{Value: "hello"})
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != "hello" {
		t.Errorf("expected \"hello\", got %v", result)
	}
}

func TestLiteralBool(t *testing.T) {
	interp, _ := newInterp(t)
	for _, val := range []bool{true, false} {
		result, err := evalSafe(interp, &expr.Literal{Value: val})
		if err != nil {
			t.Fatalf("unexpected runtime error: %v", err)
		}
		if result != val {
			t.Errorf("expected %v, got %v", val, result)
		}
	}
}

func TestLiteralNil(t *testing.T) {
	interp, _ := newInterp(t)
	result, err := evalSafe(interp, &expr.Literal{Value: nil})
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

// --- Grouping ---

func TestGrouping(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Grouping{Expr: &expr.Literal{Value: 42.0}}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != 42.0 {
		t.Errorf("expected 42.0, got %v", result)
	}
}

// --- Unary ---

func TestUnaryMinus(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Unary{Op: tok(scanner.Minus), Expr: &expr.Literal{Value: 5.0}}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != -5.0 {
		t.Errorf("expected -5.0, got %v", result)
	}
}

func TestUnaryMinusOnString(t *testing.T) {
	interp, errors := newInterp(t)
	e := &expr.Unary{Op: tok(scanner.Minus), Expr: &expr.Literal{Value: "hello"}}
	_, runtimeErr := evalSafe(interp, e)
	if runtimeErr == nil {
		t.Error("expected runtime error for unary minus on string, got none")
	}
	if len(*errors) == 0 {
		t.Error("expected error callback to be called")
	}
}

func TestUnaryBangFalse(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Unary{Op: tok(scanner.Bang), Expr: &expr.Literal{Value: false}}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != true {
		t.Errorf("expected true, got %v", result)
	}
}

func TestUnaryBangTrue(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Unary{Op: tok(scanner.Bang), Expr: &expr.Literal{Value: true}}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != false {
		t.Errorf("expected false, got %v", result)
	}
}

func TestUnaryBangNil(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Unary{Op: tok(scanner.Bang), Expr: &expr.Literal{Value: nil}}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	// nil is falsy, so !nil == true
	if result != true {
		t.Errorf("expected true (nil is falsy), got %v", result)
	}
}

func TestUnaryBangNumber(t *testing.T) {
	interp, _ := newInterp(t)
	// Numbers are truthy, so !42 == false
	e := &expr.Unary{Op: tok(scanner.Bang), Expr: &expr.Literal{Value: 42.0}}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != false {
		t.Errorf("expected false (number is truthy), got %v", result)
	}
}

// --- Binary arithmetic ---

func TestBinaryAdd(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: 1.0},
		Op:    tok(scanner.Plus),
		Right: &expr.Literal{Value: 2.0},
	}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != 3.0 {
		t.Errorf("expected 3.0, got %v", result)
	}
}

func TestBinarySubtract(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: 10.0},
		Op:    tok(scanner.Minus),
		Right: &expr.Literal{Value: 4.0},
	}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != 6.0 {
		t.Errorf("expected 6.0, got %v", result)
	}
}

func TestBinaryMultiply(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: 3.0},
		Op:    tok(scanner.Star),
		Right: &expr.Literal{Value: 4.0},
	}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != 12.0 {
		t.Errorf("expected 12.0, got %v", result)
	}
}

func TestBinaryDivide(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: 10.0},
		Op:    tok(scanner.Slash),
		Right: &expr.Literal{Value: 2.0},
	}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != 5.0 {
		t.Errorf("expected 5.0, got %v", result)
	}
}

// Dividing by zero produces +Inf in Go; this test documents current behavior.
// Ideally the interpreter should report a runtime error instead.
func TestBinaryDivideByZero(t *testing.T) {
	interp, errors := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: 1.0},
		Op:    tok(scanner.Slash),
		Right: &expr.Literal{Value: 0.0},
	}
	_, runtimeErr := evalSafe(interp, e)
	if runtimeErr == nil && len(*errors) == 0 {
		t.Log("NOTE: division by zero is not caught; result is +Inf (missing error handling)")
	}
}

// --- String concatenation ---

func TestBinaryStringConcat(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: "hello"},
		Op:    tok(scanner.Plus),
		Right: &expr.Literal{Value: " world"},
	}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("expected \"hello world\", got %v", result)
	}
}

func TestBinaryAddMixedTypes(t *testing.T) {
	interp, errors := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: "hello"},
		Op:    tok(scanner.Plus),
		Right: &expr.Literal{Value: 1.0},
	}
	_, runtimeErr := evalSafe(interp, e)
	if runtimeErr == nil {
		t.Error("expected runtime error for string + number, got none")
	}
	if len(*errors) == 0 {
		t.Error("expected error callback to be called")
	}
}

// --- Comparison operators ---

func TestBinaryGreater(t *testing.T) {
	interp, _ := newInterp(t)
	tests := []struct {
		left, right float64
		expected    bool
	}{
		{5.0, 3.0, true},
		{3.0, 5.0, false},
		{3.0, 3.0, false},
	}
	for _, tt := range tests {
		e := &expr.Binary{
			Left:  &expr.Literal{Value: tt.left},
			Op:    tok(scanner.Greater),
			Right: &expr.Literal{Value: tt.right},
		}
		result, err := evalSafe(interp, e)
		if err != nil {
			t.Fatalf("unexpected runtime error: %v", err)
		}
		if result != tt.expected {
			t.Errorf("%v > %v: expected %v, got %v", tt.left, tt.right, tt.expected, result)
		}
	}
}

func TestBinaryGreaterEqual(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: 3.0},
		Op:    tok(scanner.GreaterEqual),
		Right: &expr.Literal{Value: 3.0},
	}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != true {
		t.Errorf("expected true for 3 >= 3, got %v", result)
	}
}

func TestBinaryLess(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: 2.0},
		Op:    tok(scanner.Less),
		Right: &expr.Literal{Value: 5.0},
	}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != true {
		t.Errorf("expected true for 2 < 5, got %v", result)
	}
}

func TestBinaryLessEqual(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: 5.0},
		Op:    tok(scanner.LessEqual),
		Right: &expr.Literal{Value: 5.0},
	}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != true {
		t.Errorf("expected true for 5 <= 5, got %v", result)
	}
}

// --- Equality ---

func TestBinaryEqualEqual(t *testing.T) {
	interp, _ := newInterp(t)
	tests := []struct {
		left, right any
		expected    bool
	}{
		{1.0, 1.0, true},
		{1.0, 2.0, false},
		{"a", "a", true},
		{"a", "b", false},
		{nil, nil, true},
		{true, true, true},
		{true, false, false},
		// Different types should not be equal
		{1.0, "1", false},
		{nil, false, false},
	}
	for _, tt := range tests {
		e := &expr.Binary{
			Left:  &expr.Literal{Value: tt.left},
			Op:    tok(scanner.EqualEqual),
			Right: &expr.Literal{Value: tt.right},
		}
		result, err := evalSafe(interp, e)
		if err != nil {
			t.Fatalf("unexpected runtime error: %v", err)
		}
		if result != tt.expected {
			t.Errorf("%v == %v: expected %v, got %v", tt.left, tt.right, tt.expected, result)
		}
	}
}

func TestBinaryBangEqual(t *testing.T) {
	interp, _ := newInterp(t)
	e := &expr.Binary{
		Left:  &expr.Literal{Value: 1.0},
		Op:    tok(scanner.BangEqual),
		Right: &expr.Literal{Value: 2.0},
	}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != true {
		t.Errorf("expected true for 1 != 2, got %v", result)
	}
}

// --- Type errors for binary operators ---

func TestBinaryArithmeticTypeError(t *testing.T) {
	interp, errors := newInterp(t)
	ops := []scanner.TokenType{scanner.Minus, scanner.Star, scanner.Slash}
	for _, op := range ops {
		*errors = nil
		e := &expr.Binary{
			Left:  &expr.Literal{Value: "hello"},
			Op:    tok(op),
			Right: &expr.Literal{Value: 1.0},
		}
		_, runtimeErr := evalSafe(interp, e)
		if runtimeErr == nil {
			t.Errorf("operator %v: expected runtime error for string operand, got none", op)
		}
	}
}

func TestBinaryComparisonTypeError(t *testing.T) {
	interp, errors := newInterp(t)
	ops := []scanner.TokenType{scanner.Greater, scanner.GreaterEqual, scanner.Less, scanner.LessEqual}
	for _, op := range ops {
		*errors = nil
		e := &expr.Binary{
			Left:  &expr.Literal{Value: "hello"},
			Op:    tok(op),
			Right: &expr.Literal{Value: 1.0},
		}
		_, runtimeErr := evalSafe(interp, e)
		if runtimeErr == nil {
			t.Errorf("operator %v: expected runtime error for string operand, got none", op)
		}
	}
}

// --- stringify ---

// This test exposes the known bug: %.2f of 42.0 is "42.00", which does NOT
// end with ".0", so the integer check always fails. Whole numbers are
// formatted as "X.00" instead of "X".
func TestStringifyWholeNumber(t *testing.T) {
	interp, _ := newInterp(t)
	got := interp.stringify(42.0)
	want := "42"
	if got != want {
		t.Errorf("stringify(42.0): expected %q, got %q (known bug: check uses \".0\" but %%.2f produces \".00\")", want, got)
	}
}

func TestStringifyDecimalNumber(t *testing.T) {
	interp, _ := newInterp(t)
	got := interp.stringify(3.14)
	// With %.2f format, 3.14 → "3.14"
	if got != "3.14" {
		t.Errorf("stringify(3.14): expected \"3.14\", got %q", got)
	}
}

func TestStringifyNil(t *testing.T) {
	interp, _ := newInterp(t)
	got := interp.stringify(nil)
	if got != "nil" {
		t.Errorf("stringify(nil): expected \"nil\", got %q", got)
	}
}

func TestStringifyBool(t *testing.T) {
	interp, _ := newInterp(t)
	if got := interp.stringify(true); got != "true" {
		t.Errorf("stringify(true): expected \"true\", got %q", got)
	}
	if got := interp.stringify(false); got != "false" {
		t.Errorf("stringify(false): expected \"false\", got %q", got)
	}
}

func TestStringifyString(t *testing.T) {
	interp, _ := newInterp(t)
	if got := interp.stringify("hello"); got != "hello" {
		t.Errorf("stringify(\"hello\"): expected \"hello\", got %q", got)
	}
}

// --- isTruthy ---

func TestIsTruthy(t *testing.T) {
	interp, _ := newInterp(t)
	tests := []struct {
		value    any
		expected bool
	}{
		{nil, false},
		{false, false},
		{true, true},
		{0.0, true},  // zero number is still truthy in Lox
		{42.0, true},
		{"", true},   // empty string is truthy
		{"hello", true},
	}
	for _, tt := range tests {
		got := interp.isTruthy(tt.value)
		if got != tt.expected {
			t.Errorf("isTruthy(%v): expected %v, got %v", tt.value, tt.expected, got)
		}
	}
}

// --- Nested expressions ---

func TestNestedArithmetic(t *testing.T) {
	interp, _ := newInterp(t)
	// (1 + 2) * 3 = 9
	inner := &expr.Binary{
		Left:  &expr.Literal{Value: 1.0},
		Op:    tok(scanner.Plus),
		Right: &expr.Literal{Value: 2.0},
	}
	e := &expr.Binary{
		Left:  &expr.Grouping{Expr: inner},
		Op:    tok(scanner.Star),
		Right: &expr.Literal{Value: 3.0},
	}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != 9.0 {
		t.Errorf("expected 9.0, got %v", result)
	}
}

func TestDoubleNegation(t *testing.T) {
	interp, _ := newInterp(t)
	// !!true == true
	inner := &expr.Unary{Op: tok(scanner.Bang), Expr: &expr.Literal{Value: true}}
	e := &expr.Unary{Op: tok(scanner.Bang), Expr: inner}
	result, err := evalSafe(interp, e)
	if err != nil {
		t.Fatalf("unexpected runtime error: %v", err)
	}
	if result != true {
		t.Errorf("expected true for !!true, got %v", result)
	}
}
