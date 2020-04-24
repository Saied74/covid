package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"silverslanellc.com/covid/pkg/virusdata"
)

func (st *sT) homeHandler(w http.ResponseWriter, r *http.Request) {
	s := st.copyScreenData()
	clone, err := st.cache.Clone()
	if err != nil {
		st.errorLog.Fatal("cloning error in home handler ", err)
	}
	clone.Execute(w, s)
}

func (st *sT) genHandler(w http.ResponseWriter, r *http.Request) {
	s := st.copyScreenData()
	//set up the data structure to read and lex the input data
	var pickData virusdata.Pick
	pickData.InterimFiles = make(virusdata.Interim)

	err := r.ParseForm() //parse request, handle error
	if err != nil {
		st.serverError(w, err)
	}
	//pick out the requested graph type and the field; handle exceptions
	graphType := r.Form["graphType"] //pick graph type
	if len(graphType) == 0 {
		s.GraphType = "bar"
		st.clientError(w, http.StatusBadRequest, "graphType")
	}
	if len(graphType) != 0 {
		s.GraphType = strings.ToLower(graphType[0])
	}
	fieldType := r.Form["fieldType"] //pick the field to be plotted
	if len(fieldType) == 0 {
		s.Selected = "positive"
		st.clientError(w, http.StatusBadRequest, "fieldType")
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
	t0 := time.Now()
	inputData, err := virusdata.GetData(st.covidProjectURL) // TODO: check to see if any data was returned
	if err != nil && !strings.HasSuffix(fmt.Sprintf("%v", err), "connection refused") {
		st.errorLog.Fatal("Connection was refused with error", err)
		// s.serverError(w, err)
	}
	t1 := time.Now()
	st.infoLog.Printf("covid api call time %v ", t1.Sub(t0))

	if inputData != nil {
		pickData.LexInputData(st.pattern, inputData)  //lex the input data with the pattern
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

	plot := s.buildPlot()
	clone, err := st.cache.Clone()
	if err != nil {
		st.errorLog.Fatal("cloning error ", err)
	}
	ts, err := clone.Parse(*plot)
	if err != nil {
		st.errorLog.Printf("plot did not parse because %v ", err)
	}

	ts.Execute(w, s)
}
