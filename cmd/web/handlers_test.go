package main

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {

	st := sT{}
	st.buildST(os.Stdout)
	st.covidProjectURL = "http://localhost:8080"

	req := httptest.NewRequest("GET", `http:localhost:8080`, nil)
	w := httptest.NewRecorder()
	st.homeHandler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	passCond := resp.StatusCode == 200 && !strings.Contains(string(body), "block")

	if !passCond {
		t.Errorf(`expected 200 StatusCode and expected the word nothing in the body
      got status code %d and body %s`,
			resp.StatusCode, string(body))
	}
}

// ----------------------TestGenHanlder-------------------------------

func TestGenHanlder(t *testing.T) {

	st := sT{}
	st.buildST(os.Stdout)
	st.covidProjectURL = "http://localhost:8080"

	type genTest struct {
		genForm   *map[string][]string
		resCode   int //results from the server
		contains1 string
		contains2 string
	}

	req := httptest.NewRequest("POST", `http://localhost:8080`, nil)

	formOne := map[string][]string{
		"graphType":   []string{"Bar"},
		"fieldType":   []string{"positive"},
		"stateCheck0": []string{},
		"stateCheck1": []string{},
	}

	formTwo := map[string][]string{
		"graphType": []string{},
		"fieldType": []string{},
	}

	testPattern := []genTest{
		genTest{
			genForm:   &formOne,
			resCode:   200,
			contains1: "var AL = {",
			contains2: "var AK = {",
		},
		genTest{
			genForm:   &formTwo,
			resCode:   400,
			contains1: "x: [],",
			contains2: "y: [],",
		},
	}

	for n, seq := range testPattern {

		req.Form = *seq.genForm

		w := httptest.NewRecorder()
		st.genHandler(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		passCond := resp.StatusCode == seq.resCode &&
			strings.Contains(string(body), seq.contains1) &&
			strings.Contains(string(body), seq.contains2)
			// strings.Contains(string(body), "submit")
		if passCond {
			st.infoLog.Printf("passed genHandler test run %d", n)
		}

		if !passCond {
			t.Errorf(`expected %d StatusCode %s and %s in the body but got %d as
      StatusCode and did not get the content, here is the body %s`,
				seq.resCode, seq.contains1, seq.contains2, resp.StatusCode, string(body))
		}
	}
}
