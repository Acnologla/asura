package handler

import (
	"strings"
	"strconv"
)
type Argument interface { 
	Num() NumberArgument
	Str() StringArgument
	Mention() MentionArgument
}

type NumberArgument struct {
	min int
	max int
}

type StringArgument struct {
	min int
	max int
}

type MentionArgument struct {
	sameServer bool 
}

func (arg *NumberArgument) Validate(str string) bool {
	integer, err := strconv.Atoi(str)
	return err != nil && integer >= arg.min && integer <= arg.max
}

func (arg *StringArgument) Validate(str string) bool {
	return len(str) >= arg.min && len(str) <= arg.max
}

func (arg *MentionArgument) Validate(str string) bool {
	return strings.HasPrefix(str,"<@") &&  strings.HasPrefix(str,">")
}
