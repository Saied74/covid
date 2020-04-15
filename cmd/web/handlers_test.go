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
		errorLog:        errorLog,
		infoLog:         infoLog,
		covidProjectURL: "http://localhost:8090", //https://covidtracking.com/api/states/daily",
		State:           states,
		Short:           short,
	}

	req := httptest.NewRequest("GET", `http:localhost:8080`, nil)
	w := httptest.NewRecorder()
	s.homeHandler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	failCond := resp.StatusCode != 200 && string(body) !=
		`body <html><head></head><body><span>FOOBAR</span>
        </body></html>`

	if failCond {
		t.Errorf(`expected 200 StatusCode and <html><head></head<body><span>FOOBAR</span>
      </body></html> got status code %d and body %s`,
			resp.StatusCode, string(body))
	}

}

func TestGenHanlder(t *testing.T) {
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
		errorLog:        errorLog,
		infoLog:         infoLog,
		covidProjectURL: "http://localhost:8090", //https://covidtracking.com/api/states/daily",
		pattern:         pattern,
		State:           states,
		Short:           short,
		StateList:       []string{},
	}

	req := httptest.NewRequest("POST", `http://localhost:8080`, nil)
	req.Form = map[string][]string{
		"graphType":   []string{"Bar"},
		"fieldType":   []string{"positive"},
		"stateCheck0": []string{},
		"stateCheck1": []string{},
	}
	w := httptest.NewRecorder()
	s.genHandler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	passCond := resp.StatusCode == 200 &&
		strings.Contains(string(body), "var AL = {") &&
		strings.Contains(string(body), "var AK = {")
		// strings.Contains(string(body), "submit")

	if !passCond {
		t.Errorf(`expected 200 StatusCode and did not find "form", "fname" and
      "submit" in the body,instead got %d and found for body %s`,
			resp.StatusCode, string(body))
	}
	if passCond {
		s.infoLog.Println("passed first pass")
	}
	s.StateList = []string{}

	req = httptest.NewRequest("POST", `http://localhost:8080`, nil)
	req.Form = map[string][]string{
		"graphType": []string{},
		"fieldType": []string{},
	}
	w = httptest.NewRecorder()
	s.genHandler(w, req)

	resp = w.Result()
	// body, _ = ioutil.ReadAll(resp.Body)

	failCond := resp.StatusCode != 400 || len(s.Xdata) != 0 || len(s.Ydata) != 1

	if failCond {
		t.Errorf(`expected 200 StatusCode and did not find "form", "fname" and
	    "submit" in the body,instead got %d and found for body`,
			resp.StatusCode)
	}
}
