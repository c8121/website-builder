package builder

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

type FunctionExpression struct {
	Name string
	Args []any
}

func NewFunctionExpression(expression string, firstArg any) (*FunctionExpression, error) {
	call, err := ParseFunctionExpression(expression)
	if err != nil {
		return nil, err
	}

	args := make([]any, len(call.Args)+1)
	args[0] = firstArg
	for i, arg := range call.Args {
		args[i+1] = arg
	}

	return &FunctionExpression{
		Name: call.Name,
		Args: args,
	}, nil
}

// ParseFunctionExpression creates an FunctionExpression from a call expression like `greetFirstname("ben")`.
// Supported argument types:
//   - strings:  "text", `raw text`
//   - integers: 42, -7, 0x1F       -> int
//   - floats:   3.14, -1e3         -> float64
//   - runes:    'a'                -> rune
//   - booleans: true, false        -> bool
//   - nil
func ParseFunctionExpression(s string) (FunctionExpression, error) {
	expr, err := parser.ParseExpr(strings.TrimSpace(s))
	if err != nil {
		return FunctionExpression{}, fmt.Errorf("invalid expression %q: %w", s, err)
	}

	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return FunctionExpression{}, fmt.Errorf("%q is not a function call", s)
	}

	// Only plain identifiers are allowed as function names (no pkg.Func, no method calls)
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return FunctionExpression{}, fmt.Errorf("invalid name: %s", types.ExprString(call.Fun))
	}

	// Reject spread syntax like f(xs...)
	if call.Ellipsis.IsValid() {
		return FunctionExpression{}, fmt.Errorf("'...' not supported")
	}

	args := make([]any, 0, len(call.Args))
	for i, a := range call.Args {
		v, err := parseArg(a)
		if err != nil {
			return FunctionExpression{}, fmt.Errorf("argument %d: %w", i, err)
		}
		args = append(args, v)
	}

	return FunctionExpression{Name: ident.Name, Args: args}, nil
}

// parseArg converts a single AST expression into a Go value.
func parseArg(e ast.Expr) (any, error) {
	switch e := e.(type) {
	case *ast.BasicLit:
		return parseLit(e.Kind, e.Value)

	case *ast.UnaryExpr:
		// Signed numbers like -5 or +3.2 are parsed as unary expressions
		lit, ok := e.X.(*ast.BasicLit)
		if ok && (e.Op == token.SUB || e.Op == token.ADD) &&
			(lit.Kind == token.INT || lit.Kind == token.FLOAT) {
			return parseLit(lit.Kind, e.Op.String()+lit.Value)
		}

	case *ast.Ident:
		switch e.Name {
		case "true":
			return true, nil
		case "false":
			return false, nil
		case "nil":
			return nil, nil
		}

	case *ast.ParenExpr:
		// Allow redundant parentheses, e.g. f((42))
		return parseArg(e.X)
	}

	return nil, fmt.Errorf("argument not supported: %s", types.ExprString(e))
}

// parseLit converts a literal token into its Go value.
func parseLit(kind token.Token, lit string) (any, error) {
	switch kind {
	case token.INT:
		v, err := strconv.ParseInt(lit, 0, 0) // base 0: supports 0x, 0o, 0b and _
		if err != nil {
			return nil, err
		}
		return int(v), nil

	case token.FLOAT:
		v, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			return nil, err
		}
		return v, nil

	case token.STRING:
		return strconv.Unquote(lit) // handles escapes and raw strings

	case token.CHAR:
		s, err := strconv.Unquote(lit)
		if err != nil {
			return nil, err
		}
		r, _ := utf8.DecodeRuneInString(s)
		return r, nil
	default:
		return nil, fmt.Errorf("literal not supported: %s", lit)
	}
}

func (e *FunctionExpression) InvokeFunctionExpression(functionMap map[string]any) (any, error) { //TODO variable args

	function, exists := functionMap[e.Name]
	if !exists {
		return "", errors.New("function not found")
	}

	v := reflect.ValueOf(function)
	if v.Kind() != reflect.Func {
		return "", fmt.Errorf("%q is not a function", e.Name)
	}
	t := v.Type()

	// Anzahl der Argumente prüfen
	numIn := t.NumIn()
	if t.IsVariadic() {
		if len(e.Args) < numIn-1 {
			return "", fmt.Errorf("%q expectes %d arguments, got %d", e.Name, numIn-1, len(e.Args))
		}
	} else if len(e.Args) != numIn {
		return "", fmt.Errorf("%q expectes %d arguments, got %d", e.Name, numIn, len(e.Args))
	}

	// Convert arg types
	in := make([]reflect.Value, len(e.Args))
	for i, a := range e.Args {
		var pt reflect.Type
		if t.IsVariadic() && i >= numIn-1 {
			pt = t.In(numIn - 1).Elem() // type of ...T
		} else {
			pt = t.In(i)
		}

		av, err := toValue(a, pt)
		if err != nil {
			return "", fmt.Errorf("%q, argument %d: %w", e.Name, i, err)
		}
		in[i] = av
	}

	return splitResults(v.Call(in), t)
}

// splitResults separates a trailing error and returns the remaining results.
func splitResults(out []reflect.Value, t reflect.Type) (any, error) {
	var err error
	if n := len(out); n > 0 && t.Out(n-1) == errorType {
		if e, _ := out[n-1].Interface().(error); e != nil {
			err = e
		}
		out = out[:n-1]
	}

	switch len(out) {
	case 0:
		return nil, err
	case 1:
		return out[0].Interface(), err
	default:
		res := make([]any, len(out))
		for i, o := range out {
			res[i] = o.Interface()
		}
		return res, err
	}
}

// toValue converts an arbitrary argument into a value of parameter type pt.
func toValue(a any, pt reflect.Type) (reflect.Value, error) {
	if a == nil {
		switch pt.Kind() {
		case reflect.Pointer, reflect.Interface, reflect.Slice,
			reflect.Map, reflect.Func, reflect.Chan:
			return reflect.Zero(pt), nil
		default:
			return reflect.Value{}, fmt.Errorf("nil is not allowed for %s", pt)
		}
	}

	av := reflect.ValueOf(a)
	at := av.Type()

	if at.AssignableTo(pt) {
		return av, nil
	}
	// Only allow numeric conversions (e.g. int -> int64, float64 -> int).
	// Deliberately NOT using ConvertibleTo in general, otherwise e.g.
	// int -> string would silently turn into a rune.
	if isNumeric(at.Kind()) && isNumeric(pt.Kind()) {
		return av.Convert(pt), nil
	}
	return reflect.Value{}, fmt.Errorf("typ %s does not match %s", at, pt)
}

// isNumeric reports whether k is an integer or floating-point kind.
func isNumeric(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
