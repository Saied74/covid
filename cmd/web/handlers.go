package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"silverslanellc.com/covid/pkg/virusdata"
)

//This method that is commented out to keep things safe is for testing the
//panic recovery middleware.
// func (s *StatesType) testHandler(w http.ResponseWriter, r *http.Request) {
// 	tt, err := template.ParseFiles(s.templateFiles...)
// 	if err != nil {
// 		s.serverError(w, err)
// 	}
// 	panic("oops! something went wrong")
//
// 	tt.Execute(w, s)
// }

func (s *StatesType) homeHandler(w http.ResponseWriter, r *http.Request) {
	tt := template.Must(template.ParseFiles(s.templateFiles...))
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
	//pick out the requested graph type and the field; handle exceptions
	graphType := r.Form["graphType"] //pick graph type
	if len(graphType) == 0 {
		s.GraphType = "bar"
		s.clientError(w, http.StatusBadRequest, "graphType")
	}
	if len(graphType) != 0 {
		s.GraphType = strings.ToLower(graphType[0])
	}
	fieldType := r.Form["fieldType"] //pick the field to be plotted
	if len(fieldType) == 0 {
		s.Selected = "positive"
		s.clientError(w, http.StatusBadRequest, "fieldType")
	}
	if len(fieldType) != 0 {
		pickData.FieldName = fieldType[0]
		s.Selected = fieldType[0]
	}
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
	if len(s.StateList) == 0 {
		s.StateList = []string{"NY"}
	}
	//get the JSON file by making the API call
	inputData, err := virusdata.GetData(s.covidProjectURL) // TODO: check to see if any data was returned
	if err != nil && !strings.HasSuffix(fmt.Sprintf("%v", err), "connection refused") {
		s.errorLog.Fatal("Connection was refused with error", err)
		// s.serverError(w, err)
	}

	if inputData != nil {
		pickData.LexInputData(s.pattern, inputData)   //lex the input data with the pattern
		pickData.DateList = pickData.BuildDateIndex() //format the dates
	}
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
	if s.plotFile != "" {
		err = s.buildPlot()
		if err != nil {
			s.errorLog.Printf("plot file did not build %v", err)
		}
	}

	tt := template.Must(template.ParseFiles(s.templateFiles...)) //parse html files, handle error
	tt.Execute(w, s)
}
