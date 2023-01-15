package goawk

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

type App struct {
	Filename                string // Input file name
	Fs                      string // Field separator
	FsRegex                 string // Field separator
	MaxLines                int    //
	LineExtractionCondition bool
	//LineExtractionConditionRegex string
	//LineExtractionConditionRegexCompile *regexp.Regexp
	//LineExtractionConditionRegexArray []string
	//LineExtractionConditionRegexCompileArray []*regexp.Regexp
	LineExtractionConditionRegex2Array        [][]string
	LineExtractionConditionRegexCompile2Array [][]*regexp.Regexp

	FieldExtractionCondition             bool
	FieldExtractionConditionMap          []map[int]string
	FieldExtractionConditionRegexCompile []map[int]*regexp.Regexp

	FieldStringPickUp             bool
	FieldStringPickUpRegexCompile map[int]*regexp.Regexp

	S  []string          // Input line
	VS map[string]string // Variables for string value
	VI map[string]int    // Variables for int value
}

type Action func(app *App)

func NewApp() *App {
	//array := make([][]*regexp.Regexp,5, 5)
	//array [][]*regexp.Regexp
	return &App{
		Fs:                       " ",
		FsRegex:                  "\\s+",
		MaxLines:                 0,
		LineExtractionCondition:  false,
		FieldExtractionCondition: false,
		FieldStringPickUp:        false,
	}
}

func (app *App) Run(actions []Action) {
	errLogger := log.New(os.Stderr, "", 0)
	var fileName string

	flag.StringVar(&fileName, "i", "", "Input file name")
	flag.Parse()

	if len(fileName) == 0 {
		//fmt.Fprintf(os.Stderr, "missing required -i\n")
		flag.Usage()
		os.Exit(2) // the same exit code flag.Parse uses
	}
	app.Filename = fileName
	//app.Fs = ","
	app.VS = make(map[string]string)
	app.VI = make(map[string]int)

	firstActionName := runtime.FuncForPC(reflect.ValueOf(actions[0]).Pointer()).Name()
	if firstActionName == "main.Begin" {
		actions[0](app)
		actions = actions[1:]
	}

	lengthOfActions := len(actions)
	lastActionName := runtime.FuncForPC(reflect.ValueOf(actions[lengthOfActions-1:][0]).Pointer()).Name()
	var endAction Action

	if lastActionName == "main.End" {
		endAction = actions[lengthOfActions-1:][0]
		actions = actions[:lengthOfActions-1]
	}

	input, err := os.Open(fileName)
	defer input.Close()
	if err != nil {
		errLogger.Println("Input file does not exist.")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(input)
	var linecount int
	for scanner.Scan() {
		//fmt.Fprintf(os.Stderr, "SCAN->%s\n", scanner.Text())
		linecount++
		if app.MaxLines > 0 && linecount > app.MaxLines {
			break
		}
		if scanner.Text() == "" {
			continue
		}
		if app.LineExtractionCondition {
			match := false
			for _, v := range app.LineExtractionConditionRegexCompile2Array {
				match2 := false
				for _, v2 := range v {
					if v2.MatchString(scanner.Text()) {
						match2 = true
					} else {
						match2 = false
						break
					}
				}
				if match2 {
					match = true
					break
				} else {
					match = false

				}
			}
			if !match {
				continue
			}
		}
		if app.FieldExtractionCondition {
			field_array := regexp.MustCompile(app.FsRegex).Split(strings.TrimRight(scanner.Text(), " "), -1)
			//fmt.Printf("--%s\n", field_array)

			match := true
			for _, v := range app.FieldExtractionConditionRegexCompile {
				match2 := true
				for k2, v2 := range v {
					if v2.MatchString(field_array[k2]) {
						match2 = true
					} else {
						match2 = false
						break
					}
				}
				if !match2 {
					match = false
				} else {
					match = true
					break
				}
			}
			if !match {
				//fmt.Printf("no match\n")
				continue
			} else {
				//fmt.Printf("match\n")
			}
		}

		app.S = nil
		fmt.Printf("=>%s\n", scanner.Text())
		for index, elem := range regexp.MustCompile(app.FsRegex).Split(strings.TrimRight(scanner.Text(), " "), -1) {
			if !app.FieldStringPickUp {
				app.S = append(app.S, elem)
			} else {
				val, ok := app.FieldStringPickUpRegexCompile[index]
				if ok {
					//fmt.Print("regex sub\n")
					result := val.FindAllStringSubmatch(elem, -1)
					if len(result) > 0 {
						//fmt.Print(" sub result ok\n")
						app.S = append(app.S, result[0][1])
					} else {
						//fmt.Print(" sub result ng\n")
						app.S = append(app.S, elem)
					}
				} else {
					app.S = append(app.S, elem)
				}
			}
		}
		for _, action := range actions {
			action(app)
		}
	}

	endAction(app)
}

func (app *App) P(pattern string) (bool, error) {
	return regexp.MatchString(pattern, app.S[0])
}

//func (app *App) SetLineExtractionCondition(pattern string)  {
func (app *App) SetLineExtractionCondition(pattern_array ...string) {
	app.LineExtractionCondition = true
	var array []*regexp.Regexp
	var array2 []string
	for _, pattern := range pattern_array {
		array = append(array, regexp.MustCompile(pattern))
		array2 = append(array2, pattern)
	}
	app.LineExtractionConditionRegexCompile2Array = append(app.LineExtractionConditionRegexCompile2Array, array)
	app.LineExtractionConditionRegex2Array = append(app.LineExtractionConditionRegex2Array, array2)
}

func (app *App) SetFieldExtractionCondition(pattern_map map[int]string) {
	app.FieldExtractionCondition = true
	var dic = map[int]*regexp.Regexp{}
	for fn, regex := range pattern_map {
		//fmt.Printf("%d\n", fn)
		//fmt.Printf("%s\n", regex)
		//dic = append(dic, map[int]*regexp.Regexp {fn: regexp.MustCompile(regex)})
		dic[fn] = regexp.MustCompile(regex)
	}
	app.FieldExtractionConditionRegexCompile = append(app.FieldExtractionConditionRegexCompile, dic)
	//fmt.Print(app.FieldExtractionConditionRegexCompile )
}

func (app *App) SetFieldStringPickUpRegex(pattern_map map[int]string) {
	app.FieldStringPickUp = true
	var dic = map[int]*regexp.Regexp{}
	for fn, regex := range pattern_map {
		//fmt.Printf("%d\n", fn)
		//fmt.Printf("%s\n", regex)
		//dic = append(dic, map[int]*regexp.Regexp {fn: regexp.MustCompile(regex)})
		dic[fn] = regexp.MustCompile(regex)
	}
	//app.FieldStringPickUpRegexCompile = append(app.FieldStringPickUpRegexCompile, dic)
	app.FieldStringPickUpRegexCompile = dic
	//fmt.Print(app.FieldExtractionConditionRegexCompile )
}
