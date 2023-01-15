package main

import (
	"fmt"
	//"github.com/hideshi/goawk"
	"goawk"
)

func Action1(app *goawk.App) {
//	app.VI[app.S[0]] = app.VI[app.S[0]] + 1
    fmt.Printf("%s\n", app.S)
    for _, v := range(app.S) {
	app.VI[v] = app.VI[v] + 1
    }
}

func End(app *goawk.App) {
    for k, v := range(app.VI) {
	    fmt.Printf("%v : %v\n", k, v)
    }
}

func main() {
	//app := new(goawk.App)
	//app.Fs = " "
	//app.FsRegex = "\\s+"
	//app.MaxLines = 0

	app := goawk.NewApp()
	//app.MaxLines = 3
	//app.SetLineExtractionCondition("Ora")
	//app.SetLineExtractionCondition("Apple")
	//app.SetLineExtractionCondition("Pine","Banana")
	//app.SetLineExtractionCondition("Banana")

	//app.SetFieldExtractionCondition(map[int]string {0: "^A.*",1:"B$"})
	//app.SetFieldExtractionCondition(map[int]string {0: "BBBB",1:"CCCC"})

	//app.SetFieldStringPickUpRegex(map[int]string {0: "^(..)\\#.*",1:".+(..)$"})
	app.SetFieldStringPickUpRegex(map[int]string {0: "^(..)\\#.*"})

	actions := []goawk.Action{Action1, End}
	app.Run(actions)
}
