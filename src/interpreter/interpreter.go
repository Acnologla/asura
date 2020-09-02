package interpreter

import (
	"fmt"
	"strings"
	"reflect"
	"strconv"
	"asura/src/utils"
)

var localVars = map[string]interface{}{}


func defaultValue(name string) interface{} {
	if name[0] == '"' {
		name := name[:len(name)-1]
		return name[1:]
	}
	i, err := strconv.ParseFloat(name, 64)
	if err != nil {
		value, ok := localVars[name]
		if !ok {
			value, ok := defaultVars[name]
			if ok{
				return value
			}else{
				fmt.Printf("Var %s not defined\n", name)
				return nil
			}
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
					if tok.Value == "break" {
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
		return nil
	}
	if token == nil {
		return nil
	}
	if token.Right == nil && token.Left == nil {
		return defaultValue(token.Value)
	}
	if token.Value == "value"{
		return token.Left
	}
	if token.Value == "break" {
		return token
	}
	if token.Value == "index" {
		left := visit(token.Left)
		right := visit(token.Right).(float64)
		str, ok := left.(string)
		if !ok {
			arr := left.([]interface{})
			if int(right) >= len(arr){
				return nil
			}
			return arr[int(right)]
		}
		if int(right) >= len(str){
			return nil
		}
		return string(str[int(right)])
	}
	if token.Value == "array" {
		arr := []interface{}{}
		for _, item := range token.Right.([]*Token) {
			arr = append(arr, visit(item))
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
	if token.Value == "for" {
		forToken := token.Right.(*Token)
		visit(forToken.Left)
		var lastVal interface{}
		for {
			val := visit(token.Left)
			returned, ok := val.(*Token)
			if ok {
				if returned.Value == "forif" {
					bo, ok := returned.Left.(bool)
					if ok {
						if !bo {
							return lastVal
						} else {
							val = returned.Right
						}
					} else {
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
				if !ok {
					return val
				} else {
					if returned {
						return val
					}
				}
			}
		}
	}
	if token.Value == "forif" {
		result := visit(&Token{
			Value: "if",
			Left:  token.Left,
			Right: []*Token{},
		})
		boo, ok := result.(bool)
		var ret interface{}
		if !ok {
			ret = visit(token.Right)
		}
		if ok {
			if boo {
				ret = visit(token.Right)
			}
		}
		if result == nil {
			boo = true
		}
		return &Token{
			Value: "forif",
			Right: ret,
			Left:  boo,
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
	if token.Value == "changeArr" {
		arr := token.Left.(*Token)
		arrValue, isArr := visit(arr.Left).([]interface{})
		if isArr {
			result := visit(arr.Right).(float64)
			arrValue[int(result)] = visit(token.Right)
		}
	}
	if token.Value == "fn" {
		return token
	}
	if token.Value == "property" {
		propertyName := token.Right.(*Token)
		left := visit(token.Left)
		r := reflect.ValueOf(left)
		if !r.IsValid(){
			return nil
		}
		var acess int = -1
		var recursive Token
		if propertyName.Value == "property"{
			recursive = Token{
				Value: "property",
				Left: propertyName.Left,
				Right: propertyName.Right,
			}
			str,ok := propertyName.Left.(string)
			var isTkn *Token
			if !ok{
				newval := visit(propertyName.Left)
				if newval == nil{
					return nil
				}
				isTkn = propertyName.Left.(*Token)
				if isTkn.Value != "call"{
					isTkn = nil
				}
				str = newval.(string)
			}
			if  isTkn != nil{
				propertyName.Value = "call"
				propertyName.Left = str
				propertyName.Right = isTkn.Right
			}else{
				propertyName.Value = str
				propertyName.Left = nil
				propertyName.Right = nil
			}
		}
		if propertyName.Value == "index"{
			acess = int(visit(propertyName.Right).(float64))
			propertyName = propertyName.Left.(*Token)
		}
		if propertyName.Right == nil && propertyName.Left == nil {
			if left == nil{
				return nil
			}
			field := reflect.Indirect(r).FieldByName(propertyName.Value)
			if !field.IsValid(){
				return nil
			}
			if recursive.Value == "property"{
				return visit(&Token{
					Left: &Token{
						Left: field.Interface(),
						Value : "value",
					},
					Value: "property",
					Right: recursive.Right,
				})
			}
			return field.Interface()
		}
		if propertyName.Value == "call" {
			name, ok := propertyName.Left.(string)
			if ok {
				fun := r.MethodByName(name)
				if fun.IsValid(){
					var values []reflect.Value
					if propertyName.Right != nil {
						for i, param := range propertyName.Right.([]*Token) {
							paramType := fun.Type().In(i).Name()
							val :=  reflect.ValueOf(visit(param))
							valtType := val.Type().Name()
							if paramType == "int" && valtType == "float64"{
								val = reflect.ValueOf(int(val.Interface().(float64)))
							}else if paramType == "Snowflake" && valtType == "string"{
								val = reflect.ValueOf(utils.StringToID(val.Interface().(string)))
							}else if paramType != valtType && paramType != "Context" && paramType != "Session"{
								fmt.Printf("Parameter error %s is not %s\n",paramType,valtType)
								return nil
							}
							values = append(values,val)
						}
					}
					result := fun.Call(values)
					var results []interface{}
					for _,val := range result{
						results = append(results,val.Interface())
					}
					if acess >= 0 {
						if acess >= len(results){
							return  nil
						}
						if recursive.Value == "property"{
							return visit(&Token{
								Left: &Token{
									Left: results[acess],
									Value: "value",
								},
								Value: "property",
								Right: recursive.Right,
							})
						}
						return results[acess]
					}
					if recursive.Value == "property"{
						return visit(&Token{
							Left: &Token{
								Left: results[0],
								Value: "value",
							},
							Value: "property",
							Right: recursive.Right,
						})
					}
					return results[0]
				}
			}
			return nil
		}
	}
	if token.Value == "call" {
		name := token.Left.(string)
		var values []interface{}
		if token.Right != nil {
			for _, param := range token.Right.([]*Token) {
				values = append(values, visit(param))
			}
		}
		fn, ok := defaultVars[name]
		if ok {
			newVals := []reflect.Value{}
			for _,val := range values{
				if val != nil{
					newVals = append(newVals, reflect.ValueOf(val))
				}
			}
			 returned := reflect.ValueOf(fn).Call(newVals)
			 if reflect.ValueOf(fn).Type().NumOut() == 1{
				 return returned[0].Interface()
			 }else{
				 vals := []interface{}{}
				 for _,val := range returned{
					 vals = append(vals,val.Interface())
				 }
				 return vals
			 }
		}
		valInterface, ok := localVars[name]
		if ok {
			val := valInterface.(*Token)
			oldvars := map[string]interface{}{}
			for key, value := range localVars {
				oldvars[key] = value
			}
			for i, param := range val.Left.([]string) {
				if  i < len(values){
					localVars[param] = values[i]
				} else{
					localVars[param] = nil
				}
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
			retToken, ok := ret.(*Token)
			if ok {
				if retToken.Value == "fn" {
					return ret
				}
			}
			for key, _ := range localVars {
				_, exists := oldvars[key]
				if !exists {
					delete(localVars, key)
				}
			}
			return ret
		} else {
			return name
		}
	}
	if token.Value == "+" {
		first := visit(token.Left)
		second := visit(token.Right)
		int1, ok := first.(float64)
		int2, ok2 := second.(float64)
		if !ok || !ok2 {
			if ok {
				first = fmt.Sprintf("%.3f", first.(float64))
				if strings.HasSuffix(first.(string),".000"){
					first = first.(string)[:len(first.(string))-4]
				}
			}
			if ok2 {
				second = fmt.Sprintf("%.3f", second.(float64))
				if strings.HasSuffix(second.(string),".000"){
					second = second.(string)[:len(second.(string))-4]
				}
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
				Value: "fn",
			}
		}
	}
	return nil
}
func Interpret(token *Token, params map[string]interface{}) interface{} {
	for key,val := range params{
		localVars[key] = val
	}
	//debugToken(token)
	return visit(token.Right)
}

func debugToken(token interface{}){	
	compound,isCompound := token.([]*Token)
	if isCompound{
		fmt.Println("[")
		for _,tkn := range compound{
			debugToken(tkn)
			fmt.Println(",\n")
		}
		fmt.Println("]")
		return 
	}
	tkn,isToken := token.(*Token)
	if isToken{
		if tkn.Left == nil && tkn.Right == nil{
			fmt.Println(tkn.Value)
			return
		}
		fmt.Printf("{%s:{\nLeft:",tkn.Value)
		debugToken(tkn.Left)
		fmt.Print(",\nRight:")
		debugToken(tkn.Right)
		fmt.Println("}\n}")
	}else{
		fmt.Println(token)
	}
}