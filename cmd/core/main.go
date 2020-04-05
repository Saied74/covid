//reading files from the website https://covidtracking.com/data/
//in JSON formt and outputting them to be plotted by excel

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"silverslanellc.com/covid/pkg/virusdata"
)

type csvRecordType [][]string

var outputFileName string

//writes the csv record to the output file
func writeRecords(fileName string, records csvRecordType) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Failed to create file for writing: %v \n", err)
		os.Exit(1)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	for _, item := range records {
		err := w.Write(item)
		if err != nil {
			fmt.Printf("Fail on write: %v \n", err)
		}
	}
}

//structures the data for writing csv records
func buildOutputRecords(p virusdata.Pick) csvRecordType {
	var csvRecords csvRecordType
	var outputLine []string
	outputLine = append(outputLine, "Date")
	for _, state := range p.StateList {
		outputLine = append(outputLine, state)
	}
	csvRecords = append(csvRecords, outputLine)
	outputLine = []string{}
	for _, date := range p.DateList {
		outputLine = append(outputLine, date)
		for _, state := range p.StateList {
			outputLine = append(outputLine, p.InterimFiles[date][state])
		}
		csvRecords = append(csvRecords, outputLine)
		outputLine = []string{}
	}
	return csvRecords
}

func main() {
	var pickData virusdata.Pick //for holding individual datapoints
	pickData.InterimFiles = make(virusdata.Interim)

	//get pattern for the JSON file
	pattern, err := virusdata.GetPattern("../../config/pattern.csv")
	if err != nil {
		log.Fatal("reading pattern", err)
	}

	//get graph output requirement data
	outReq, err := virusdata.GetPattern("../../config/outreq.csv")
	if err != nil {
		log.Fatal("reading output requirements", err)
	}

	//process requirements
	for i, item := range outReq {
		if len(item) < 2 {
			log.Println("Short line was read from the outreq file line:", i)
		}
		switch item[0] {
		case "state":
			pickData.StateList = item[1:]
		case "field":
			pickData.FieldName = item[1]
		case "file":
			outputFileName = item[1]
		}
	}
	inputData := virusdata.GetData()
	pickData.LexInputData(pattern, inputData)
	pickData.DateList = pickData.BuildDateIndex()
	csvRecords := buildOutputRecords(pickData)
	writeRecords("../../data/"+outputFileName, csvRecords)
	for _, key := range pickData.DateList {
		fmt.Println(key, pickData.InterimFiles[key])
	}
}
