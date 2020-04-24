package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestEditStateType(t *testing.T) {
	// var err error
	var st sT
	var buff bytes.Buffer
	st.buildST(&buff)
	s := st.copyScreenData()

	st.errorLog.Printf("foo")
	st.infoLog.Printf("bar")
	content := fmt.Sprintf("%s", &buff)

	passCond := len(s.State) == 50 && len(s.Short) == 50 &&
		strings.Contains(content, "foo") &&
		strings.Contains(content, "bar")

	if !passCond {
		t.Errorf("expected 50 for states, got %d for States and %d for Short",
			len(s.State), len(s.Short))
	}
}
