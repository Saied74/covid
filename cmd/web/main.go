//for information, read the Readme.MD file

package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"silverslanellc.com/covid/pkg/virusdata"
)

//ScreenType is exported for use in templates
type ScreenType struct {
	State     []string   //states for the left side of the web page
	Short     []string   //state abbreviations for the API parsing
	Fields    []string   //for field selector on top of the page
	Xdata     []string   //for plotting the x axis of the graph
	Ydata     [][]string //for plotting the y axis of the graph
	GraphType string     //graph type as specified on the web page
	StateList []string   //list of salected states
	Selected  string     //which field was selected to be extracted from data
}

type sT struct {
	state           []string   //states for the left side of the web page
	short           []string   //state abbreviations for the API parsing
	fields          []string   //for field selector on top of the page
	pattern         [][]string //pattern for the lexer
	appHome         string     //the home address of the applicatoin
	patternFile     string     //name of the pattern file for parsing
	csvOutputFile   string     //name of the file for formatting the csv output
	covidProjectURL string     //URL for the covid tracking project
	templateFiles   []string   //names of files to be parsed for the web page
	cache           *template.Template
	errorLog        *log.Logger
	infoLog         *log.Logger
}

//getFields estracts the fields to be displayed in the drop down menue on
//the web page for the field to be plotted.  It ignores date (which is for the
//horizontal axis and state which is shown on the left side of the screen)
func (st *sT) getFields() {
	st.fields = []string{}
	for _, row := range st.pattern {
		if len(row) > 1 {
			if row[0] == "attribute" && row[1] != "date" && row[1] != "state" {
				st.fields = append(st.fields, row[1])
			}
		}
	}
}

//This is pulled out so durng unit testing s struct can be modified.
func (st *sT) buildST(out io.Writer) {
	var err error
	infoLog := getInfoLogger(out)
	errorLog := getErrorLogger(out)
	st.state = states
	st.short = short
	st.templateFiles = templateFiles
	st.patternFile = patternFile
	st.cache, err = newTemplateCache(templateFiles)
	if err != nil {
		st.errorLog.Fatal("template cash did not build ", err)
	}
	st.pattern, err = virusdata.GetPattern(patternFile)
	if err != nil {
		log.Fatal("reading pattern", err)
	}
	st.getFields()
	st.errorLog = errorLog()
	st.infoLog = infoLog()
}

//The same with routes, it is for testability reason
func (st *sT) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/home", st.homeHandler)
	mux.HandleFunc("/generate", st.genHandler)
	return mux
}

func main() {
	var err error
	var st sT

	//Flags
	covidProjectURL := flag.String("cpURL",
		"https://covidtracking.com/api/states/daily", "covid project URL")
	ipAddress := flag.String("ipa", ":8080", "server ip address")

	st.buildST(os.Stdout)
	st.covidProjectURL = *covidProjectURL

	mux := st.routes()
	srv := &http.Server{
		Addr:     *ipAddress,
		ErrorLog: st.errorLog,
		Handler:  st.recoverPanic(st.logRequest(mux)),
	}

	st.infoLog.Printf("Starting server on %s", *ipAddress)

	err = srv.ListenAndServe()
	st.errorLog.Fatal(err)
}
