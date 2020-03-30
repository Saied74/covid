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

//type for picking individual data points
type pick struct {
	date       string
	state      string
	fieldName  string
	fieldValue string
}
type stateString []string
type interim map[string]map[string]string
type csvRecordType [][]string
type bundle struct {
	interimFiles interim
	pickFile     pick
}

var outputFileName string

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

func (b *bundle) lexInputData(pattern [][]string, inputData *string,
	states stateString) {
	var start, done bool
	// b.pickFile.fieldName = "death"

	item := lexer.Lex(pattern, *inputData)

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
				tmpDate := strings.Split(newItem.ItemValue, "T")
				if len(tmpDate) != 2 {
					log.Fatal("encounted a badly formatted date", newItem.ItemValue)
				}
				b.pickFile.date = tmpDate[0]
			case "state":
				b.pickFile.state = newItem.ItemValue
				b.pickFile.state = strings.TrimPrefix(b.pickFile.state, `"`)
				b.pickFile.state = strings.TrimSuffix(b.pickFile.state, `"`)
			case b.pickFile.fieldName:
				b.pickFile.fieldValue = newItem.ItemValue
				if newItem.ItemValue == "null" {
					b.pickFile.fieldValue = ""
				}
			}
		}
		if !start && states.inSlice(b.pickFile.state) {
			b.processItem()
		}
		if done {
			break
		}
	}
}

func (s *stateString) inSlice(candidate string) bool {
	for _, element := range *s {
		if element == candidate {
			return true
		}
	}
	return false
}

func (b *bundle) processItem() {
	_, ok := b.interimFiles[b.pickFile.date]
	if ok {
		b.interimFiles[b.pickFile.date][b.pickFile.state] = b.pickFile.fieldValue
		return
	}
	b.interimFiles[b.pickFile.date] = map[string]string{
		b.pickFile.state: b.pickFile.fieldValue,
	}
	return
}

func (b *bundle) buildDateIndex() []string {
	var dateIndex []string
	for key := range b.interimFiles {
		dateIndex = append(dateIndex, key)
	}
	sort.Strings(dateIndex)
	return dateIndex
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

func (b *bundle) buildOutputRecords(s stateString, dateIndex []string) *csvRecordType {
	var csvRecords csvRecordType
	var outputLine []string
	outputLine = append(outputLine, "Date")
	for _, state := range s {
		outputLine = append(outputLine, state)
	}
	csvRecords = append(csvRecords, outputLine)
	outputLine = []string{}
	for _, date := range dateIndex {
		outputLine = append(outputLine, date)
		for _, state := range s {
			outputLine = append(outputLine, b.interimFiles[date][state])
		}
		csvRecords = append(csvRecords, outputLine)
		outputLine = []string{}
	}
	return &csvRecords
}

func main() {
	// var pickData pick      //for holding individual datapoints
	var states stateString //for holding the list of states
	var bundleFiles bundle
	bundleFiles.interimFiles = make(interim)
	// var interimFile = make(interim)
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
			bundleFiles.pickFile.fieldName = item[1]
		case "file":
			outputFileName = item[1]
		}
	}
	inputData := buf.String()
	bundleFiles.lexInputData(pattern, &inputData, states)
	dateIndex := bundleFiles.buildDateIndex()
	csvRecords := *bundleFiles.buildOutputRecords(states, dateIndex)
	writeRecords("../data/"+outputFileName, csvRecords)
	for _, key := range dateIndex {
		fmt.Println(key, bundleFiles.interimFiles[key])
	}
}
