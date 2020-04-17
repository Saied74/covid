//for information, read the Readme.MD file

package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

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
	ipAddress       string     //server ip address and port number
	errorLog        *log.Logger
	infoLog         *log.Logger
}

//getFields estracts the fields to be displayed in the drop down menue on
//the web page for the field to be plotted.  It ignores date (which is for the
//horizontal axis and state which is shown on the left side of the screen)
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

//This is pulled out so durng unit testing s struct can be modified.
func (s *StatesType) editStateType(out io.Writer) {
	infoLog := getInfoLogger(out)
	errorLog := getErrorLogger(out)
	s.State = states
	s.Short = short
	s.errorLog = errorLog()
	s.infoLog = infoLog()
}

//The same with routes, it is for testability reason
func (s *StatesType) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/home", s.homeHandler)
	// mux.HandleFunc("/test", s.testHandler)
	mux.HandleFunc("/generate", s.genHandler)
	return mux
}

func main() {
	var err error
	var s StatesType

	//user can change the name of the configuration file using this flag
	config := flag.String("c", "config.csv", "Configuratoin file name")
	//user can provide an environment variable pointing to the directory
	//containing the configuration file
	environ := flag.String("e", "search order", "Env Variable for Config file location")
	flag.Parse()
	err = s.setUp(*config, *environ)
	if err != nil {
		s.errorLog.Fatal("Did not succeed configuring ", err)
	}
	err = s.validateConfigs()
	if err != nil {
		s.errorLog.Fatal("configs did not validate ", err)
	}
	//get the pattern for parsing JSON file
	s.pattern, err = virusdata.GetPattern(s.patternFile) //("../../config/pattern.csv")
	if err != nil {
		log.Fatal("reading pattern", err)
	}
	s.editStateType(os.Stdout)
	s.getFields()

	mux := s.routes()
	srv := &http.Server{
		Addr:     s.ipAddress,
		ErrorLog: s.errorLog,
		Handler:  s.recoverPanic(s.logRequest(mux)),
	}

	s.infoLog.Printf("Starting server on %s", s.ipAddress)

	err = srv.ListenAndServe()
	s.errorLog.Fatal(err)
}
