package interpreter

import (
	"fmt"
	"log"
	"strconv"
)

var localVars = map[string]interface{}{}

func defaultValue(name string) interface{} {
	if name[0] == '"' {
		name := name[:len(name)-1]
		return name[1:]
	}
	i, err := strconv.ParseFloat(name,64)
	if err != nil {
		value, ok := localVars[name]
		if !ok {
			log.Fatalf("Var %s not defined",name)
		}
		return value
	}
	return i
}

func visit(intToken interface{}) interface{} {
	token, ok := intToken.(*Token)
	if !ok {
		str, ok2 := intToken.(string)
		if ok2 {
			return defaultValue(str)
		} else {
			var ret interface{}
			var lastVal interface{}
			tokens := intToken.([]*Token)
			for _, token := range tokens {
				ret = visit(token)
				tok, ok := ret.(*Token)
				if ok {
					if tok.Value == "break"{
						return lastVal
					}
					if tok.Value == "ret" {
						return ret
					}
				}
				lastVal = ret
			}
			return ret
		}
	}
	if token.Right == nil && token.Left == nil {
		return defaultValue(token.Value)
	}
	if token.Value == "break"{
		return token
	}
	if token.Value == "index"{
		left := visit(token.Left)
		right := visit(token.Right).(float64)
		str,ok := left.(string)
		if !ok{
			arr := left.([]interface{})
			return arr[int(right)]
		}
		return string(str[int(right)])
	}
	if token.Value == "array"{
		arr := []interface{}{}
		for _,item := range token.Right.([]*Token){
			arr = append(arr,visit(item))
		}
		return arr
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
	if token.Value == "for"{
		forToken := token.Right.(*Token)
		visit(forToken.Left)
		var lastVal interface{}
		for {
			val := visit(token.Left)
			returned, ok := val.(*Token)
			if ok{
				if returned.Value == "forif"{
				bo,ok := returned.Left.(bool)
				if ok{
					if !bo{
						return lastVal
					}else{
						val = returned.Right
					}
				}else{
					val = returned.Right
				}
			}
			}
			visit(forToken.Right)
			lastVal = val
		}
	}
	if token.Value == "conditions" {
		for _, condition := range token.Right.([]*Token) {
			if condition.Value == "else" {
				return visit(condition.Right)
			} else {
				val := visit(condition)
				returned, ok := val.(bool)
				if !ok{
					return val
				}else{
					if returned{
						return val
					}
				}
			}
		}
	}
	if token.Value == "forif"{
		result := visit(&Token{
			Value: "if",
			Left: token.Left,
			Right: []*Token{},
		})
		boo, ok := result.(bool)
		var ret interface{}
		if !ok {
			ret = visit(token.Right)
		}
		if ok{
			if boo{
				ret = visit(token.Right)
			}	
		}
		if result == nil{
			boo = true
		}
		return &Token{
			Value: "forif",
			Right: ret,
			Left: boo,
		}
	}
	if token.Value == "if" {
		condition := token.Left.(*Token)
		if condition.Value == "==" {
			if visit(condition.Left) == visit(condition.Right) {
				return visit(token.Right)
			}
		} else if condition.Value == "!=" {
			if visit(condition.Left) != visit(condition.Right) {
				return visit(token.Right)
			}
		} else if condition.Value == ">" {
			if visit(condition.Left).(float64) > visit(condition.Right).(float64) {
				return visit(token.Right)
			}
		} else if condition.Value == "<" {
			if visit(condition.Left).(float64) < visit(condition.Right).(float64) {
				return visit(token.Right)
			}
		}
		return false
	}
	if token.Value == "call" {
		name := token.Left.(string)
		var values []interface{}
		if token.Right != nil {
			for _, param := range token.Right.([]*Token) {
				values = append(values, visit(param))
			}
		}
		if name == "print" {
			fmt.Println(values...)
			return nil
		}
		if name == "len"{
			arr := values[0].([]interface{})
			return float64(len(arr))
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
		first := visit(token.Left)
		second := visit(token.Right)
		int1,ok := first.(float64)
		int2,ok2 := second.(float64)
		if !ok || !ok2{
			if ok{ 
				first = fmt.Sprintf("%f",first.(float64))
			}
			if ok2{
				second = fmt.Sprintf("%f",second.(float64))
			}
			return first.(string) + second.(string)
		}
		return int1 + int2
	}
	if token.Value == "-" {
		return visit(token.Left).(float64) - visit(token.Right).(float64)
	}
	if token.Value == "*" {
		return visit(token.Left).(float64) * visit(token.Right).(float64)
	}
	if token.Value == "/" {
		return visit(token.Left).(float64) / visit(token.Right).(float64)
	}
	_, isCompound := token.Right.([]*Token)
	if isCompound {
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
func Interpret(token *Token) interface{} {
	return visit(token.Right)
}
