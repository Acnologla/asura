package interpreter

import (
	"strconv"
	"log"
	"asura/src/utils"
	"unicode"
)

var operators = []string{"+","-","=","==","*","/","(",")","**", ":=",":"}
var keywords = []string{"fn","if","else"}
var breakers = []string{";","\n"," "}

type Lexem struct{
	Type int
	Value string
}
/*

Type 1 = Bool
Type 2 = Int
Type 3 = String
Type 4 = Break
Type 5 = Operator
Type 6 = Keyword
Type 7 = Identifier

*/
func lex(code string,i int) (*Lexem, int){

	actual := string(code[i])
	if actual == " "{
		return nil, (i+1)
	}
	if utils.Includes(breakers,actual){
		return &Lexem{
			Type: 4,
		}, (i+1)
	}

	if actual == "\"" || actual == "'"{
		breaker := actual
		var str string
		i++
		for ;len(code) != i;i++{
			if string(code[i]) == breaker{
				break
			}
			str += string(code[i])
		}
		i++
		return &Lexem{
			Type: 2,
			Value: str,
		}, i 
	}
	 _, err := strconv.Atoi(actual)
	 if err == nil{
		var number = actual	 
		i++
		 for ;len(code) != i;i++{
			if utils.Includes(breakers,string(code[i])){
				_, err := strconv.Atoi(number)
				if err != nil{
					log.Fatal("Invalid number")
				}
				return &Lexem{
					Value: number,
					Type: 2,
				},i 
			}
			number += string(code[i])	 
		 }
	 }
	 if !unicode.IsNumber(rune(code[i])) && !unicode.IsLetter(rune(code[i])){
		var operator string
		for ;len(code) != i;i++{
			if unicode.IsNumber(rune(code[i])) || unicode.IsLetter(rune(code[i])) || utils.Includes(breakers,string(code[i])){
				break
			}
			operator += string(code[i])
		}
		if utils.Includes(operators,operator){
			return &Lexem{
				Type: 5,
				Value: operator, 
			},i
		}else{
			log.Fatal("Invalid operator")
			return nil,i
		}
	 }
	 var acc string
	 for ;len(code) != i;i++{
		if  (!unicode.IsNumber(rune(code[i])) && !unicode.IsLetter(rune(code[i]))) || utils.Includes(breakers,string(code[i])){
			break
		}
		acc += string(code[i])
	}
	if utils.Includes(keywords,acc){
		return &Lexem{
			Type: 6,
			Value: acc,
		},i
	}
	return &Lexem{
		Type: 7,
		Value: acc,
	},i
}

func Lex(code string) []*Lexem{
	var lexems = []*Lexem{}
	code = code + "\n"
	for i:=0; i < len(code)-1;{
		lexem,value := lex(code,i)
		i = value
		if lexem != nil{
			lexems = append(lexems,lexem)
		}
	}
	return lexems
}