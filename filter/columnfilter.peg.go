package filter

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulefilterExpression
	ruleandFilterExpression
	rulecolumnFilterExpression
	rulerelation
	ruleexpression
	rulecolumnSpecifier
	ruleliteral
	rulestringLiteral
	rulestringContent
	rulerelationOperator
	ruleequals
	rulenotEquals
	rulematchesRegexp
	ruleandOperator
	ruleorOperator
	rulews
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	rulePegText
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
)

var rul3s = [...]string{
	"Unknown",
	"filterExpression",
	"andFilterExpression",
	"columnFilterExpression",
	"relation",
	"expression",
	"columnSpecifier",
	"literal",
	"stringLiteral",
	"stringContent",
	"relationOperator",
	"equals",
	"notEquals",
	"matchesRegexp",
	"andOperator",
	"orOperator",
	"ws",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"PegText",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type ColumnFilterParser struct {
	Expr *Expression

	Buffer string
	buffer []rune
	rules  [33]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *ColumnFilterParser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *ColumnFilterParser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *ColumnFilterParser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *ColumnFilterParser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *ColumnFilterParser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.Expr = &Expression{}
		case ruleAction1:
			p.Expr = p.Expr.PushLeftGoRight(OpOr)
		case ruleAction2:
			p.Expr = p.Expr.GoUp()
		case ruleAction3:
			p.Expr = p.Expr.PushLeftGoRight(OpAnd)
		case ruleAction4:
			p.Expr = p.Expr.GoUp()
		case ruleAction5:
			p.Expr = p.Expr.GoLeft()
		case ruleAction6:
			p.Expr = p.Expr.GoUp()
		case ruleAction7:
			p.Expr.SetType(TypeRelation)
		case ruleAction8:
			p.Expr = p.Expr.GoRight()
		case ruleAction9:
			p.Expr = p.Expr.GoUp()
		case ruleAction10:
			p.Expr.SetColumn(buffer[begin:end])
		case ruleAction11:
			p.Expr.SetString(buffer[begin:end])
		case ruleAction12:
			p.Expr.Op = OpEquals
		case ruleAction13:
			p.Expr.Op = OpNotEquals
		case ruleAction14:
			p.Expr.Op = OpMatchesRegexp

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *ColumnFilterParser) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 filterExpression <- <(Action0 andFilterExpression (ws+ orOperator ws+ Action1 andFilterExpression Action2)* !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[ruleAction0]() {
					goto l0
				}
				if !_rules[ruleandFilterExpression]() {
					goto l0
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					if !_rules[rulews]() {
						goto l3
					}
				l4:
					{
						position5, tokenIndex5 := position, tokenIndex
						if !_rules[rulews]() {
							goto l5
						}
						goto l4
					l5:
						position, tokenIndex = position5, tokenIndex5
					}
					if !_rules[ruleorOperator]() {
						goto l3
					}
					if !_rules[rulews]() {
						goto l3
					}
				l6:
					{
						position7, tokenIndex7 := position, tokenIndex
						if !_rules[rulews]() {
							goto l7
						}
						goto l6
					l7:
						position, tokenIndex = position7, tokenIndex7
					}
					if !_rules[ruleAction1]() {
						goto l3
					}
					if !_rules[ruleandFilterExpression]() {
						goto l3
					}
					if !_rules[ruleAction2]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				{
					position8, tokenIndex8 := position, tokenIndex
					if !matchDot() {
						goto l8
					}
					goto l0
				l8:
					position, tokenIndex = position8, tokenIndex8
				}
				add(rulefilterExpression, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 andFilterExpression <- <(columnFilterExpression (ws+ andOperator ws+ Action3 columnFilterExpression Action4)*)> */
		func() bool {
			position9, tokenIndex9 := position, tokenIndex
			{
				position10 := position
				if !_rules[rulecolumnFilterExpression]() {
					goto l9
				}
			l11:
				{
					position12, tokenIndex12 := position, tokenIndex
					if !_rules[rulews]() {
						goto l12
					}
				l13:
					{
						position14, tokenIndex14 := position, tokenIndex
						if !_rules[rulews]() {
							goto l14
						}
						goto l13
					l14:
						position, tokenIndex = position14, tokenIndex14
					}
					if !_rules[ruleandOperator]() {
						goto l12
					}
					if !_rules[rulews]() {
						goto l12
					}
				l15:
					{
						position16, tokenIndex16 := position, tokenIndex
						if !_rules[rulews]() {
							goto l16
						}
						goto l15
					l16:
						position, tokenIndex = position16, tokenIndex16
					}
					if !_rules[ruleAction3]() {
						goto l12
					}
					if !_rules[rulecolumnFilterExpression]() {
						goto l12
					}
					if !_rules[ruleAction4]() {
						goto l12
					}
					goto l11
				l12:
					position, tokenIndex = position12, tokenIndex12
				}
				add(ruleandFilterExpression, position10)
			}
			return true
		l9:
			position, tokenIndex = position9, tokenIndex9
			return false
		},
		/* 2 columnFilterExpression <- <relation> */
		func() bool {
			position17, tokenIndex17 := position, tokenIndex
			{
				position18 := position
				if !_rules[rulerelation]() {
					goto l17
				}
				add(rulecolumnFilterExpression, position18)
			}
			return true
		l17:
			position, tokenIndex = position17, tokenIndex17
			return false
		},
		/* 3 relation <- <(Action5 expression Action6 ws* Action7 relationOperator ws* Action8 expression Action9)> */
		func() bool {
			position19, tokenIndex19 := position, tokenIndex
			{
				position20 := position
				if !_rules[ruleAction5]() {
					goto l19
				}
				if !_rules[ruleexpression]() {
					goto l19
				}
				if !_rules[ruleAction6]() {
					goto l19
				}
			l21:
				{
					position22, tokenIndex22 := position, tokenIndex
					if !_rules[rulews]() {
						goto l22
					}
					goto l21
				l22:
					position, tokenIndex = position22, tokenIndex22
				}
				if !_rules[ruleAction7]() {
					goto l19
				}
				if !_rules[rulerelationOperator]() {
					goto l19
				}
			l23:
				{
					position24, tokenIndex24 := position, tokenIndex
					if !_rules[rulews]() {
						goto l24
					}
					goto l23
				l24:
					position, tokenIndex = position24, tokenIndex24
				}
				if !_rules[ruleAction8]() {
					goto l19
				}
				if !_rules[ruleexpression]() {
					goto l19
				}
				if !_rules[ruleAction9]() {
					goto l19
				}
				add(rulerelation, position20)
			}
			return true
		l19:
			position, tokenIndex = position19, tokenIndex19
			return false
		},
		/* 4 expression <- <(columnSpecifier / literal)> */
		func() bool {
			position25, tokenIndex25 := position, tokenIndex
			{
				position26 := position
				{
					position27, tokenIndex27 := position, tokenIndex
					if !_rules[rulecolumnSpecifier]() {
						goto l28
					}
					goto l27
				l28:
					position, tokenIndex = position27, tokenIndex27
					if !_rules[ruleliteral]() {
						goto l25
					}
				}
			l27:
				add(ruleexpression, position26)
			}
			return true
		l25:
			position, tokenIndex = position25, tokenIndex25
			return false
		},
		/* 5 columnSpecifier <- <('$' <[0-9]+> Action10)> */
		func() bool {
			position29, tokenIndex29 := position, tokenIndex
			{
				position30 := position
				if buffer[position] != rune('$') {
					goto l29
				}
				position++
				{
					position31 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l29
					}
					position++
				l32:
					{
						position33, tokenIndex33 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l33
						}
						position++
						goto l32
					l33:
						position, tokenIndex = position33, tokenIndex33
					}
					add(rulePegText, position31)
				}
				if !_rules[ruleAction10]() {
					goto l29
				}
				add(rulecolumnSpecifier, position30)
			}
			return true
		l29:
			position, tokenIndex = position29, tokenIndex29
			return false
		},
		/* 6 literal <- <stringLiteral> */
		func() bool {
			position34, tokenIndex34 := position, tokenIndex
			{
				position35 := position
				if !_rules[rulestringLiteral]() {
					goto l34
				}
				add(ruleliteral, position35)
			}
			return true
		l34:
			position, tokenIndex = position34, tokenIndex34
			return false
		},
		/* 7 stringLiteral <- <('"' <stringContent> '"' Action11)> */
		func() bool {
			position36, tokenIndex36 := position, tokenIndex
			{
				position37 := position
				if buffer[position] != rune('"') {
					goto l36
				}
				position++
				{
					position38 := position
					if !_rules[rulestringContent]() {
						goto l36
					}
					add(rulePegText, position38)
				}
				if buffer[position] != rune('"') {
					goto l36
				}
				position++
				if !_rules[ruleAction11]() {
					goto l36
				}
				add(rulestringLiteral, position37)
			}
			return true
		l36:
			position, tokenIndex = position36, tokenIndex36
			return false
		},
		/* 8 stringContent <- <((!'"' .) / ('\\' '"'))+> */
		func() bool {
			position39, tokenIndex39 := position, tokenIndex
			{
				position40 := position
				{
					position43, tokenIndex43 := position, tokenIndex
					{
						position45, tokenIndex45 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l45
						}
						position++
						goto l44
					l45:
						position, tokenIndex = position45, tokenIndex45
					}
					if !matchDot() {
						goto l44
					}
					goto l43
				l44:
					position, tokenIndex = position43, tokenIndex43
					if buffer[position] != rune('\\') {
						goto l39
					}
					position++
					if buffer[position] != rune('"') {
						goto l39
					}
					position++
				}
			l43:
			l41:
				{
					position42, tokenIndex42 := position, tokenIndex
					{
						position46, tokenIndex46 := position, tokenIndex
						{
							position48, tokenIndex48 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l48
							}
							position++
							goto l47
						l48:
							position, tokenIndex = position48, tokenIndex48
						}
						if !matchDot() {
							goto l47
						}
						goto l46
					l47:
						position, tokenIndex = position46, tokenIndex46
						if buffer[position] != rune('\\') {
							goto l42
						}
						position++
						if buffer[position] != rune('"') {
							goto l42
						}
						position++
					}
				l46:
					goto l41
				l42:
					position, tokenIndex = position42, tokenIndex42
				}
				add(rulestringContent, position40)
			}
			return true
		l39:
			position, tokenIndex = position39, tokenIndex39
			return false
		},
		/* 9 relationOperator <- <(equals / notEquals / matchesRegexp)> */
		func() bool {
			position49, tokenIndex49 := position, tokenIndex
			{
				position50 := position
				{
					position51, tokenIndex51 := position, tokenIndex
					if !_rules[ruleequals]() {
						goto l52
					}
					goto l51
				l52:
					position, tokenIndex = position51, tokenIndex51
					if !_rules[rulenotEquals]() {
						goto l53
					}
					goto l51
				l53:
					position, tokenIndex = position51, tokenIndex51
					if !_rules[rulematchesRegexp]() {
						goto l49
					}
				}
			l51:
				add(rulerelationOperator, position50)
			}
			return true
		l49:
			position, tokenIndex = position49, tokenIndex49
			return false
		},
		/* 10 equals <- <('=' '=' Action12)> */
		func() bool {
			position54, tokenIndex54 := position, tokenIndex
			{
				position55 := position
				if buffer[position] != rune('=') {
					goto l54
				}
				position++
				if buffer[position] != rune('=') {
					goto l54
				}
				position++
				if !_rules[ruleAction12]() {
					goto l54
				}
				add(ruleequals, position55)
			}
			return true
		l54:
			position, tokenIndex = position54, tokenIndex54
			return false
		},
		/* 11 notEquals <- <((('!' '=') / ('<' '>')) Action13)> */
		func() bool {
			position56, tokenIndex56 := position, tokenIndex
			{
				position57 := position
				{
					position58, tokenIndex58 := position, tokenIndex
					if buffer[position] != rune('!') {
						goto l59
					}
					position++
					if buffer[position] != rune('=') {
						goto l59
					}
					position++
					goto l58
				l59:
					position, tokenIndex = position58, tokenIndex58
					if buffer[position] != rune('<') {
						goto l56
					}
					position++
					if buffer[position] != rune('>') {
						goto l56
					}
					position++
				}
			l58:
				if !_rules[ruleAction13]() {
					goto l56
				}
				add(rulenotEquals, position57)
			}
			return true
		l56:
			position, tokenIndex = position56, tokenIndex56
			return false
		},
		/* 12 matchesRegexp <- <('~' '=' Action14)> */
		func() bool {
			position60, tokenIndex60 := position, tokenIndex
			{
				position61 := position
				if buffer[position] != rune('~') {
					goto l60
				}
				position++
				if buffer[position] != rune('=') {
					goto l60
				}
				position++
				if !_rules[ruleAction14]() {
					goto l60
				}
				add(rulematchesRegexp, position61)
			}
			return true
		l60:
			position, tokenIndex = position60, tokenIndex60
			return false
		},
		/* 13 andOperator <- <(('a' 'n' 'd') / ('A' 'N' 'D') / ('&' '&'))> */
		func() bool {
			position62, tokenIndex62 := position, tokenIndex
			{
				position63 := position
				{
					position64, tokenIndex64 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l65
					}
					position++
					if buffer[position] != rune('n') {
						goto l65
					}
					position++
					if buffer[position] != rune('d') {
						goto l65
					}
					position++
					goto l64
				l65:
					position, tokenIndex = position64, tokenIndex64
					if buffer[position] != rune('A') {
						goto l66
					}
					position++
					if buffer[position] != rune('N') {
						goto l66
					}
					position++
					if buffer[position] != rune('D') {
						goto l66
					}
					position++
					goto l64
				l66:
					position, tokenIndex = position64, tokenIndex64
					if buffer[position] != rune('&') {
						goto l62
					}
					position++
					if buffer[position] != rune('&') {
						goto l62
					}
					position++
				}
			l64:
				add(ruleandOperator, position63)
			}
			return true
		l62:
			position, tokenIndex = position62, tokenIndex62
			return false
		},
		/* 14 orOperator <- <(('o' 'r') / ('O' 'R') / ('|' '|'))> */
		func() bool {
			position67, tokenIndex67 := position, tokenIndex
			{
				position68 := position
				{
					position69, tokenIndex69 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l70
					}
					position++
					if buffer[position] != rune('r') {
						goto l70
					}
					position++
					goto l69
				l70:
					position, tokenIndex = position69, tokenIndex69
					if buffer[position] != rune('O') {
						goto l71
					}
					position++
					if buffer[position] != rune('R') {
						goto l71
					}
					position++
					goto l69
				l71:
					position, tokenIndex = position69, tokenIndex69
					if buffer[position] != rune('|') {
						goto l67
					}
					position++
					if buffer[position] != rune('|') {
						goto l67
					}
					position++
				}
			l69:
				add(ruleorOperator, position68)
			}
			return true
		l67:
			position, tokenIndex = position67, tokenIndex67
			return false
		},
		/* 15 ws <- <(' ' / '\t')> */
		func() bool {
			position72, tokenIndex72 := position, tokenIndex
			{
				position73 := position
				{
					position74, tokenIndex74 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l75
					}
					position++
					goto l74
				l75:
					position, tokenIndex = position74, tokenIndex74
					if buffer[position] != rune('\t') {
						goto l72
					}
					position++
				}
			l74:
				add(rulews, position73)
			}
			return true
		l72:
			position, tokenIndex = position72, tokenIndex72
			return false
		},
		/* 17 Action0 <- <{ p.Expr = &Expression{} }> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		/* 18 Action1 <- <{ p.Expr = p.Expr.PushLeftGoRight(OpOr) }> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 19 Action2 <- <{ p.Expr = p.Expr.GoUp() }> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 20 Action3 <- <{ p.Expr = p.Expr.PushLeftGoRight(OpAnd) }> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 21 Action4 <- <{ p.Expr = p.Expr.GoUp() }> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 22 Action5 <- <{ p.Expr = p.Expr.GoLeft() }> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 23 Action6 <- <{ p.Expr = p.Expr.GoUp() }> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 24 Action7 <- <{ p.Expr.SetType(TypeRelation) }> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 25 Action8 <- <{ p.Expr = p.Expr.GoRight() }> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 26 Action9 <- <{ p.Expr = p.Expr.GoUp() }> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		nil,
		/* 28 Action10 <- <{ p.Expr.SetColumn(buffer[begin:end]) }> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 29 Action11 <- <{ p.Expr.SetString(buffer[begin:end]) }> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 30 Action12 <- <{ p.Expr.Op = OpEquals }> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 31 Action13 <- <{ p.Expr.Op = OpNotEquals }> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 32 Action14 <- <{ p.Expr.Op = OpMatchesRegexp }> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
	}
	p.rules = _rules
}
