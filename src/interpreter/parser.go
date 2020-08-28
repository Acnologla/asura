package interpreter

import (
	"log"
)

type Token struct{
	Left interface{}
	Value  string
	Right interface{}
}

type Parser struct{
	Lexems []*Lexem
	Error bool
	I int
}

func (this *Parser) CreateToken(left interface{},value string,right interface{}) * Token{
	return &Token{
		Left: left,
		Value: value,
		Right: right,
	}
}

func (this *Parser) Current() *Lexem{
	if len(this.Lexems) == this.I{
		return &Lexem{}
	}
	return this.Lexems[this.I]
}

func (this *Parser) Eat(lexemType int) *Lexem{
	if len(this.Lexems) == this.I{
		return &Lexem{}
	}
	if this.Lexems[this.I].Type != lexemType{
		log.Fatal("Invalid type")
	}else{
		this.I++
		return this.Lexems[this.I-1]
	}
	return nil
}


func (this *Parser) Factor() *Token{
	current := this.Current()
	if current.Type == 5 && current.Value == "("{
		expr := this.Expr()
		this.Eat(5)
		return expr
	}
	this.Eat(2)
	return this.CreateToken(nil,current.Value,nil)
}

func (this *Parser) Mult()  *Token{
	old := this.Factor()
	for ;this.Current().Type == 5 && (this.Current().Value == "*" || this.Current().Value == "/");{
		op := this.Eat(5)
		old = this.CreateToken(old,op.Value,this.Factor())
	}
	return old
}

func (this *Parser) Expr() *Token{
	old := this.Mult()
	for ;this.Current().Type == 5 && (this.Current().Value == "+" || this.Current().Value == "-");{
		op:= this.Eat(5)
		old =  this.CreateToken(old,op.Value,this.Mult())
	}
	return old
}

func (this *Parser) Parse() *Token{
	current := this.Current()
	if current.Type == 7{
		name := this.Eat(7)
		if this.Current().Type == 5{
			operator := this.Eat(5)
			if  operator.Value == ":="{
				return this.CreateToken(name.Value,":=",this.Expr())
			}
		}
	}
	if current.Type == 2 || current.Type == 3{
		return this.Expr()
	}
	return nil
}



func Parse(lexems []*Lexem) *Token{
	parser := &Parser{
		Lexems: lexems,
	}
	finalToken := parser.CreateToken(nil,"main",[]*Token{})
	for ;parser.I < len(lexems);{
		if parser.Error{
			return nil
		}else{
			finalToken.Right = append(finalToken.Right.([]*Token),parser.Parse())
		}
	}
	return finalToken
}