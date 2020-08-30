package interpreter

import (
	"strconv"
	"reflect"
	"fmt"
)

func visit(intToken interface{}) interface{}{
	token := intToken.(*Token)
	if token.Right == nil && token.Left == nil{
		if token.Value[0] == '"'{
			return token.Value
		}
		i,err := strconv.Atoi(token.Value)
		if err != nil {
			return token.Value
		}
		return i
	}
	if token.Value == "main"{
		tokens := token.Right.([]*Token)
		for _,newToken := range tokens{
			visit(newToken)
		}
	}
	if token.Value == "call"{
		name := token.Left.(string)
		var values []interface{}
		for _,param := range token.Right.([]*Token){
			values = append(values,reflect.ValueOf(visit(param)))
		}
		if name == "print"{
			fmt.Println(values...)
		}
	}
	if token.Value == "+"{
		return visit(token.Left).(int) + visit(token.Right).(int)
	}
	if token.Value == "-"{
		return visit(token.Left).(int) - visit(token.Right).(int)
	}
	if token.Value == "*"{
		return visit(token.Left).(int) * visit(token.Right).(int)
	}
	if token.Value == "/"{
		return visit(token.Left).(int) / visit(token.Right).(int)
	}
	return nil
}
func Interpret(token *Token){
	visit(token)
}