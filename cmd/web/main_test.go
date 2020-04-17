package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestEditStateType(t *testing.T) {
	// var err error
	var s StatesType
	var buff bytes.Buffer

	s.editStateType(&buff)
	s.errorLog.Printf("foo")
	s.infoLog.Printf("bar")
	content := fmt.Sprintf("%s", &buff)

	passCond := len(s.State) == 50 && len(s.Short) == 50 &&
		strings.Contains(content, "foo") &&
		strings.Contains(content, "bar")

	if !passCond {
		t.Errorf("expected 50 for states, got %d for States and %d for Short",
			len(s.State), len(s.Short))
	}
}
