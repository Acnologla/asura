package interpreter

import (
	"asura/src/utils"
	"fmt"
)

type Token struct {
	Left  interface{}
	Value string
	Right interface{}
}

type Parser struct {
	Lexems []*Lexem
	Error  string
	I      int
}

var stringTypes = map[int]string{
	1: "Boolean",
	2: "Int",
	3: "String",
	4: "Newline",
	5: "Token",
	6: "Keyword",
	7: "Identifier",
}

var comparators = []string{">", "==", "!=", "<"}

func (this *Parser) CreateToken(left interface{}, value string, right interface{}) *Token {
	return &Token{
		Left:  left,
		Value: value,
		Right: right,
	}
}

func (this *Parser) Current() *Lexem {
	if len(this.Lexems) == this.I {
		return &Lexem{}
	}
	return this.Lexems[this.I]
}

func (this *Parser) Eat(lexemType int) *Lexem {
	if len(this.Lexems) == this.I {
		if lexemType != 4 {
			this.Error = fmt.Sprintf("Invalid type expected: %s\n Value: EOF", stringTypes[lexemType])
		}
		return &Lexem{}
	}
	if this.Lexems[this.I].Type != lexemType {
		this.Error = fmt.Sprintf("Invalid type expected: %s, got %s\n Value: %s", stringTypes[lexemType], stringTypes[this.Lexems[this.I].Type], this.Lexems[this.I].Value)
	} else {
		this.I++
		return this.Lexems[this.I-1]
	}
	return nil
}

func (this *Parser) accAcess(name *Lexem, operator *Lexem) *Token {
	var tok *Token
	this.JumpBreaks()
	if operator.Value == "." {
		tok = this.CreateToken(name.Value, "property", this.Add())
	} else {
		n := this.Parse()
		tok = this.CreateToken(name.Value, "index", n)
		this.Eat(5)
		this.JumpBreaks()
		if this.Current().Type == 5 && this.Current().Value == "=" {
			this.Eat(5)
			value := this.Parse()
			return this.CreateToken(tok, "changeArr", value)
		}
	}
	for this.Current().Type == 5 && (this.Current().Value == "." || this.Current().Value == "[") {
		if this.Current().Value == "." {
			this.Eat(5)
			tok = this.CreateToken(tok, "property", this.Add())
		} else {
			this.Eat(5)
			n := this.Parse()
			tok = this.CreateToken(tok, "index", n)
			this.Eat(5)
			this.JumpBreaks()
			if this.Current().Type == 5 && this.Current().Value == "=" {
				this.Eat(5)
				value := this.Parse()
				return this.CreateToken(tok, "changeArr", value)
			}
		}
	}
	return tok
}

func (this *Parser) AccID() *Token {
	name := this.Eat(7)
	if this.Current().Type == 5 && this.Current().Value != ")" {
		operator := this.Eat(5)
		var value interface{} = name.Value
		if operator.Value == "." || operator.Value == "[" {
			tok := this.accAcess(name, operator)
			if this.Current().Type == 5  && this.Current().Value != ")"{
				value = tok
				operator = this.Eat(5)
			} else {
				return tok
			}
		}
		if operator.Value == "++" {
			return this.CreateToken(value, "=", this.CreateToken(value, "+", this.CreateToken(nil, "1", nil)))
		}
		if operator.Value == "--" {
			return this.CreateToken(value, "=", this.CreateToken(value, "-", this.CreateToken(nil, "1", nil)))
		}
		if operator.Value == "="{
			return this.CreateToken(value,"=",this.Parse())
		}
		if operator.Value == ":=" {
			return this.CreateToken(value, ":=", this.Parse())
		}
		if operator.Value == "(" {
			var tok *Token
			if this.Current().Type == 5 && this.Current().Value == ")" {
				tok = this.CreateToken(value, "call", nil)
			} else {
				params := []*Token{this.Parse()}
				for this.Current().Type == 5 && this.Current().Value == "," {
					this.Eat(5)
					params = append(params, this.Parse())
				}
				tok = this.CreateToken(value, "call", params)
			}
			this.Eat(5)
			for this.Current().Type == 5 && (this.Current().Value == "." || this.Current().Value == "[") {
				if this.Current().Type == 5 && this.Current().Value == "." {
					this.Eat(5)
					this.JumpBreaks()
					tok = this.CreateToken(tok, "property", this.Add())
				}
				if this.Current().Type == 5 && this.Current().Value == "[" {
					this.Eat(5)
					this.JumpBreaks()
					nTok := this.Add()
					this.Eat(5)
					tok = this.CreateToken(tok, "index", nTok)
				}
			}
			return tok
		}
		this.I--
		return this.CreateToken(nil, name.Value, nil)
	}
	return this.CreateToken(nil, name.Value, nil)
}

func (this *Parser) Factor() *Token {
	current := this.Current()
	if current.Type == 5 && current.Value == "(" {
		this.Eat(5)
		expr := this.Expr()
		this.Eat(5)
		return expr
	}
	if current.Type == 7 {
		return this.AccID()
	} else if current.Type == 3 {
		this.Eat(3)
		return this.CreateToken(nil, "\""+current.Value+"\"", nil)
	} else {
		this.Eat(2)
	}
	return this.CreateToken(nil, current.Value, nil)
}

func (this *Parser) Mult() *Token {
	old := this.Factor()
	for this.Current().Type == 5 && (this.Current().Value == "*" || this.Current().Value == "/" || this.Current().Value == "%") {
		op := this.Eat(5)
		old = this.CreateToken(old, op.Value, this.Factor())
	}
	return old
}

func (this *Parser) Add() *Token {
	old := this.Mult()
	for this.Current().Type == 5 && (this.Current().Value == "+" || this.Current().Value == "-") {
		op := this.Eat(5)
		old = this.CreateToken(old, op.Value, this.Mult())
	}
	return old
}

func (this *Parser) Expr() *Token {
	old := this.Add()
	for this.Current().Type == 5 && utils.Includes(comparators, this.Current().Value) {
		op := this.Eat(5)
		old = this.CreateToken(old, op.Value, this.Add())
	}
	return old
}

func (this *Parser) Compound() (tokens []*Token) {
	for this.I < len(this.Lexems) {
		result := this.Parse()
		if result != nil {
			if result.Value == "}" {
				return
			} else {
				tokens = append(tokens, result)
			}
		}
	}
	return
}

func (this *Parser) JumpBreaks() {
	for this.Current().Type == 4 {
		this.Eat(4)
	}
}

func (this *Parser) Parse() *Token {
	current := this.Current()
	if current.Type == 5 {
		if current.Value == "{" {
			this.Eat(5)
			this.JumpBreaks()
			object := map[string]*Token{}
			key := this.Eat(this.Current().Type)
			this.Eat(5)
			value := this.Parse()
			object[key.Value] = value
			for ;this.Current().Type == 5 && this.Current().Value == ",";{
				this.JumpBreaks()
				this.Eat(5)
				key = this.Eat(this.Current().Type)
				this.Eat(5)
				value = this.Parse()
				object[key.Value] = value
			}
			this.Eat(5)
			return this.CreateToken(nil,"map",object)
		}
		if current.Value == "(" {
			return this.Expr()
		}
		if current.Value == "[" {
			this.Eat(5)
			arr := []*Token{this.Parse()}
			for this.Current().Type == 5 && this.Current().Value == "," {
				this.Eat(5)
				arr = append(arr, this.Parse())
			}
			this.Eat(5)
			return this.CreateToken(nil, "array", arr)
		}
		if current.Value == "}" {
			this.Eat(5)
			return this.CreateToken(nil, "}", nil)
		}
	}
	if current.Type == 6 {
		keyword := this.Eat(6)
		if keyword.Value == "break" {
			return this.CreateToken("break", "break", nil)
		}
		if keyword.Value == "import"{
			val := this.Eat(3).Value
			return this.CreateToken(val,"import",nil)
		}
		if keyword.Value == "fn" {
			var name string
			if this.Current().Type == 7 {
				name = this.Eat(7).Value
			} else {
				name = "fn"
			}
			this.Eat(5)
			params := []string{}
			if this.Current().Type != 5 {
				params = append(params, this.Current().Value)
				this.Eat(7)
				for this.Current().Type == 5 && this.Current().Value == "," {
					this.Eat(5)
					params = append(params, this.Current().Value)
					this.Eat(7)
				}
			}
			this.Eat(5)
			this.Eat(5)
			fnTok := this.CreateToken(params, name, this.Compound())
			return fnTok
		}
		if keyword.Value == "if" {
			condition := this.Parse()
			this.Eat(5)
			tok := this.CreateToken(condition, "if", this.Compound())
			this.JumpBreaks()
			if this.Current().Type == 6 && this.Current().Value == "else" {
				this.Eat(6)
				this.Eat(5)
				elseToken := this.CreateToken(nil, "else", this.Compound())
				return this.CreateToken(nil, "conditions", []*Token{tok, elseToken})
			}
			return tok
		}
		if keyword.Value == "for" {
			var first *Token
			if this.Current().Type != 4 {
				first = this.Parse()
			}
			this.Eat(4)
			condition := this.Parse()
			this.Eat(4)
			var third *Token
			if this.Current().Type != 5 || this.Current().Value != "{" {
				third = this.Parse()
			}
			this.Eat(5)
			condToken := this.CreateToken(condition, "forif", this.Compound())
			forToken := this.CreateToken(first, "forToken", third)
			return this.CreateToken(condToken, "for", forToken)
		}
		if keyword.Value == "ret" {
			if this.Current().Type != 4 {
				return this.CreateToken(nil, "ret", this.Parse())
			} else {
				return this.CreateToken(nil, "ret", nil)
			}
		}
	}
	if current.Type == 2 || current.Type == 3 || current.Type == 7 {
		return this.Expr()
	} else {
		this.Eat(4)
	}
	return nil
}

func Parse(lexems []*Lexem) *Token {
	parser := &Parser{
		Lexems: lexems,
	}
	finalToken := parser.CreateToken(nil, "main", []*Token{})
	for parser.I < len(lexems) {
		if parser.Error != "" {
			return parser.CreateToken(parser.Error, "err", nil)
		} else {
			result := parser.Parse()
			if result != nil {
				finalToken.Right = append(finalToken.Right.([]*Token), result)
			}
		}
	}
	if parser.Error != "" {
		return parser.CreateToken(parser.Error, "err", nil)
	}
	return finalToken
}
