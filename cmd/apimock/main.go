//this is a mock of the covid tracking project API
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type appType struct {
	testData  *[]byte
	ipAddress string
	errorLog  *log.Logger
	infoLog   *log.Logger
}

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func readData(fileName string) (*[]byte, error) {

	content, err := ioutil.ReadFile(fileName) //get the whole file
	if err != nil {
		return nil, fmt.Errorf("open error %v on file %s", err, fileName)
	}
	return &content, nil
}
func (app *appType) homeHandler(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	w.Write(*app.testData)
}

func main() {
	cntPtr, err := readData("./data.json")
	if err != nil {
		errorLog.Fatal("data.json datafile read failed because", err)
	}
	app := appType{
		testData:  cntPtr,
		ipAddress: "localhost:8090",
		errorLog:  errorLog,
		infoLog:   infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/home", app.homeHandler)

	srv := &http.Server{
		Addr:     app.ipAddress,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", app.ipAddress)

	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
