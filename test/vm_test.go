package test

import (
	"testing"
	"asura/src/vm"
)

func TestLexer(t *testing.T) {
	code := "a := 2"
	lexems := vm.Lex(code)
	if lexems[0].Value != "a" || lexems[0].Type != 7{
		t.Errorf("Invalid lexem %d",0)
	}
	if lexems[1].Value != ":=" || lexems[1].Type != 5{
		t.Errorf("Invalid lexem %d",1)
	}
	if lexems[2].Value != "2" || lexems[2].Type != 2{
		t.Errorf("Invalid lexem %d",2)
	}
}
