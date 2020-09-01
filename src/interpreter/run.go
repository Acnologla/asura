package interpreter

import "fmt"

func Run(code string, params map[string]interface{}) interface{} {
	lexems := Lex(code)
	token := Parse(lexems)
	if token.Value == "err"{
		fmt.Println(token.Left.(string))
		return token.Left.(string)
	}
	return Interpret(token,params)
}
