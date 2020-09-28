package interpreter

import (
	"asura/src/utils"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

var localVars = []map[string]interface{}{}
var sf = 0

func toArrInterface(value interface{}) ([]interface{}, bool) {
	str, ok := value.(string)
	if ok {
		arr := make([]interface{}, len(str))
		for i := 0; i < len(str); i++ {
			arr[i] = string(str[i])
		}
		return arr, true
	}
	s := reflect.ValueOf(value)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Slice {
		print("InterfaceSlice() given a non-slice type")
		return []interface{}{}, false
	}
	arr := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		arr[i] = s.Index(i).Interface()
	}
	return arr, true
}

func updateSf(name string,val interface{}){
	for i := len(localVars)-1; 0 <= i;i--{
		_,exists := localVars[i][name]
		if exists{
			localVars[i][name] = val
			break
		}
	}
} 

func getSf(name string) (interface{},bool){
	for i := len(localVars)-1; 0 <= i;i--{
		val,exists := localVars[i][name]
		if exists{
			return val, true
		}
	}
	return nil, false
}

func defaultValue(name string) interface{} {
	if name[0] == '"' {
		name := name[:len(name)-1]
		return name[1:]
	}
	i, err := strconv.ParseFloat(name, 64)
	if err != nil {
		value, ok := getSf(name)
		if !ok {
			value, ok := defaultVars[name]
			if ok {
				return value
			} else {
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
	if token.Value == "value" {
		return token.Left
	}
	if token.Value == "break" {
		return token
	}
	if token.Value == "index" {
		right, ok := visit(token.Right).(float64)
		if !ok {
			_, ok := visit(token.Right).(string)
			if ok {
				return visit(&Token{
					Value: "property",
					Left:  token.Left,
					Right: token.Right,
				})
			}
		}
		left := visit(token.Left)
		str, ok := left.(string)
		if !ok {
			arr, _ := toArrInterface(left)
			if int(right) >= len(arr) {
				return nil
			}
			return arr[int(right)]
		}
		if int(right) >= len(str) {
			return nil
		}
		return string(str[int(right)])
	}
	if token.Value == "map" {
		object := map[string]interface{}{}
		for key, value := range token.Right.(map[string]*Token) {
			object[key] = visit(value)
		}
		return object
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
	if token.Value == "="{
		val,isTokn := token.Left.(*Token)
		if isTokn{
			if val.Value == "property"{
				return visit(&Token{
					Left: val.Left,
					Value: "property",
					Right: &Token{
						Left: val.Right.(*Token).Value,
						Value: "changeObj",
						Right: token.Right,
					},
				})
			}
		}
		name := token.Left.(string)
		updateSf(name,visit(token.Right))
	}
	if token.Value == ":=" {
		name := token.Left.(string)
		localVars[sf][name] = visit(token.Right)
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
	if token.Value == "==" {
		if visit(token.Left) == visit(token.Right) {
			return true
		}
		return false
	}
	if token.Value == "!=" {
		if visit(token.Left) != visit(token.Right) {
			return true
		}
		return false
	}
	if token.Value == ">" {
		if visit(token.Left).(float64) > visit(token.Right).(float64) {
			return true
		}
		return false
	}
	if token.Value == "<" {
		if visit(token.Left).(float64) < visit(token.Right).(float64) {
			return true
		}
		return false
	}
	if token.Value == "if" {
		condition, ok := visit(token.Left).(bool)
		if !ok {
			return false
		}
		if condition {
			return visit(token.Right)
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
		objValue, isObj := visit(arr.Left).(map[string]interface{})
		if isObj{
			result := visit(arr.Right).(string)
			objValue[result] = visit(token.Right)
		}
	}
	if token.Value == "fn" {
		return token
	}
	if token.Value == "property" {
		propertyName := token.Right.(*Token)
		left := visit(token.Left)
		r := reflect.ValueOf(left)
		if !r.IsValid() {
			return nil
		}
		var acess int = -1
		var recursive Token
		if propertyName.Value == "index" {
			_, ok := visit(propertyName.Right).(string)
			if ok {
				propertyName.Value = "property"
				acess = -2
			}
		}
		change := false
		var old *Token
		if propertyName.Value == "changeObj"{
			change = true
			old = &Token{
				Right: propertyName.Right,
			}
			propertyName.Value = propertyName.Left.(string)
			propertyName.Left = nil
			propertyName.Right = nil 
		}
		if propertyName.Value == "property" {
			recursive = Token{
				Value: "property",
				Left:  propertyName.Left,
				Right: propertyName.Right,
			}
			str, ok := propertyName.Left.(string)
			var isTkn *Token
			if !ok {
				newval := visit(propertyName.Left)
				if newval == nil {
					return nil
				}
				isTkn = propertyName.Left.(*Token)
				if isTkn.Value != "call" {
					isTkn = nil
				}
				str = newval.(string)
			}
			if isTkn != nil {
				propertyName.Value = "call"
				propertyName.Left = str
				propertyName.Right = isTkn.Right
			} else {
				propertyName.Value = str
				propertyName.Left = nil
				propertyName.Right = nil
			}
		}
		if propertyName.Value == "index" && acess == -1 {
			acess = int(visit(propertyName.Right).(float64))
			propertyName = propertyName.Left.(*Token)
		}
		if propertyName.Right == nil && propertyName.Left == nil {
			if left == nil {
				return nil
			}
			if strings.HasPrefix(propertyName.Value, "\"") && strings.HasSuffix(propertyName.Value, "\"") {
				propertyName.Value = propertyName.Value[1 : len(propertyName.Value)-1]
			}
			if reflect.Indirect(r).Type().Kind() == reflect.Map {
				for _, e := range reflect.Indirect(r).MapKeys() {
					if e.Interface() == reflect.ValueOf(propertyName.Value).Interface() {
						if change{
							reflect.Indirect(r).SetMapIndex(e,reflect.ValueOf(visit(old.Right)))
							return nil
						}
						return reflect.Indirect(r).MapIndex(e).Interface()
					}
					return nil
				}
			}
			field := reflect.Indirect(r).FieldByName(propertyName.Value)
			if !field.IsValid() {
				return nil
			}
			if field.Type().Name() == "int" {
				field = reflect.ValueOf(float64(field.Interface().(int)))
			}
			if field.Type().Name() == "uint" {
				field = reflect.ValueOf(float64(field.Interface().(uint)))
			}
			if recursive.Value == "property" {
				return visit(&Token{
					Left: &Token{
						Left:  field.Interface(),
						Value: "value",
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
				var fun reflect.Value
				if r.Type().Kind() == reflect.Map{
					for _, e := range reflect.Indirect(r).MapKeys() {
						if e.Interface() == reflect.ValueOf(name).Interface() {
							fun = r.MapIndex(e)
						}
					}
				}else{
					fun = r.MethodByName(name)
				}
				if fun.IsValid() {
					fnToken,isToken := fun.Interface().(*Token)
					var result =  []reflect.Value{}
					if isToken{
						result = append(result,reflect.ValueOf(visit(&Token{
							Value: "call",
							Left: fnToken,
							Right: propertyName.Right,
						})))
					}else{
						var values []reflect.Value
					if propertyName.Right != nil {
						for i, param := range propertyName.Right.([]*Token) {
							paramType := fun.Type().In(i).Name()
							val := reflect.ValueOf(visit(param))
							valtType := val.Type().Name()
							if paramType == "int" && valtType == "float64" {
								val = reflect.ValueOf(int(val.Interface().(float64)))
							} else if paramType == "Snowflake" && valtType == "string" {
								val = reflect.ValueOf(utils.StringToID(val.Interface().(string)))
							} else if paramType != valtType && paramType != "Context" && paramType != "Session" {
								fmt.Printf("Parameter error %s is not %s\n", paramType, valtType)
								return nil
							}
							values = append(values, val)
						}
					}
						result = fun.Call(values)
					}
					var results []interface{}
					for _, val := range result {
						if val.Type().Name() == "int" {
							results = append(results, reflect.ValueOf(float64(val.Interface().(int))).Interface())
						} else if val.Type().Name() == "uint" {
							results = append(results, reflect.ValueOf(float64(val.Interface().(uint))).Interface())
						} else {
							results = append(results, val.Interface())
						}
					}
					if len(results) == 0 {
						return nil
					}
					if acess >= 0 {
						resultsArr, isArr := toArrInterface(results[0])
						if !isArr {
							return nil
						}
						if acess >= len(resultsArr) {
							return nil
						}
						if recursive.Value == "property" {
							return visit(&Token{
								Left: &Token{
									Left:  resultsArr[acess],
									Value: "value",
								},
								Value: "property",
								Right: recursive.Right,
							})
						}
						return resultsArr[acess]
					}
					if recursive.Value == "property" {
						return visit(&Token{
							Left: &Token{
								Left:  results[0],
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
	if token.Value == "import" {
		importPath := token.Left.(string)
		content, err := ioutil.ReadFile(fmt.Sprintf("./src/interpreter/libs/%s.acnl", importPath))
		if err != nil {
			fmt.Println("Invalid file")
			return nil
		}
		lexems := Lex(string(content))
		token := Parse(lexems)
		visit(token.Right)
		return nil
	}
	if token.Value == "call" {
		var valFn *Token
		name, ok := token.Left.(string)
		if !ok {
			v := visit(token.Left)
			val2, ok := v.(*Token)
			if ok {
				if val2.Value == "fn" {
					valFn = val2
				}
			}
		}
		var values []interface{}
		if token.Right != nil {
			for _, param := range token.Right.([]*Token) {
				values = append(values, visit(param))
			}
		}
		fn, ok := defaultVars[name]
		if ok {
			newVals := []reflect.Value{}
			for _, val := range values {
				if val != nil {
					newVals = append(newVals, reflect.ValueOf(val))
				}
			}
			returned := reflect.ValueOf(fn).Call(newVals)
			if reflect.ValueOf(fn).Type().NumOut() == 1 {
				return returned[0].Interface()
			} else {
				vals := []interface{}{}
				for _, val := range returned {
					vals = append(vals, val.Interface())
				}
				return vals
			}
		}
		valInterface, ok := getSf(name)
		if valFn != nil && !ok {
			valInterface = valFn
			ok = true
		}
		if ok {
			val := valInterface.(*Token)
			sf++
			localVars = append(localVars,map[string]interface{}{})
			for i, param := range val.Left.([]string) {
				if i < len(values) {
					localVars[sf][param] = values[i]
				} else {
					localVars[sf][param] = nil
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
					for key,val := range localVars[sf]{
						localVars[0][key] = val
					}
					for key,val := range localVars[sf-1]{
						localVars[0][key] = val
					}
				}
			}
			sf--
			localVars = localVars[:len(localVars)-1]
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
				if strings.HasSuffix(first.(string), ".000") {
					first = first.(string)[:len(first.(string))-4]
				}
			}
			if ok2 {
				second = fmt.Sprintf("%.3f", second.(float64))
				if strings.HasSuffix(second.(string), ".000") {
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
	if token.Value == "%"{
		return int(visit(token.Left).(float64)) % int(visit(token.Right).(float64))
	}
	_, isCompound := token.Right.([]*Token)
	if isCompound {
		params, ok := token.Left.([]string)
		if ok {
			localVars[sf][token.Value] = &Token{
				Left:  params,
				Right: token.Right,
				Value: "fn",
			}
		}
	}
	return nil
}

func Interpret(token *Token, params map[string]interface{}) interface{} {
	if len(localVars) == 0{
		localVars = append(localVars,map[string]interface{}{})
	}
	for key, val := range params {
		localVars[0][key] = val
	}
	//debugToken(token)
	return visit(token.Right)
}

func debugToken(token interface{}) {
	compound, isCompound := token.([]*Token)
	if isCompound {
		fmt.Println("[")
		for _, tkn := range compound {
			debugToken(tkn)
			fmt.Println(",\n")
		}
		fmt.Println("]")
		return
	}
	tkn, isToken := token.(*Token)
	if isToken  && tkn != nil{
		if tkn.Left == nil && tkn.Right == nil {
			fmt.Println(tkn.Value)
			return
		}
		fmt.Printf("{%s:{\nLeft:", tkn.Value)
		debugToken(tkn.Left)
		fmt.Print(",\nRight:")
		debugToken(tkn.Right)
		fmt.Println("}\n}")
	} else {
		fmt.Println(token)
	}
}
