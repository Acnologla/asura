package interpreter

import (
	"asura/src/utils"
	"log"
)

type Token struct {
	Left  interface{}
	Value string
	Right interface{}
}

type Parser struct {
	Lexems []*Lexem
	Error  bool
	I      int
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
		return &Lexem{}
	}
	if this.Lexems[this.I].Type != lexemType {
		log.Fatalf("Invalid type expected: %d, got %d\n Value: %s", lexemType, this.Lexems[this.I].Type, this.Lexems[this.I].Value)
	} else {
		this.I++
		return this.Lexems[this.I-1]
	}
	return nil
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
		this.Eat(7)
		if this.Current().Type == 5 && this.Current().Value == "."{
			this.Eat(5)
			return this.CreateToken(current.Value,"property",this.Parse())
		}
	} else if current.Type == 3 {
		this.Eat(3)
		return this.CreateToken(nil, "\""+current.Value+"\"", nil)
	} else {
		this.Eat(2)
	}
	if this.Current().Type == 5 && this.Current().Value == "[" {
		this.Eat(5)
		tok := this.Parse()
		newTok := this.CreateToken(current.Value, "index", tok)
		this.Eat(5)
		this.JumpBreaks()
		if this.Current().Type == 5 && this.Current().Value == "="{
			this.Eat(5)
			value := this.Parse()
			return this.CreateToken(newTok,"changeArr",value)
		}
		return newTok
	}
	return this.CreateToken(nil, current.Value, nil)
}

func (this *Parser) Mult() *Token {
	old := this.Factor()
	for this.Current().Type == 5 && (this.Current().Value == "*" || this.Current().Value == "/") {
		op := this.Eat(5)
		old = this.CreateToken(old, op.Value, this.Factor())
	}
	return old
}

func (this *Parser) Expr() *Token {
	old := this.Mult()
	for this.Current().Type == 5 && (this.Current().Value == "+" || this.Current().Value == "-") {
		op := this.Eat(5)
		old = this.CreateToken(old, op.Value, this.Mult())
	}
	if this.Current().Type == 5 && utils.Includes(comparators, this.Current().Value) {
		oldVal := this.Eat(5)
		return this.CreateToken(old, oldVal.Value, this.Parse())
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

func (this *Parser) JumpBreaks(){
	for this.Current().Type == 4 {
		this.Eat(4)
	}
}

func (this *Parser) Parse() *Token {
	current := this.Current()
	if current.Type == 5 {
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
		if keyword.Value == "fn" {
			var name string
			if this.Current().Type == 7{
				name = this.Eat(7).Value
			}else{
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
			first := this.Parse()
			this.Eat(4)
			condition := this.Parse()
			this.Eat(4)
			third := this.Parse()
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
	if current.Type == 7 {
		name := this.Eat(7)
		if this.Current().Type == 5 && this.Current().Value != ")" {
			operator := this.Eat(5)
			if operator.Value == "++" {
				return this.CreateToken(name.Value, ":=", this.CreateToken(name.Value, "+", this.CreateToken(nil, "1", nil)))
			}
			if operator.Value == "--" {
				return this.CreateToken(name.Value, ":=", this.CreateToken(name.Value, "-", this.CreateToken(nil, "1", nil)))
			}
			if operator.Value == ":=" {
				return this.CreateToken(name.Value, ":=", this.Parse())
			}
			if operator.Value == "(" {
				var tok *Token
				if this.Current().Type == 5 && this.Current().Value == ")" {
					tok = this.CreateToken(name.Value, "call", nil)
				} else {
					params := []*Token{this.Parse()}
					for this.Current().Type == 5 && this.Current().Value == "," {
						this.Eat(5)
						params = append(params, this.Parse())
					}
					tok = this.CreateToken(name.Value, "call", params)
				}
				this.Eat(5)
				if this.Current().Type == 5 && this.Current().Value == "."{
					this.Eat(5)
					return this.CreateToken(tok,"property",this.Parse())
				}
				return tok
			}
			this.I -= 2
			return this.Expr()
		}
		return this.CreateToken(nil, name.Value, nil)
	}
	if current.Type == 2 || current.Type == 3 {
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
		if parser.Error {
			return nil
		} else {
			result := parser.Parse()
			if result != nil {
				finalToken.Right = append(finalToken.Right.([]*Token), result)
			}
		}
	}
	return finalToken
}
