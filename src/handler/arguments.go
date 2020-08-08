package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ValidateNumber(str string, min int, max int) error {
	integer, err := strconv.Atoi(str)
	if err != nil {
		return errors.New(fmt.Sprintf("Expected a number but instead got %s", str))
	}
	if integer < min || integer > max {
		return errors.New(fmt.Sprintf("Expected %s to be a number between %d and %d", str, min, max))
	}
	return nil
}

func ValidateString(str string, min int, max int) error {
	if len(str) < min || len(str) > max {
		return errors.New(fmt.Sprintf("Expected %s to be a text with size between %d and %d", str, min, max))
	}
	return nil
}

func ValidateUserMention(str string) error {
	if !strings.HasPrefix(str, "<@!") || !strings.HasPrefix(str, ">") {
		return errors.New(fmt.Sprintf("Expected a user mention but instead got %s", str))
	}
	return nil
}
