//reading files from the website https://covidtracking.com/data/
//in JSON formt and outputting them to be plotted by excel

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	lexer "github.com/Saied74/Lexer2"
)

var pickData struct {
	date       string
	state      string
	fieldName  string
	fieldValue string
}

var states []string
var outputFileName string
var interimFile = make(map[string]map[string]string)
var dateIndex []string
var csvRecords [][]string

func getPattern(fileName string) ([][]string, error) {
	var pattern [][]string
	content, err := ioutil.ReadFile(fileName) //get the whole file
	if err != nil {
		return [][]string{}, fmt.Errorf("open error %v on file %s", err, fileName)
	}
	pat1 := strings.Split(string(content), "\n") //split into lines
	for _, pat2 := range pat1 {                  //scan the lines
		pat3 := strings.Split(pat2, "|") //split into comma seperated fields
		pattern = append(pattern, pat3)  //append to the output
	}
	return pattern, nil
}

func inSlice(candidate string) bool {
	for _, element := range states {
		if element == candidate {
			return true
		}
	}
	return false
}

func processItem() {
	_, ok := interimFile[pickData.date]
	if ok {
		interimFile[pickData.date][pickData.state] = pickData.fieldValue
		return
	}
	interimFile[pickData.date] = map[string]string{
		pickData.state: pickData.fieldValue,
	}
	return
}

func buildDateIndex() {
	for key := range interimFile {
		dateIndex = append(dateIndex, key)
	}
	sort.Strings(dateIndex)
}

func writeRecords(fileName string, records [][]string) {
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

func buildOutputRecords() {
	var outputLine []string
	outputLine = append(outputLine, "Date")
	for _, state := range states {
		outputLine = append(outputLine, state)
	}
	csvRecords = append(csvRecords, outputLine)
	outputLine = []string{}
	for _, date := range dateIndex {
		outputLine = append(outputLine, date)
		for _, state := range states {
			outputLine = append(outputLine, interimFile[date][state])
		}
		csvRecords = append(csvRecords, outputLine)
		outputLine = []string{}
	}
}

func main() {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	// Make request
	response, err := client.Get("https://covidtracking.com/api/states/daily")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Copy data from the response to a byte buffer
	var buf bytes.Buffer
	n, err := io.Copy(&buf, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Number of bytes copied to STDOUT:", n)

	//get pattern for the JSON file
	pattern, err := getPattern("../config/pattern.csv")
	if err != nil {
		log.Fatal("reading pattern", err)
	}

	//get graph output requirement data
	outReq, err := getPattern("../config/outreq.csv")
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
			states = item[1:]
		case "field":
			pickData.fieldName = item[1]
		case "file":
			outputFileName = item[1]
		}
	}
	// fmt.Println("States:", states)

	pickData.fieldName = "death"

	item := lexer.Lex(pattern, buf.String())
	// var itemKey, itemValue string
	var start, done bool
	for {
		newItem := <-item
		switch newItem.ItemKey {
		case "nodeType":
			start = true
		case "object":
			start = false
		case "EOF":
			done = true
		}
		if start {
			switch newItem.ItemKey {
			case "dateChecked":
				pickData.date = newItem.ItemValue
			case "state":
				pickData.state = newItem.ItemValue
				pickData.state = strings.TrimPrefix(pickData.state, `"`)
				pickData.state = strings.TrimSuffix(pickData.state, `"`)
			case pickData.fieldName:
				pickData.fieldValue = newItem.ItemValue
				if newItem.ItemValue == "null" {
					pickData.fieldValue = ""
				}
			}
		}
		if !start && inSlice(pickData.state) {
			processItem()
			// fmt.Println(pickData)
		}
		if done {
			break
		}
	}
	buildDateIndex()
	buildOutputRecords()
	writeRecords("../data/"+outputFileName, csvRecords)
	for _, key := range dateIndex {
		fmt.Println(key, interimFile[key])
	}
}
