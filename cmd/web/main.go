package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"silverslanellc.com/covid/pkg/virusdata"
)

//StatesType is exported for use in templates
type StatesType struct {
	State           []string   //states for the left side of the web page
	Short           []string   //state abbreviations for the API parsing
	Fields          []string   //for field selector on top of the page
	Xdata           []string   //for plotting the x axis of the graph
	Ydata           [][]string //for plotting the y axis of the graph
	GraphType       string     //graph type as specified on the web page
	StateList       []string   //list of salected states
	Selected        string     //which field was selected to be extracted from data
	pattern         [][]string //pattern for the lexer
	appHome         string     //the home address of the applicatoin
	patternFile     string     //name of the pattern file for parsing
	csvOutputFile   string     //name of the file for formatting the csv output
	covidProjectURL string     //URL for the covid tracking project
	templateFiles   []string   //names of files to be parsed for the web page
	plotFile        string     //one of the templateFiles must also be a plot file
}

func (s *StatesType) getFields() {
	s.Fields = []string{}
	for _, row := range s.pattern {
		if len(row) > 1 {
			if row[0] == "attribute" && row[1] != "date" && row[1] != "state" {
				s.Fields = append(s.Fields, row[1])
			}
		}
	}
}

func main() {
	var err error
	var s StatesType
	//user can change the name of the configuration file using this flag
	config := flag.String("c", "config.csv", "Configuratoin file name")
	//user can provide an environment variable pointing to the directory
	//containing the configuration file
	environ := flag.String("e", ".", "Env Variable for Config file location")
	flag.Parse()
	err = s.setUp(*config, *environ)
	if err != nil {
		log.Fatal("Did not succeed configuring ", err)
	}
	err = s.validateConfigs()
	if err != nil {
		log.Fatal("configs did not validate ", err)
	}
	fmt.Println("App Home:\t", s.appHome)
	fmt.Println("Pattern File:\t", s.patternFile)
	fmt.Println("CSV Output File:\t", s.csvOutputFile)
	fmt.Println("Covid Project URL:\t", s.covidProjectURL)
	fmt.Println("Template Files:\t", s.templateFiles)
	fmt.Println("Plot file:\t", s.plotFile)
	// os.Exit(1)

	//get the pattern for parsing JSON file
	s.pattern, err = virusdata.GetPattern(s.patternFile) //("../../config/pattern.csv")
	if err != nil {
		log.Fatal("reading pattern", err)
	}
	s.getFields()
	s.State = states
	s.Short = short
	// s.Fields = fields // TODO: this will have to come here

	mux := http.NewServeMux()
	mux.HandleFunc("/home", s.homeHandler)
	mux.HandleFunc("/generate", s.genHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
