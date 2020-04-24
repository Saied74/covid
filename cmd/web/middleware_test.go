package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogRequest(t *testing.T) {
	var buff bytes.Buffer
	st := sT{
		infoLog: getInfoLogger(&buff)(),
	}
	req := httptest.NewRequest("GET", `http://localhost:8080/generate`, nil)
	w := httptest.NewRecorder()

	testHandle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	st.logRequest(testHandle).ServeHTTP(w, req)

	output := fmt.Sprintf("%s", &buff)
	successCond := strings.Contains(output, "HTTP/1.1 GET /generate") &&
		!strings.Contains(output, "HTTP/1.1 POST /generate")

	if !successCond {
		t.Errorf("fail, expected content HTTP/1.1 GET /generate got %s", output)
	}

	//---------------------  Round 2 -------------------------------------
	var buff2 bytes.Buffer
	st.infoLog = getInfoLogger(&buff2)()

	req2 := httptest.NewRequest("POST", `http://localhost:8080/home`, nil)

	st.logRequest(testHandle).ServeHTTP(w, req2)

	output = fmt.Sprintf("%s", &buff2)
	successCond = strings.Contains(output, "HTTP/1.1 POST /home") &&
		!strings.Contains(output, "HTTP/1.1 GET /generate")

	if !successCond {
		t.Errorf("fail, expected content HTTP/1.1 POST /home got %s", output)
	}
}

func TestRecoverPanic(t *testing.T) {
	var buff bytes.Buffer
	st := sT{
		infoLog:  getInfoLogger(&buff)(),
		errorLog: getErrorLogger(&buff)(),
	}
	req := httptest.NewRequest("GET", `http://localhost:8080/generate`, nil)
	w := httptest.NewRecorder()

	testHandle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Oops! something went wrong")
	})

	st.recoverPanic(testHandle).ServeHTTP(w, req)
	output := fmt.Sprintf("%s", &buff)
	if !strings.Contains(output, "Oops! something went wrong") {
		t.Errorf("fail, expected content of Oops! something went wrong got %s",
			output)
	}
}
