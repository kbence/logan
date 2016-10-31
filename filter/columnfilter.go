package filter

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/kbence/logan/parser"
)

type Operator uint8

const (
	OpEquals Operator = iota
	OpNotEquals
	OpMatchesRegexp
	OpOr
	OpAnd
)

var operatorStrings = []string{"==", "!=", "OR", "AND"}

type ExpressionType uint8

const (
	TypeLiteral ExpressionType = iota
	TypeColumn
	TypeRelation
	TypeLogical
)

type Expression struct {
	Type    ExpressionType
	Parent  *Expression
	Left    *Expression
	Right   *Expression
	Op      Operator
	Literal string
}

func (e *Expression) String() string {
	buf := bytes.NewBufferString("[Expression ")

	if e == nil {
		return "nil"
	}

	switch e.Type {
	case TypeLiteral:
		buf.WriteString(fmt.Sprintf("\"%s\"", e.Literal))
		break

	case TypeColumn:
		buf.WriteString(fmt.Sprintf("$%s", e.Literal))
		break

	case TypeRelation:
		buf.WriteString(fmt.Sprintf("%s %s %s", e.Left.String(),
			operatorStrings[int(e.Op)], e.Right.String()))
		break

	case TypeLogical:
		buf.WriteString(fmt.Sprintf("%s %s %s", e.Left.String(),
			operatorStrings[int(e.Op)], e.Right.String()))
		break
	}

	buf.WriteString("]")
	return buf.String()
}

func (e *Expression) GoLeft() *Expression {
	e.Left = &Expression{Parent: e}
	return e.Left
}

func (e *Expression) GoRight() *Expression {
	e.Right = &Expression{Parent: e}
	return e.Right
}

func (e *Expression) GoUp() *Expression {
	return e.Parent
}

func (e *Expression) PushLeftGoRight(op Operator) *Expression {
	newExpr := &Expression{Type: TypeLogical, Parent: e.Parent, Left: e, Op: op}
	newExpr.Right = &Expression{Parent: newExpr}

	if e.Parent != nil {
		e.Parent.Right = newExpr
	}

	return newExpr.Right
}

func (e *Expression) SetType(t ExpressionType) {
	e.Type = t
}

func (e *Expression) SetColumn(columnId string) {
	e.SetType(TypeColumn)
	e.Literal = columnId
}

func (e *Expression) SetString(str string) {
	e.SetType(TypeLiteral)
	e.Literal = str
}

func (e *Expression) EvaluateString(line *parser.LogLine) string {
	switch e.Type {
	case TypeLiteral:
		return e.Literal

	case TypeColumn:
		column, err := strconv.ParseInt(e.Literal, 10, 31)
		if err != nil || int(column) > len(line.Columns) {
			break
		}
		return line.Columns[column-1]
	}

	return ""
}

func (e *Expression) EvaluateBool(line *parser.LogLine) bool {
	switch e.Type {
	case TypeRelation:
		switch e.Op {
		case OpEquals:
			return e.Left.EvaluateString(line) == e.Right.EvaluateString(line)

		case OpNotEquals:
			return e.Left.EvaluateString(line) != e.Right.EvaluateString(line)

		case OpMatchesRegexp:
			re := regexp.MustCompile(e.Right.EvaluateString(line))
			return re.MatchString(e.Left.EvaluateString(line))
		}
		break

	case TypeLogical:
		switch e.Op {
		case OpAnd:
			return e.Left.EvaluateBool(line) && e.Right.EvaluateBool(line)

		case OpOr:
			return e.Left.EvaluateBool(line) || e.Right.EvaluateBool(line)
		}
	}

	return false
}

type ColumnFilter struct {
	expr *Expression
}

func NewColumnFilter(filterExpression string) *ColumnFilter {
	parser := &ColumnFilterParser{Buffer: filterExpression}
	parser.Init()
	parser.Parse()
	parser.Execute()

	return &ColumnFilter{expr: parser.Expr}
}

func (f *ColumnFilter) String() string {
	return f.expr.String()
}

func (f *ColumnFilter) Match(line *parser.LogLine) bool {
	return f.expr.EvaluateBool(line)
}
