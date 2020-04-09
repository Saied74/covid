package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"silverslanellc.com/covid/pkg/virusdata"
)

func (s *StatesType) homeHandler(w http.ResponseWriter, r *http.Request) {
	tt, err := template.ParseFiles(s.templateFiles...)
	if err != nil {
		s.serverError(w, err)
	}
	tt.Execute(w, s)
}

func (s *StatesType) genHandler(w http.ResponseWriter, r *http.Request) {
	//set up the data structure to read and lex the input data
	var pickData virusdata.Pick
	pickData.InterimFiles = make(virusdata.Interim)

	err := r.ParseForm() //parse request, handle error
	if err != nil {
		s.serverError(w, err)
	}
	//pick out the requested graph type and the field
	graphType := r.Form["graphType"] //pick graph type
	if len(graphType) == 0 {
		s.clientError(w, http.StatusBadRequest, "graphType")
	}
	s.GraphType = strings.ToLower(graphType[0])

	fieldType := r.Form["fieldType"] //pick the field to be plotted
	if len(fieldType) == 0 {
		s.clientError(w, http.StatusBadRequest, "fieldType")
	}
	pickData.FieldName = fieldType[0]
	s.Selected = fieldType[0]

	//now pick the states requested
	s.StateList = []string{}
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
	inputData, err := virusdata.GetData(s.covidProjectURL) // TODO: check to see if any data was returned
	if err != nil {
		s.serverError(w, err)
	}

	pickData.LexInputData(s.pattern, inputData)   //lex the input data with the pattern
	pickData.DateList = pickData.BuildDateIndex() //format the dates

	s.Xdata = pickData.DateList

	s.Ydata = [][]string{}
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
	err = s.buildPlot()
	if err != nil {
		s.errorLog.Printf("plot file did not build %v", err)
	}

	tt, err := template.ParseFiles(s.templateFiles...) //parse html files, handle error
	if err != nil {
		s.serverError(w, err)
	}
	tt.Execute(w, s)
}
