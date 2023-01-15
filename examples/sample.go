package main

import (
	"fmt"
	//"github.com/hideshi/goawk"
	"goawk"
	"strconv"
)

func Begin(app *goawk.App) {
	fmt.Println("Execute Begin")
	fmt.Println(app.Filename)
}

func Action1(app *goawk.App) {
	fmt.Println("Execute Action1")
	fmt.Println("Input file is processed one by one per line, then the text is splitted with field separator")
	for _, elem := range app.S[1:] {
		fmt.Printf("%#v\n", elem)
		v, err := strconv.Atoi(elem)
		if err != nil {
			fmt.Println(err)
		} else {
			app.VI["sum"] = app.VI["sum"] + v
		}
	}
}

func Action2(app *goawk.App) {
	fmt.Println("Execute Action2")
    pattern := ".*20.*"
	matched, _ := app.P(pattern)
	if matched == true {
		fmt.Printf("The text %s matched with the patten %s\n", app.S[0], pattern)
	}
}

func End(app *goawk.App) {
	fmt.Println("Execute End")
	fmt.Printf("Sum of the input is %d\n", app.VI["sum"])
}

func main() {
	app := new(goawk.App)
	actions := []goawk.Action{Begin, Action1, Action2, End}
	app.Run(actions)
}
