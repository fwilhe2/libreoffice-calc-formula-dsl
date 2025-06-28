package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Expression Types

type Expr interface {
	Compile(argMap map[string]string) string
	String() string
}

type Number struct {
	Value string
}

type Variable struct {
	Name string
}

type BinaryOp struct {
	Left     Expr
	Operator string
	Right    Expr
}

type FunctionCall struct {
	Name string
	Args []Expr
}

func (n Number) Compile(_ map[string]string) string {
	return n.Value
}
func (n Number) String() string {
	return n.Value
}

func (v Variable) Compile(argMap map[string]string) string {
	if val, ok := argMap[v.Name]; ok {
		return val
	}
	return v.Name
}
func (v Variable) String() string {
	return v.Name
}

func (b BinaryOp) Compile(argMap map[string]string) string {
	return "(" + b.Left.Compile(argMap) + b.Operator + b.Right.Compile(argMap) + ")"
}
func (b BinaryOp) String() string {
	return "(" + b.Left.String() + " " + b.Operator + " " + b.Right.String() + ")"
}

func (f FunctionCall) Compile(argMap map[string]string) string {
	formula, ok := formulas[f.Name]
	if !ok {
		panic("Unknown function: " + f.Name)
	}
	if len(f.Args) != len(formula.Args) {
		panic("Argument count mismatch for function " + f.Name)
	}

	localArgs := map[string]string{}
	for i, arg := range formula.Args {
		localArgs[arg] = f.Args[i].Compile(argMap)
	}

	for k, v := range argMap {
		if _, ok := localArgs[k]; !ok {
			localArgs[k] = v
		}
	}

	return formula.Body.Compile(localArgs)
}

func (f FunctionCall) String() string {
	args := []string{}
	for _, a := range f.Args {
		args = append(args, a.String())
	}
	return f.Name + "(" + strings.Join(args, ", ") + ")"
}

// Formula Definition
type Formula struct {
	Args []string
	Body Expr
}

var constants = map[string]string{}
var formulas = map[string]Formula{}

func main() {
	source := `
		let TAX_RATE = 0.19
		define net_price(gross_price) = gross_price / (1 + TAX_RATE)
		define discount(price, percent) = price * (1 - percent / 100)
		define final_price(price, percent) = net_price(discount(price, percent))
	`

	parseDSL(source)

	fmt.Println("\nPretty-Printed AST for 'final_price':")
	f := formulas["final_price"]
	PrintAST(f.Body, "")

	compiled := compileFormula("final_price", []string{"A1", "A2"})
	fmt.Println("\nLibreOffice Calc Formula:\n", compiled)
}

func parseDSL(src string) {
	lines := strings.Split(src, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasPrefix(line, "let ") {
			parts := strings.SplitN(line[len("let "):], "=", 2)
			constants[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		} else if strings.HasPrefix(line, "define ") {
			parseFormula(line)
		}
	}
}

func parseFormula(line string) {
	parts := strings.SplitN(line[len("define "):], "=", 2)
	signature := strings.TrimSpace(parts[0])
	body := strings.TrimSpace(parts[1])

	open := strings.Index(signature, "(")
	close := strings.Index(signature, ")")
	name := signature[:open]
	args := strings.Split(signature[open+1:close], ",")

	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}

	parser := NewParser(body)
	expr := parser.ParseExpression()
	formulas[name] = Formula{Args: args, Body: expr}
}

func compileFormula(name string, argValues []string) string {
	formula, ok := formulas[name]
	if !ok {
		panic("Unknown formula: " + name)
	}
	if len(argValues) != len(formula.Args) {
		panic("Argument count mismatch for " + name)
	}

	argMap := map[string]string{}
	for i, arg := range formula.Args {
		argMap[arg] = argValues[i]
	}
	for k, v := range constants {
		argMap[k] = v
	}

	result := formula.Body.Compile(argMap)
	result = strings.ReplaceAll(result, " ", "")
	return "=" + result
}

// AST Printer
func PrintAST(expr Expr, indent string) {
	switch e := expr.(type) {
	case Number:
		fmt.Println(indent + "Number: " + e.Value)
	case Variable:
		fmt.Println(indent + "Variable: " + e.Name)
	case BinaryOp:
		fmt.Println(indent + "BinaryOp: " + e.Operator)
		PrintAST(e.Left, indent+"  ")
		PrintAST(e.Right, indent+"  ")
	case FunctionCall:
		fmt.Println(indent + "FunctionCall: " + e.Name)
		for _, arg := range e.Args {
			PrintAST(arg, indent+"  ")
		}
	default:
		fmt.Println(indent + "Unknown Node")
	}
}

// Parser Implementation
type Parser struct {
	input string
	pos   int
}

func NewParser(input string) *Parser {
	return &Parser{input: input}
}

func (p *Parser) ParseExpression() Expr {
	return p.parseAddSub()
}

func (p *Parser) parseAddSub() Expr {
	expr := p.parseMulDiv()
	for {
		p.skipWhitespace()
		if strings.HasPrefix(p.input[p.pos:], "+") {
			p.pos++
			right := p.parseMulDiv()
			expr = BinaryOp{Left: expr, Operator: "+", Right: right}
		} else if strings.HasPrefix(p.input[p.pos:], "-") {
			p.pos++
			right := p.parseMulDiv()
			expr = BinaryOp{Left: expr, Operator: "-", Right: right}
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) parseMulDiv() Expr {
	expr := p.parsePrimary()
	for {
		p.skipWhitespace()
		if strings.HasPrefix(p.input[p.pos:], "*") {
			p.pos++
			right := p.parsePrimary()
			expr = BinaryOp{Left: expr, Operator: "*", Right: right}
		} else if strings.HasPrefix(p.input[p.pos:], "/") {
			p.pos++
			right := p.parsePrimary()
			expr = BinaryOp{Left: expr, Operator: "/", Right: right}
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) parsePrimary() Expr {
	p.skipWhitespace()
	if strings.HasPrefix(p.input[p.pos:], "(") {
		p.pos++
		expr := p.ParseExpression()
		p.skipWhitespace()
		p.pos++ // skip ')'
		return expr
	}
	start := p.pos
	for p.pos < len(p.input) && (unicode.IsLetter(rune(p.input[p.pos])) || unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '.' || p.input[p.pos] == '_') {
		p.pos++
	}
	token := p.input[start:p.pos]
	p.skipWhitespace()
	if p.pos < len(p.input) && p.input[p.pos] == '(' {
		p.pos++
		args := []Expr{}
		for {
			p.skipWhitespace()
			if p.pos < len(p.input) && p.input[p.pos] == ')' {
				p.pos++
				break
			}
			arg := p.ParseExpression()
			args = append(args, arg)
			p.skipWhitespace()
			if p.pos < len(p.input) && p.input[p.pos] == ',' {
				p.pos++
			}
		}
		return FunctionCall{Name: token, Args: args}
	}

	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return Number{Value: token}
	}
	return Variable{Name: token}
}

func (p *Parser) skipWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}
