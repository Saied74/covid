package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type placeType int
type resultType int

type onePattern struct {
	configData *string     //actual data to be put into the configFile
	configFile string      //config file name for this test run
	filePlace  placeType   //where the file is to be stored per Readme.MD
	expect     resultType  //either error or not error
	notError   *StatesType //if the result is not error, what should it be
}

type patternType []onePattern

const (
	inDir    placeType = iota //in the same directory as the application
	inGOPATH                  //in $GOPATH/src/covid
	inENV                     //in directory pointed to by the enviroment variable in -e flag
)

const (
	gotError resultType = iota
	noError
)

var fileOne = `appHome|$HOME/Documents/gocode/src/covid
# just a comment
patternFile|config/pattern.csv
csvOutputFile|config/outreq.csv
covidProjectURL|https://covidtracking.com/api/states/daily
templateFile|ui/html/base.page.tmpl
templateFile|ui/html/plot.partial.tmpl
plotFile|ui/html/plot.partial.tmpl
ipAddress|:8080`

var fileTwo = `appHome|
patternFile|
csvOutputFile|


covidProjectURL|
templateFile|
templateFile|
plotFile|

ipAddress|`

var fileThree = `appHome|$HOME/Documents/gocode/src/covid
patternFile|
csvOutputFile|config/outreq.csv
covidProjectURL|
templateFile|ui/html/base.page.tmpl
templateFile|
# another random comment
plotFile|
ipAddress|:8080`

var fileFour = `appHome|
patternFile|config/pattern.csv
csvOutputFile|
covidProjectURL|https://covidtracking.com/api/states/daily
templateFile|ui/html/base.page.tmpl
templateFile|
plotFile|ui/html/plot.partial.tmpl
ipAddress|`

var resultOne = StatesType{
	appHome:         "$HOME/Documents/gocode/src/covid",
	patternFile:     "config/pattern.csv",
	csvOutputFile:   "config/outreq.csv",
	covidProjectURL: "https://covidtracking.com/api/states/daily",
	templateFiles: []string{"ui/html/base.page.tmpl",
		"ui/html/plot.partial.tmpl"},
	plotFile:  "ui/html/plot.partial.tmpl",
	ipAddress: ":8080",
}

var resultTwo = StatesType{
	appHome:         "",
	patternFile:     "",
	csvOutputFile:   "",
	covidProjectURL: "",
	templateFiles:   []string{},
	plotFile:        "",
	ipAddress:       "",
}

var resultThree = StatesType{
	appHome:         "$HOME/Documents/gocode/src/covid",
	patternFile:     "",
	csvOutputFile:   "config/outreq.csv",
	covidProjectURL: "",
	templateFiles:   []string{"ui/html/base.page.tmpl"},
	plotFile:        "",
	ipAddress:       ":8080",
}

var resultFour = StatesType{
	appHome:         "",
	patternFile:     "config/pattern.csv",
	csvOutputFile:   "",
	covidProjectURL: "https://covidtracking.com/api/states/daily",
	templateFiles:   []string{"ui/html/base.page.tmpl"},
	plotFile:        "ui/html/plot.partial.tmpl",
	ipAddress:       "",
}

var appDir = "/Users/asadolahseghatoleslami/Documents/gocode/src/covid"

//using the LICESE file and config directory as test in the ddirectory
// $GOPATH/covid

func TestFileExists(t *testing.T) {
	dir := "../.." //os.Getenv("GOPATH")
	fileList := []struct {
		fileName string
		expect   bool
	}{
		{
			fileName: filepath.Join(dir, "LICENSE"),
			expect:   true,
		},
		{
			fileName: filepath.Join(dir, "Readme.MD"),
			expect:   true,
		},
		{
			fileName: filepath.Join(dir, "config"),
			expect:   false,
		},
		{
			fileName: filepath.Join(dir, "go.module"),
			expect:   false,
		},
		{
			fileName: filepath.Join(dir, "gitigonre"),
			expect:   false,
		},
	}

	for _, fileItem := range fileList {
		check := fileExists(fileItem.fileName)
		if check != fileItem.expect {
			t.Errorf("testing %s, expected %t got %t", fileItem.fileName,
				fileItem.expect, check)
		}
	}
}

func TestSetUp(t *testing.T) {
	var err error
	pattern := patternType{
		{
			configData: &fileOne,
			configFile: "configTest.csv",
			filePlace:  inDir,
			expect:     noError,
			notError:   &resultOne,
		},
		{
			configData: &fileOne,
			configFile: "configTest.csv",
			filePlace:  inGOPATH,
			expect:     noError,
			notError:   &resultOne,
		},
		{
			configData: &fileOne,
			configFile: "configTest.csv",
			filePlace:  inENV,
			expect:     noError,
			notError:   &resultOne,
		},
		{
			configData: &fileTwo,
			configFile: "configTest.csv",
			filePlace:  inDir,
			expect:     noError,
			notError:   &resultTwo,
		},
		{
			configData: &fileThree,
			configFile: "configTest.csv",
			filePlace:  inGOPATH,
			expect:     noError,
			notError:   &resultThree,
		},
		{
			configData: &fileFour,
			configFile: "configTest.csv",
			filePlace:  inENV,
			expect:     noError,
			notError:   &resultFour,
		},
	}

	for n, item := range pattern {
		s := StatesType{}
		switch item.filePlace {
		case inDir:
			{
				writeFile(filepath.Join(".", item.configFile), item.configData)
				err = s.setUp(item.configFile, "")
				if err != nil {
					t.Errorf("setUp call failed in itermation %d", n)
				}
				err := s.matchUp(item.notError)
				if err != nil {
					t.Errorf("in iteration %d matchUp failed becuse %v", n, err)
				}
				err = os.Remove(filepath.Join(".", item.configFile))
				if err != nil {
					t.Errorf("did not remove %s because %v",
						filepath.Join(".", item.configFile), err)
				}
				s = StatesType{}
			}

		case inGOPATH:
			{
				path := os.Getenv("GOPATH")
				if len(path) == 0 {
					t.Errorf("did not get an environment variable from $GOPATH")
					break
				}
				path = filepath.Join(path, "src/covid")
				err = writeFile(filepath.Join(path, item.configFile), item.configData)
				if err != nil {
					t.Errorf("could not write file %s because %v",
						filepath.Join(path, item.configFile), err)
				}
				err = s.setUp(item.configFile, "")
				if err != nil {
					t.Errorf("setUp failed in iteration %d because %v", n, err)
				}
				err := s.matchUp(item.notError)
				if err != nil {
					t.Errorf("in iteration %d matchUp failed because %v", n, err)
				}
				err = os.Remove(filepath.Join(path, item.configFile))
				if err != nil {
					t.Errorf("did not remove %s because %v",
						filepath.Join(path, item.configFile), err)
				}
				s = StatesType{}
			}

		case inENV:
			{
				env := "TESTPATH"
				err = os.Setenv(env, "../../")
				if err != nil {
					t.Errorf("could not set envrionment variable %s because %v",
						env, err)
				}
				path := os.Getenv(env)
				err = writeFile(filepath.Join(path, item.configFile), item.configData)
				if err != nil {
					t.Errorf("could not write file %s because %v",
						filepath.Join(path, item.configFile), err)
				}
				err = s.setUp(item.configFile, env)
				if err != nil {
					t.Errorf("setUp failed in iteration %d because %v", n, err)
				}
				err := s.matchUp(item.notError)
				if err != nil {
					t.Errorf("in iteration %d matchUp failed because %v", n, err)
				}
				err = os.Remove(filepath.Join(path, item.configFile))
				if err != nil {
					t.Errorf("did not remove %s because %v",
						filepath.Join(path, item.configFile), err)
				}
				s = StatesType{}
			}
		}

		if err != nil {
			t.Errorf("test failed for the final result of %v", err)
		}
	}
}

func (s *StatesType) matchUp(item *StatesType) error {
	if item.appHome != s.appHome {
		return fmt.Errorf("for appHome expected %s got %s",
			item.appHome, s.appHome)
	}
	if item.patternFile != s.patternFile {
		return fmt.Errorf("for patternFile expected %s got %s",
			item.patternFile, s.patternFile)
	}
	if item.csvOutputFile != s.csvOutputFile {
		return fmt.Errorf("for csvOutputFile expted %s got %s",
			item.csvOutputFile, s.csvOutputFile)
	}
	if item.covidProjectURL != s.covidProjectURL {
		return fmt.Errorf("for covidProjectURL expected %s got %s",
			item.covidProjectURL, s.covidProjectURL)
	}
	ok := true
	for _, item2 := range item.templateFiles {
		ok = inSlice(item2, s.templateFiles)
	}
	if !ok {
		return fmt.Errorf("for templateFile expected %v got %v",
			item.templateFiles, s.templateFiles)
	}
	if item.plotFile != s.plotFile {
		return fmt.Errorf("for plotFile expected %s got %s",
			item.plotFile, s.plotFile)
	}
	if item.ipAddress != s.ipAddress {
		return fmt.Errorf("for ipAddress expected %s got %s",
			item.ipAddress, s.ipAddress)
	}
	return nil
}

var newFileZero = `appHome/Users/asadolahseghatoleslami/Documents`

var newFileOne = `appHome|Users/asadolahseghatoleslami/Documents/gocode/src/covid`

var newFileTwo = `appHome|HOME/Documents/gocode/src/covid`

var newFileThree = `appHome|$HOME/Documents/gocode/src/covid
patternFile|config/pattern.csv`

var newFileFour = `appHome|$HOME/Documents/gocode/src/covid
# just a comment
patternFile|config/pattern.csv
csvOutputFile|config/outreq.csv
covidProjectURL|https://covidtracking.com/api/states/daily
templateFile|ui/html/base.page.tmpl
plotFile|ui/html/base.page.tmpl`

var newFileFive = `appHome|$HOME/Documents/gocode/src/covid
# just a comment
patternFile|config/pattern.csv
csvOutputFile|config/outreq.csv
covidProjectURL|https://covidtracking.com/api/states/daily
templateFile|ui/html/base.page.tmpl
templateFile|ui/html/plot.partial.tmpl
plotFile|ui/html/plot.partial.tmpl`

var newFileSix = `appHome|$HOME/Documents/gocode/src/covid
patternFile|config/pattern.csv
csvOutputFile|config/outreq.csv
covidProjectURL|https://covidtracking.com/api/states/daily
templateFile|ui/html/base.page.tmpl
templateFile|ui/html/plot.partial.tmpl
plotFile|ui/html/plot.partial.tmpl
ipAddress|:8080`

var newResult = []StatesType{
	StatesType{ //zero
		appHome:         "",
		patternFile:     "",
		csvOutputFile:   "",
		covidProjectURL: "",
		templateFiles:   []string{},
		plotFile:        "",
		ipAddress:       "",
	},
	StatesType{ //one
		appHome:         "",
		patternFile:     "",
		csvOutputFile:   "",
		covidProjectURL: "",
		templateFiles:   []string{},
		plotFile:        "",
		ipAddress:       "",
	},
	StatesType{ //two
		appHome:         "",
		patternFile:     "",
		csvOutputFile:   "",
		covidProjectURL: "",
		templateFiles:   []string{},
		plotFile:        "",
		ipAddress:       "",
	}, //three
	StatesType{
		appHome:         "/Users/asadolahseghatoleslami/Documents/gocode/src/covid",
		patternFile:     "/Users/asadolahseghatoleslami/Documents/gocode/src/covid/config/pattern.csv",
		csvOutputFile:   "",
		covidProjectURL: "",
		templateFiles:   []string{},
		plotFile:        "",
		ipAddress:       "",
	},
	StatesType{ //four
		appHome:         "/Users/asadolahseghatoleslami/Documents/gocode/src/covid",
		patternFile:     "/Users/asadolahseghatoleslami/Documents/gocode/src/covid/config/pattern.csv",
		csvOutputFile:   "/Users/asadolahseghatoleslami/Documents/gocode/src/covid/config/outreq.csv",
		covidProjectURL: "https://covidtracking.com/api/states/daily",
		templateFiles:   []string{"/Users/asadolahseghatoleslami/Documents/gocode/src/covid/ui/html/base.page.tmpl"},
		plotFile:        "/Users/asadolahseghatoleslami/Documents/gocode/src/covid/ui/html/base.page.tmpl",
		ipAddress:       "",
	},
	StatesType{ //five
		appHome:         "/Users/asadolahseghatoleslami/Documents/gocode/src/covid",
		patternFile:     "/Users/asadolahseghatoleslami/Documents/gocode/src/covid/config/pattern.csv",
		csvOutputFile:   "/Users/asadolahseghatoleslami/Documents/gocode/src/covid/config/outreq.csv",
		covidProjectURL: "https://covidtracking.com/api/states/daily",
		templateFiles: []string{"/Users/asadolahseghatoleslami/ui/html/base.page.tmpl",
			"/Users/asadolahseghatoleslami/Documents/gocode/src/covid/ui/html/plot.partial.tmpl"},
		plotFile:  "/Users/asadolahseghatoleslami/Documents/gocode/src/covid/ui/html/plot.partial.tmpl",
		ipAddress: "",
	},
	StatesType{ //six
		appHome:         "/Users/asadolahseghatoleslami/Documents/gocode/src/covid",
		patternFile:     "/Users/asadolahseghatoleslami/Documents/gocode/src/covid/config/pattern.csv",
		csvOutputFile:   "/Users/asadolahseghatoleslami/Documents/gocode/src/covid/config/outreq.csv",
		covidProjectURL: "https://covidtracking.com/api/states/daily",
		templateFiles: []string{"/Users/asadolahseghatoleslami/Documents/gocode/src/covid/ui/html/base.page.tmpl",
			"/Users/asadolahseghatoleslami/Documents/gocode/src/covid/ui/html/plot.partial.tmpl"},
		plotFile:  "/Users/asadolahseghatoleslami/Documents/gocode/src/covid/ui/html/plot.partial.tmpl",
		ipAddress: ":8080",
	},
}

func TestValidateConfigs(t *testing.T) {
	var err error
	pattern := patternType{
		{
			configData: &newFileZero,
			configFile: "configTest.csv",
			filePlace:  inDir,
			expect:     noError,
			notError:   &newResult[0],
		},
		{
			configData: &newFileOne,
			configFile: "configTest.csv",
			filePlace:  inDir,
			expect:     noError,
			notError:   &newResult[1],
		},
		{
			configData: &newFileTwo,
			configFile: "configTest.csv",
			filePlace:  inGOPATH,
			expect:     noError,
			notError:   &newResult[2],
		},
		{
			configData: &newFileThree,
			configFile: "configTest.csv",
			filePlace:  inENV,
			expect:     noError,
			notError:   &newResult[3],
		},
		{
			configData: &newFileFour,
			configFile: "configTest.csv",
			filePlace:  inDir,
			expect:     noError,
			notError:   &newResult[4],
		},
		{
			configData: &newFileFive,
			configFile: "configTest.csv",
			filePlace:  inGOPATH,
			expect:     noError,
			notError:   &newResult[5],
		},
		{
			configData: &newFileSix,
			configFile: "configTest.csv",
			filePlace:  inENV,
			expect:     noError,
			notError:   &newResult[6],
		},
	}
	for n, item := range pattern {
		s := StatesType{}
		writeFile(filepath.Join(".", item.configFile), item.configData) // TODO: handle error
		s.setUp(item.configFile, "")
		err = s.validateConfigs()
		errString := fmt.Sprintf("%v", err)
		if err != nil {
			switch n {
			case 0:
				if !strings.Contains(errString, "no appHome") {
					t.Errorf("in iteration %d validation error in TestValidateConfigs was %v", n, err)
				}
			case 1:
				if !strings.Contains(errString, "malformed appHome") {
					t.Errorf("in iteration %d validation error in TestValidateConfigs was %v", n, err)
				}
			case 2:
				if !strings.Contains(errString, "malformed appHome") {
					t.Errorf("in iteration %d validation error in TestValidateConfigs was %v", n, err)
				}
			case 3:
				if !strings.Contains(errString, "no csvOutputFile") {
					t.Errorf("in iteration %d validation error in TestValidateConfigs was %v", n, err)
				}
			case 4:
				if !strings.Contains(errString, "no ip address") {
					t.Errorf("in iteration %d validation error in TestValidateConfigs was %v", n, err)
				}
			}
		}
		err = s.matchUp(item.notError)
		switch n {
		case 1:
			if err != nil {
				t.Errorf("test failed in iteration %d of TestValidateConfigs due to %v",
					n, err)
			}
		case 2:
			if err != nil {
				t.Errorf("test failed in iteration %d of TestValidateConfigs due to %v",
					n, err)
			}
		case 3:
			if err != nil {
				t.Errorf("test failed in iteration %d of TestValidateConfigs due to %v",
					n, err)
			}
		case 4:
			if err != nil {
				t.Errorf("test failed in iteration %d of TestValidateConfigs due to %v",
					n, err)
			}
		case 5:
			if err != nil {
				t.Errorf("test failed in iteration %d of TestValidateConfigs due to %v",
					n, err)
			}
		case 6:
			if err != nil {
				t.Errorf("test failed in iteration %d of TestValidateConfigs due to %v",
					n, err)
			}
		}
		err = os.Remove(filepath.Join(".", item.configFile))
		if err != nil {
			t.Errorf("did not remove %s because %v in TestValidateConfigs in iteration %d",
				filepath.Join(".", item.configFile), err, n)
		}
	}
}
