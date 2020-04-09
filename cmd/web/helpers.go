//helper functions for the handlers

package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

func (s *StatesType) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	s.errorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (s *StatesType) clientError(w http.ResponseWriter, status int, element string) {
	s.errorLog.Printf("%s was not found", element)

	http.Error(w, http.StatusText(status), status)
}

func (s *StatesType) notFound(w http.ResponseWriter) {
	s.clientError(w, http.StatusNotFound, "")
}

//it proved difficult to create the traces for multiple lines for Plotly
//using Go template language.  So, I decided to build the file manually
//with the buildPlot function here.  For the format of this file, you can
//checkout the Plotly JavaScript webpages.
func (s *StatesType) buildPlot() error {
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

	f, err := os.Create(s.plotFile)
	if err != nil {
		return fmt.Errorf("plotFile %s was not created because %v", s.plotFile, err)
	}
	defer f.Close()
	_, err = f.WriteString(plot)
	if err != nil {
		return fmt.Errorf("plotFilee was not written beause %v", err)
	}
	return nil
}
