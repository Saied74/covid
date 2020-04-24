//helper functions for the handlers

package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

func (s *sT) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	s.errorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (s *sT) clientError(w http.ResponseWriter, status int, element string) {
	s.errorLog.Printf("%s was not found", element)

	http.Error(w, http.StatusText(status), status)
}

//I am building loggers as closurs so there can be only one of each
//type of logger but used both in the test and in the normal runtime
//reason for the "out" is of the io.Writer type is so I can use os.Stdout
//for regular logging and bytes.Buffer for testing

func getInfoLogger(out io.Writer) func() *log.Logger {
	infoLog := log.New(out, "INFO\t", log.Ldate|log.Ltime)
	return func() *log.Logger {
		return infoLog
	}
}

func getErrorLogger(out io.Writer) func() *log.Logger {
	errorLog := log.New(out, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return func() *log.Logger {
		return errorLog
	}
}

//it proved difficult to create the traces for multiple lines for Plotly
//using Go template language.  So, I decided to build the file manually
//with the buildPlot function here.  For the format of this file, you can
//checkout the Plotly JavaScript webpages.
func (s *ScreenType) buildPlot() *string {
	plot := "{{ define \"plotdata\" }}"
	for n, state := range s.StateList {
		plot += "\nvar " + state + " = {\n  x: ["
		for _, xdata := range s.Xdata {
			plot += "\"" + xdata + "\"" + ", "
		}
		plot = strings.TrimSuffix(plot, ", ")
		plot += "],\n  y: ["
		for _, ydata := range s.Ydata[n] {
			plot += ydata + ", "
		}
		plot = strings.TrimSuffix(plot, ", ")
		plot += "],\n  type: "
		plot += "\"" + s.GraphType + "\"" + ",\n"
		plot += "name: \"" + state + "\"\n};"
	}
	plot += "\nvar data = ["
	for _, state := range s.StateList {
		plot += state + ", "
	}
	plot = strings.TrimSuffix(plot, ", ")
	plot += "];\n"
	plot += "{{end}}"

	return &plot
}

func newTemplateCache(tmpls []string) (*template.Template, error) {

	ts, err := template.ParseFiles(tmpls...)
	if err != nil {
		return nil, fmt.Errorf("template did not parse in newTemplateCache %v", err)
	}

	return ts, nil
}

func (s *sT) copyScreenData() *ScreenType {
	var sst = ScreenType{}
	sst.State = s.state
	sst.Short = s.short
	sst.Fields = s.fields
	return &sst
}
