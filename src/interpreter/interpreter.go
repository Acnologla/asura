package interpreter

import (
	"fmt"
	"log"
	"strconv"
)

var localVars = map[string]interface{}{}
func defaultValue(name string) interface{} {
	if name[0] == '"' {
		return name
	}
	i, err := strconv.Atoi(name)
	if err != nil {
		value, ok := localVars[name]
		if !ok {
			log.Fatal("Var not defined")
		}
		return value
	}
	return i
}


func visit(intToken interface{}) interface{} {
	token, ok := intToken.(*Token)	
	if !ok{
		str, ok2 := intToken.(string)
		if ok2{
			return defaultValue(str)
		}else{
			var ret interface{}
			tokens := intToken.([]*Token)
			for _, token := range tokens {
				ret = visit(token)
				tok, ok := ret.(*Token)
				if ok {
					if tok.Value == "ret" {
						ret = tok.Right
						break
					}
				}
			}
			return ret
		}
	}
	if token.Right == nil && token.Left == nil {
		return defaultValue(token.Value)
	}
	if token.Value == "ret" {
		return &Token{
			Value: "ret",
			Right: visit(token.Right),
		}
	}
	if token.Value == ":=" {
		name := token.Left.(string)
		localVars[name] = visit(token.Right)
	}
	if token.Value == "if"{
		condition := token.Left.(*Token)
		if condition.Value == "=="{
			if visit(condition.Left) == visit(condition.Right){
				return visit(token.Right)
			}
		}else if condition.Value == "!="{
			if visit(condition.Left) != visit(condition.Right){
				return visit(token.Right)
			}
		}else if condition.Value == ">"{
			if visit(condition.Left).(int) > visit(condition.Right).(int) {
				return visit(token.Right)
			}
		}else if condition.Value == "<"{
			if visit(condition.Left).(int) < visit(condition.Right).(int){
				return	visit(token.Right)
			}
		}
		return nil
	}
	if token.Value == "call" {
		name := token.Left.(string)
		var values []interface{}
		if token.Right != nil{
			for _, param := range token.Right.([]*Token) {
				values = append(values, visit(param))
			}
		}
		if name == "print" {
			fmt.Println(values...)
		}
		valInterface, ok := localVars[name]
		if ok {
			val := valInterface.(*Token)
			oldvars := map[string]interface{}{}
			for key, value := range localVars {
				oldvars[key] = value
			}
			for i, param := range val.Left.([]string) {
				localVars[param] = values[i]
			}
			var ret interface{}
			tokens := val.Right.([]*Token)
			for _, token := range tokens {
				ret = visit(token)
				tok, ok := ret.(*Token)
				if ok {
					if tok.Value == "ret" {
						ret = tok.Right
						break
					}
				}
			}
			localVars = oldvars
			return ret
		}
	}
	if token.Value == "+" {
		return visit(token.Left).(int) + visit(token.Right).(int)
	}
	if token.Value == "-" {
		return visit(token.Left).(int) - visit(token.Right).(int)
	}
	if token.Value == "*" {
		return visit(token.Left).(int) * visit(token.Right).(int)
	}
	if token.Value == "/" {
		return visit(token.Left).(int) / visit(token.Right).(int)
	}
	_, isCompound := token.Right.([]*Token)
	if isCompound{
		params, ok := token.Left.([]string)
		if ok {
			localVars[token.Value] = &Token{
				Left:  params,
				Right: token.Right,
			}
		}
	}
	return nil
}
func Interpret(token *Token) interface{}{
	return visit(token.Right)
}
