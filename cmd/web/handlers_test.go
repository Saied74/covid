package main

import (
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"silverslanellc.com/covid/pkg/virusdata"
)

func TestHomeHandler(t *testing.T) {

	dummyTmpl := `{{define "plotdata"}}nothing{{end}}`
	writeFile("../../ui/html/plot.partial.tmpl", &dummyTmpl)
	defer os.Remove("../../ui/html/plot.partial.tmpl")

	s := StatesType{
		templateFiles: []string{"../../ui/html/base.page.tmpl",
			"../../ui/html/plot.partial.tmpl"},
		plotFile:        "../../ui/html/plot.partial.tmpl",
		errorLog:        getErrorLogger(os.Stdout)(),
		infoLog:         getInfoLogger(os.Stdout)(),
		covidProjectURL: "http://localhost:8090", //https://covidtracking.com/api/states/daily",
		State:           states,
		Short:           short,
	}

	req := httptest.NewRequest("GET", `http:localhost:8080`, nil)
	w := httptest.NewRecorder()
	s.homeHandler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	passCond := resp.StatusCode == 200 && strings.Contains(string(body), "nothing")

	if !passCond {
		t.Errorf(`expected 200 StatusCode and expected the word nothing in the body
      got status code %d and body %s`,
			resp.StatusCode, string(body))
	}
}

//----------------------TestGenHanlder-------------------------------

func TestGenHanlder(t *testing.T) {

	type genTest struct {
		genForm   *map[string][]string
		resCode   int //results from the server
		contains1 string
		contains2 string
	}

	dummyTmpl := `{{define "plotdata"}}nothing{{end}}`
	writeFile("../../ui/html/plot.partial.tmpl", &dummyTmpl)
	defer os.Remove("../../ui/html/plot.partial.tmpl")

	pattern, err := virusdata.GetPattern("../../config/pattern.csv")
	if err != nil {
		log.Fatal("reading pattern", err)
	}

	s := StatesType{
		templateFiles: []string{"../../ui/html/base.page.tmpl",
			"../../ui/html/plot.partial.tmpl"},
		plotFile:        "../../ui/html/plot.partial.tmpl",
		errorLog:        getErrorLogger(os.Stdout)(),
		infoLog:         getInfoLogger(os.Stdout)(),
		covidProjectURL: "http://localhost:8090", //https://covidtracking.com/api/states/daily",
		pattern:         pattern,
		State:           states,
		Short:           short,
		StateList:       []string{},
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
		s.genHandler(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		passCond := resp.StatusCode == seq.resCode &&
			strings.Contains(string(body), seq.contains1) &&
			strings.Contains(string(body), seq.contains2)
			// strings.Contains(string(body), "submit")
		if passCond {
			s.infoLog.Printf("passed genHandler test run %d", n)
		}

		if !passCond {
			t.Errorf(`expected %d StatusCode %s and %s in the body but got %d as
      StatusCode and did not get the content, here is the body %s`,
				seq.resCode, seq.contains1, seq.contains2, resp.StatusCode, string(body))
		}
	}
}
