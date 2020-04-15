//odds and ends functions

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"silverslanellc.com/covid/pkg/virusdata"
)

func inSlice(item string, target []string) bool {
	for _, item2 := range target {
		if item == item2 {
			return true
		}
	}
	return false
}

func fileExists(f string) bool {
	info, err := os.Lstat(f)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (s *StatesType) setUp(config string, environ string) error {
	var fileName string
	var err error
	//process config file in app directory
	fileName = "./" + config
	if fileExists(fileName) {
		err = s.processConfig(fileName)
		if err != nil {
			return fmt.Errorf("reading from app directory %v", err)
		}
		return nil
	}
	//process config file in $GOPATH/src/covid ddirectory
	fileName = os.Getenv("GOPATH") + "/src/covid/" + config
	if fileExists(fileName) {
		err = s.processConfig(fileName)
		if err != nil {
			return fmt.Errorf("reading from $GOPATH/src/covid %v", err)
		}
		return nil
	}
	//process config file in directory pointed to by environment variable
	environ = strings.TrimPrefix(environ, "$")
	environ = strings.ToUpper(environ)
	fileName = os.Getenv(environ)
	if len(fileName) == 0 {
		return fmt.Errorf("no environment variable with name %s was found", environ)
	}
	fileName += "/" + config
	if fileExists(fileName) {
		err = s.processConfig(fileName)
		if err != nil {
			return fmt.Errorf("reading from environment variable %s, %v", environ, err)
		}
		return nil
	}
	return fmt.Errorf(" did not find config file from environment variable %s, and filename %s", environ, fileName)
}

func (s *StatesType) processConfig(fileName string) error {
	pattern, err := virusdata.GetPattern(fileName)
	if err != nil {
		return fmt.Errorf("reading config file %s, with error %v", fileName, err)
	}
	for _, line := range pattern {
		if len(line) < 2 {
			continue
		}
		if strings.HasPrefix(line[0], "#") {
			continue
		}
		switch line[0] {
		case "appHome":
			s.appHome = line[1]
		case "patternFile":
			s.patternFile = line[1]
		case "csvOutputFile":
			s.csvOutputFile = line[1]
		case "covidProjectURL":
			s.covidProjectURL = line[1]
		case "templateFile":
			s.templateFiles = append(s.templateFiles, line[1])
		case "plotFile":
			s.plotFile = line[1]
		case "ipAddress":
			s.ipAddress = line[1]
		}
	}
	return nil
}

func (s *StatesType) validateConfigs() error {
	var plotMatch bool //when true, plot file is in the list of templateFiles

	//check for a valid appHome
	if len(s.appHome) == 0 {
		return fmt.Errorf("no appHome verable was found in the config file")
	}
	//if neither environment variable nor absolute
	if !strings.HasPrefix(s.appHome, "/") && !strings.HasPrefix(s.appHome, "$") {
		s.appHome = ""
		return fmt.Errorf("malformed appHome in config file %s", s.appHome)
	}
	if strings.HasPrefix(s.appHome, "$") { //in case of environment variable
		appHome := strings.Split(s.appHome, "/")
		if len(appHome) < 2 {
			return fmt.Errorf("appHome veriable did not split right, %v", appHome)
		}
		home := strings.TrimPrefix(appHome[0], "$")
		homeAddr := os.Getenv(strings.ToUpper(home))
		if len(homeAddr) == 0 {
			return fmt.Errorf("no appHome environment variable was returned for %s", home)
		}
		appHome[0] = homeAddr
		s.appHome = strings.Join(appHome, "/")
	}

	//check the pattern file
	if len(s.patternFile) == 0 {
		return fmt.Errorf("no patternFile variable was provided in the config file")
	}
	s.patternFile = filepath.Join(s.appHome, s.patternFile)
	if !fileExists(s.patternFile) {
		return fmt.Errorf("Pattern file %s was not found", s.patternFile)
	}

	//ceck the csv output file
	if len(s.csvOutputFile) == 0 {
		return fmt.Errorf("no csvOutputFile was provided in the config file")
	}
	s.csvOutputFile = filepath.Join(s.appHome, s.csvOutputFile)
	if !fileExists(s.csvOutputFile) {
		return fmt.Errorf("CSV ouput file %s was not found", s.csvOutputFile)
	}

	//check the covid project url
	if len(s.covidProjectURL) == 0 {
		return fmt.Errorf("no covidProjectURL was provided in the config file")
	}
	if !strings.HasPrefix(s.covidProjectURL, "https://") {
		return fmt.Errorf("malformed covidProjectURL %s", s.covidProjectURL)
	}

	//check the plot file, fileExists will be checked later
	if len(s.plotFile) == 0 {
		return fmt.Errorf("no plotFile was provided in the config file")
	}
	for _, templFile := range s.templateFiles {
		if s.plotFile == templFile {
			plotMatch = true
		}
	}
	if !plotMatch {
		return fmt.Errorf("plotfile %s was not in the list of template files %v",
			s.plotFile, s.templateFiles)
	}
	s.plotFile = filepath.Join(s.appHome, s.plotFile)

	if !fileExists(s.plotFile) {
		plotContent := `{{ define "plotdata" }} {{ end }}`
		writeFile(s.plotFile, &plotContent)
	}

	//check templateFiles
	if len(s.templateFiles) == 0 {
		return fmt.Errorf("no templateFiles was provided in the config file")
	}
	// log.Println("before the range", s.templateFiles)
	for n, tmp := range s.templateFiles {
		s.templateFiles[n] = filepath.Join(s.appHome, tmp)
		if !fileExists(s.templateFiles[n]) {
			return fmt.Errorf("template file %s was not found", s.templateFiles[n])
		}
	}
	if len(s.ipAddress) == 0 {
		return fmt.Errorf("no ip address field was provided in config file")
	}

	return nil
}
