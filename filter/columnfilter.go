package filter

import (
	"strconv"

	"github.com/kbence/logan/parser"
)

type Operator uint8

const (
	OpEquals Operator = iota
)

type ExpressionType uint8

const (
	TypeLiteral ExpressionType = iota
	TypeColumn
	TypeRelation
)

type Expression struct {
	Type    ExpressionType
	Parent  *Expression
	Left    *Expression
	Right   *Expression
	Op      Operator
	Literal string
}

func (e *Expression) GoLeft() *Expression {
	e.Left = &Expression{Parent: e}
	return e.Left
}

func (e *Expression) GoRight() *Expression {
	e.Right = &Expression{Parent: e}
	return e.Right
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
	switch e.Op {
	case OpEquals:
		return e.Left.EvaluateString(line) == e.Right.EvaluateString(line)
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

func (f *ColumnFilter) Match(line *parser.LogLine) bool {
	return f.expr.EvaluateBool(line)
}
