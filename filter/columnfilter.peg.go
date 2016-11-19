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

func (node *node32) Print(buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
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
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokens32
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
	p.tokens32.PrintSyntaxTree(p.Buffer)
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
	p.Reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.Reset()

	_rules, tree := p.rules, tokens32{tree: make([]token32, math.MaxInt16)}
	p.Parse = func(rule ...int) error {
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
		/* 0 filterExpression <- <(Action0 andFilterExpression (ws orOperator ws Action1 andFilterExpression Action2)* !.)> */
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
					if !_rules[ruleorOperator]() {
						goto l3
					}
					if !_rules[rulews]() {
						goto l3
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
					position4, tokenIndex4 := position, tokenIndex
					if !matchDot() {
						goto l4
					}
					goto l0
				l4:
					position, tokenIndex = position4, tokenIndex4
				}
				add(rulefilterExpression, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 andFilterExpression <- <(columnFilterExpression (ws andOperator ws Action3 columnFilterExpression Action4)*)> */
		func() bool {
			position5, tokenIndex5 := position, tokenIndex
			{
				position6 := position
				if !_rules[rulecolumnFilterExpression]() {
					goto l5
				}
			l7:
				{
					position8, tokenIndex8 := position, tokenIndex
					if !_rules[rulews]() {
						goto l8
					}
					if !_rules[ruleandOperator]() {
						goto l8
					}
					if !_rules[rulews]() {
						goto l8
					}
					if !_rules[ruleAction3]() {
						goto l8
					}
					if !_rules[rulecolumnFilterExpression]() {
						goto l8
					}
					if !_rules[ruleAction4]() {
						goto l8
					}
					goto l7
				l8:
					position, tokenIndex = position8, tokenIndex8
				}
				add(ruleandFilterExpression, position6)
			}
			return true
		l5:
			position, tokenIndex = position5, tokenIndex5
			return false
		},
		/* 2 columnFilterExpression <- <relation> */
		func() bool {
			position9, tokenIndex9 := position, tokenIndex
			{
				position10 := position
				if !_rules[rulerelation]() {
					goto l9
				}
				add(rulecolumnFilterExpression, position10)
			}
			return true
		l9:
			position, tokenIndex = position9, tokenIndex9
			return false
		},
		/* 3 relation <- <(Action5 expression Action6 ws Action7 relationOperator ws Action8 expression Action9)> */
		func() bool {
			position11, tokenIndex11 := position, tokenIndex
			{
				position12 := position
				if !_rules[ruleAction5]() {
					goto l11
				}
				if !_rules[ruleexpression]() {
					goto l11
				}
				if !_rules[ruleAction6]() {
					goto l11
				}
				if !_rules[rulews]() {
					goto l11
				}
				if !_rules[ruleAction7]() {
					goto l11
				}
				if !_rules[rulerelationOperator]() {
					goto l11
				}
				if !_rules[rulews]() {
					goto l11
				}
				if !_rules[ruleAction8]() {
					goto l11
				}
				if !_rules[ruleexpression]() {
					goto l11
				}
				if !_rules[ruleAction9]() {
					goto l11
				}
				add(rulerelation, position12)
			}
			return true
		l11:
			position, tokenIndex = position11, tokenIndex11
			return false
		},
		/* 4 expression <- <(columnSpecifier / literal)> */
		func() bool {
			position13, tokenIndex13 := position, tokenIndex
			{
				position14 := position
				{
					position15, tokenIndex15 := position, tokenIndex
					if !_rules[rulecolumnSpecifier]() {
						goto l16
					}
					goto l15
				l16:
					position, tokenIndex = position15, tokenIndex15
					if !_rules[ruleliteral]() {
						goto l13
					}
				}
			l15:
				add(ruleexpression, position14)
			}
			return true
		l13:
			position, tokenIndex = position13, tokenIndex13
			return false
		},
		/* 5 columnSpecifier <- <('$' <[0-9]+> Action10)> */
		func() bool {
			position17, tokenIndex17 := position, tokenIndex
			{
				position18 := position
				if buffer[position] != rune('$') {
					goto l17
				}
				position++
				{
					position19 := position
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l17
					}
					position++
				l20:
					{
						position21, tokenIndex21 := position, tokenIndex
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l21
						}
						position++
						goto l20
					l21:
						position, tokenIndex = position21, tokenIndex21
					}
					add(rulePegText, position19)
				}
				if !_rules[ruleAction10]() {
					goto l17
				}
				add(rulecolumnSpecifier, position18)
			}
			return true
		l17:
			position, tokenIndex = position17, tokenIndex17
			return false
		},
		/* 6 literal <- <stringLiteral> */
		func() bool {
			position22, tokenIndex22 := position, tokenIndex
			{
				position23 := position
				if !_rules[rulestringLiteral]() {
					goto l22
				}
				add(ruleliteral, position23)
			}
			return true
		l22:
			position, tokenIndex = position22, tokenIndex22
			return false
		},
		/* 7 stringLiteral <- <('"' <stringContent> '"' Action11)> */
		func() bool {
			position24, tokenIndex24 := position, tokenIndex
			{
				position25 := position
				if buffer[position] != rune('"') {
					goto l24
				}
				position++
				{
					position26 := position
					if !_rules[rulestringContent]() {
						goto l24
					}
					add(rulePegText, position26)
				}
				if buffer[position] != rune('"') {
					goto l24
				}
				position++
				if !_rules[ruleAction11]() {
					goto l24
				}
				add(rulestringLiteral, position25)
			}
			return true
		l24:
			position, tokenIndex = position24, tokenIndex24
			return false
		},
		/* 8 stringContent <- <((!'"' .) / ('\\' '"'))+> */
		func() bool {
			position27, tokenIndex27 := position, tokenIndex
			{
				position28 := position
				{
					position31, tokenIndex31 := position, tokenIndex
					{
						position33, tokenIndex33 := position, tokenIndex
						if buffer[position] != rune('"') {
							goto l33
						}
						position++
						goto l32
					l33:
						position, tokenIndex = position33, tokenIndex33
					}
					if !matchDot() {
						goto l32
					}
					goto l31
				l32:
					position, tokenIndex = position31, tokenIndex31
					if buffer[position] != rune('\\') {
						goto l27
					}
					position++
					if buffer[position] != rune('"') {
						goto l27
					}
					position++
				}
			l31:
			l29:
				{
					position30, tokenIndex30 := position, tokenIndex
					{
						position34, tokenIndex34 := position, tokenIndex
						{
							position36, tokenIndex36 := position, tokenIndex
							if buffer[position] != rune('"') {
								goto l36
							}
							position++
							goto l35
						l36:
							position, tokenIndex = position36, tokenIndex36
						}
						if !matchDot() {
							goto l35
						}
						goto l34
					l35:
						position, tokenIndex = position34, tokenIndex34
						if buffer[position] != rune('\\') {
							goto l30
						}
						position++
						if buffer[position] != rune('"') {
							goto l30
						}
						position++
					}
				l34:
					goto l29
				l30:
					position, tokenIndex = position30, tokenIndex30
				}
				add(rulestringContent, position28)
			}
			return true
		l27:
			position, tokenIndex = position27, tokenIndex27
			return false
		},
		/* 9 relationOperator <- <(equals / notEquals / matchesRegexp)> */
		func() bool {
			position37, tokenIndex37 := position, tokenIndex
			{
				position38 := position
				{
					position39, tokenIndex39 := position, tokenIndex
					if !_rules[ruleequals]() {
						goto l40
					}
					goto l39
				l40:
					position, tokenIndex = position39, tokenIndex39
					if !_rules[rulenotEquals]() {
						goto l41
					}
					goto l39
				l41:
					position, tokenIndex = position39, tokenIndex39
					if !_rules[rulematchesRegexp]() {
						goto l37
					}
				}
			l39:
				add(rulerelationOperator, position38)
			}
			return true
		l37:
			position, tokenIndex = position37, tokenIndex37
			return false
		},
		/* 10 equals <- <('=' '=' Action12)> */
		func() bool {
			position42, tokenIndex42 := position, tokenIndex
			{
				position43 := position
				if buffer[position] != rune('=') {
					goto l42
				}
				position++
				if buffer[position] != rune('=') {
					goto l42
				}
				position++
				if !_rules[ruleAction12]() {
					goto l42
				}
				add(ruleequals, position43)
			}
			return true
		l42:
			position, tokenIndex = position42, tokenIndex42
			return false
		},
		/* 11 notEquals <- <((('!' '=') / ('<' '>')) Action13)> */
		func() bool {
			position44, tokenIndex44 := position, tokenIndex
			{
				position45 := position
				{
					position46, tokenIndex46 := position, tokenIndex
					if buffer[position] != rune('!') {
						goto l47
					}
					position++
					if buffer[position] != rune('=') {
						goto l47
					}
					position++
					goto l46
				l47:
					position, tokenIndex = position46, tokenIndex46
					if buffer[position] != rune('<') {
						goto l44
					}
					position++
					if buffer[position] != rune('>') {
						goto l44
					}
					position++
				}
			l46:
				if !_rules[ruleAction13]() {
					goto l44
				}
				add(rulenotEquals, position45)
			}
			return true
		l44:
			position, tokenIndex = position44, tokenIndex44
			return false
		},
		/* 12 matchesRegexp <- <('~' '=' Action14)> */
		func() bool {
			position48, tokenIndex48 := position, tokenIndex
			{
				position49 := position
				if buffer[position] != rune('~') {
					goto l48
				}
				position++
				if buffer[position] != rune('=') {
					goto l48
				}
				position++
				if !_rules[ruleAction14]() {
					goto l48
				}
				add(rulematchesRegexp, position49)
			}
			return true
		l48:
			position, tokenIndex = position48, tokenIndex48
			return false
		},
		/* 13 andOperator <- <(('a' 'n' 'd') / ('A' 'N' 'D') / ('&' '&'))> */
		func() bool {
			position50, tokenIndex50 := position, tokenIndex
			{
				position51 := position
				{
					position52, tokenIndex52 := position, tokenIndex
					if buffer[position] != rune('a') {
						goto l53
					}
					position++
					if buffer[position] != rune('n') {
						goto l53
					}
					position++
					if buffer[position] != rune('d') {
						goto l53
					}
					position++
					goto l52
				l53:
					position, tokenIndex = position52, tokenIndex52
					if buffer[position] != rune('A') {
						goto l54
					}
					position++
					if buffer[position] != rune('N') {
						goto l54
					}
					position++
					if buffer[position] != rune('D') {
						goto l54
					}
					position++
					goto l52
				l54:
					position, tokenIndex = position52, tokenIndex52
					if buffer[position] != rune('&') {
						goto l50
					}
					position++
					if buffer[position] != rune('&') {
						goto l50
					}
					position++
				}
			l52:
				add(ruleandOperator, position51)
			}
			return true
		l50:
			position, tokenIndex = position50, tokenIndex50
			return false
		},
		/* 14 orOperator <- <(('o' 'r') / ('O' 'R') / ('|' '|'))> */
		func() bool {
			position55, tokenIndex55 := position, tokenIndex
			{
				position56 := position
				{
					position57, tokenIndex57 := position, tokenIndex
					if buffer[position] != rune('o') {
						goto l58
					}
					position++
					if buffer[position] != rune('r') {
						goto l58
					}
					position++
					goto l57
				l58:
					position, tokenIndex = position57, tokenIndex57
					if buffer[position] != rune('O') {
						goto l59
					}
					position++
					if buffer[position] != rune('R') {
						goto l59
					}
					position++
					goto l57
				l59:
					position, tokenIndex = position57, tokenIndex57
					if buffer[position] != rune('|') {
						goto l55
					}
					position++
					if buffer[position] != rune('|') {
						goto l55
					}
					position++
				}
			l57:
				add(ruleorOperator, position56)
			}
			return true
		l55:
			position, tokenIndex = position55, tokenIndex55
			return false
		},
		/* 15 ws <- <(' ' / '\t')> */
		func() bool {
			position60, tokenIndex60 := position, tokenIndex
			{
				position61 := position
				{
					position62, tokenIndex62 := position, tokenIndex
					if buffer[position] != rune(' ') {
						goto l63
					}
					position++
					goto l62
				l63:
					position, tokenIndex = position62, tokenIndex62
					if buffer[position] != rune('\t') {
						goto l60
					}
					position++
				}
			l62:
				add(rulews, position61)
			}
			return true
		l60:
			position, tokenIndex = position60, tokenIndex60
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
