package test

import (
	"github.com/acnologla/interpreter"
	"testing"
)

func TestLexer(t *testing.T) {
	code := "a := 2"
	lexems := interpreter.Lex(code)
	if lexems[0].Value != "a" || lexems[0].Type != 7 {
		t.Errorf("Invalid lexem %d", 0)
	}
	if lexems[1].Value != ":=" || lexems[1].Type != 5 {
		t.Errorf("Invalid lexem %d", 1)
	}
	if lexems[2].Value != "2" || lexems[2].Type != 2 {
		t.Errorf("Invalid lexem %d", 2)
	}
	t.Run("TestParser", func(t *testing.T) {
		token := interpreter.Parse(lexems)
		tokens := token.Right.([]*interpreter.Token)
		currentToken := tokens[0]
		if currentToken.Left.(string) != "a" {
			t.Errorf("Invalid token %d", 0)
		}
		if currentToken.Value != ":=" {
			t.Errorf("Invalid token %d", 1)
		}
		if currentToken.Right.(*interpreter.Token).Value != "2" {
			t.Errorf("Invalid token %d", 2)
		}
		t.Run("TestInterpreter", func(t *testing.T) {
			result := interpreter.Interpret(token, map[string]interface{}{})
			if result != nil {
				t.Errorf("Interpreter should return nil")
			}
		})
	})
}
