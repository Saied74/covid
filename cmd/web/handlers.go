package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"silverslanellc.com/covid/pkg/virusdata"
)

//the second file generates the data for multiple graphs.  Using the go
//templates proved to be too complicaated (or impossible)
var files = []string{
	"../../ui/html/base.page.tmpl",    //the base (and only) page
	"../../ui/html/plot.partial.tmpl", //generated JS section for multiple graphs
}

func (s *StatesType) homeHandler(w http.ResponseWriter, r *http.Request) {
	tt, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal("group files did not parse ", err)
	}
	tt.Execute(w, s)
}

func (s *StatesType) genHandler(w http.ResponseWriter, r *http.Request) {
	//set up the data structure to read and lex the input data
	var pickData virusdata.Pick
	pickData.InterimFiles = make(virusdata.Interim)
	//clear thsese variables so the "append" function works rigth
	s.Xdata = []string{}
	s.Ydata = [][]string{}
	s.StateList = []string{}

	err := r.ParseForm() //parse request, handle error
	if err != nil {
		log.Println("form did not parse", err)
	}
	//pick out the requested graph type and the field
	graphType := r.Form["graphType"] //pick graph type
	if len(graphType) == 0 {
		log.Fatal("Got bad graphType from the web page", graphType) // TODO: get rid of the fatal and also send a message to the screen
	}
	s.GraphType = strings.ToLower(graphType[0])

	fieldType := r.Form["fieldType"] //pick the field to be plotted
	if len(fieldType) == 0 {
		log.Fatal("Got bad fieldType from the web page", fieldType) // TODO: get rid of the fatal and also send a message to the screen
	}
	pickData.FieldName = strings.ToLower(fieldType[0])
	s.Selected = fieldType[0]

	//now pick the states requested
	for i := range s.Short { //Short because that is how the api responds
		candidate := "stateCheck" + strconv.Itoa(i)
		for key := range r.Form {
			if key == candidate {
				pickData.StateList = append(pickData.StateList, s.Short[i])
			}
		}
	}
	s.StateList = pickData.StateList

	//get the JSON file by making the API call
	inputData := virusdata.GetData() // TODO: check to see if any data was returned

	//get the pattern for parsing JSON file
	// TODO: use the same step to populate StateType.Fields
	pattern, err := virusdata.GetPattern("../../config/pattern.csv")
	if err != nil {
		log.Fatal("reading pattern", err)
	}

	pickData.LexInputData(pattern, inputData)     //lex the input data with the pattern
	pickData.DateList = pickData.BuildDateIndex() //format the dates

	s.Xdata = pickData.DateList
	var yLine []string
	for _, state := range s.StateList {
		for _, date := range s.Xdata {
			yLine = append(yLine, pickData.InterimFiles[date][state])
		}
		s.Ydata = append(s.Ydata, yLine)
		yLine = []string{}
	}
	//build the plot file to be parsed with the other template
	// TODO: handle exceptions better as discssed elsewhere instead of log.Fatal
	plot := s.buildPlot()
	log.Println("Plot: ", plot)
	f, err := os.Create("../../ui/html/plot.partial.tmpl")
	if err != nil {
		log.Fatal("Could not create ../../ui/html/plot.partial.tmp", err)
	}
	defer f.Close()
	_, err = f.WriteString(plot)
	if err != nil {
		log.Fatal("could not write the plot file", err)
	}
	tt, err := template.ParseFiles(files...) //parse html files, handle error
	if err != nil {
		log.Fatal("group files did not parse ", err)
	}
	tt.Execute(w, s)
}

//it proved difficult to create the traces for multiple lines for Plotly
//using Go template language.  So, I decided to build the file manually
//with the buildPlot function here.  For the format of this file, you can
//checkout the Plotly JavaScript webpages.
func (s *StatesType) buildPlot() string {
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
	return plot
}
