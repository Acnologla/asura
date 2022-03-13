package test

import (
	"testing"
)

func TestTrueIsTrue(t *testing.T) {
	if true == false {
		t.Error("True is not true")
	}
}
