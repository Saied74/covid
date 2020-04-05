package main

import (
	"log"
	"net/http"
)

//StatesType is exported for use in templates
type StatesType struct {
	State []string //states for the left side of the web page
	Short []string //state abbreviations for the API parsing
	// Tick      []bool     //maybe we don't need it
	Fields    []string   //for field selector on top of the page
	Xdata     []string   //for plotting the x axis of the graph
	Ydata     [][]string //for plotting the y axis of the graph
	GraphType string     //graph type as specified on the web page
	StateList []string
	Selected  string //which field was selected to be extracted from data
}

func main() {

	var s StatesType
	s.State = states
	s.Short = short
	s.Fields = fields

	mux := http.NewServeMux()
	mux.HandleFunc("/home", s.homeHandler)
	mux.HandleFunc("/generate", s.genHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
