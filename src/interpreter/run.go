package interpreter


func Run(code string) interface{}{
	lexems := Lex(code)
	token := Parse(lexems)
	return Interpret(token)
}