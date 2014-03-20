package main

//TODO
// check the OS
// add ability to give path to (in/out) file
// 

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	importOption  bool
	exportOption  bool
	verboseOption bool
)

const (
	fileName = "ENV.env"
)

type envObj struct {
	name  string
	value string
}

func main() {

	argsWithoutProg := os.Args[1:]

	for _, element := range argsWithoutProg {

		if strings.Contains(element, "verbose") {
			verboseOption = true
		}

		if strings.Contains(element, "import") {
			importOption = true
			importer()
		} else if strings.Contains(element, "export") {
			exportOption = true
			exporter()
		}

	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func exporter() {

	envs := getEnv()
	mapB, _ := json.Marshal(envs)

	err := ioutil.WriteFile(fileName, mapB, 0644)
	check(err)

}

func importer() {

	filebyte, err := ioutil.ReadFile(fileName)
	check(err)

	// var deps []envObj

	var dat map[string]interface{}
	err = json.Unmarshal(filebyte, &dat)
	check(err)

	fmt.Println("importing", len(dat), "variables")

	for k, v := range dat {

		if str, ok := v.(string); ok {
			/* act on str */
			if verboseOption {
				fmt.Println("SETTING", k, "=", str)
			}
			os.Setenv(k, str)

		} else {
			/* not string */
		}

	}
	fmt.Println("DONE")
}

func getEnv() map[string]string {
	items := make(map[string]string)
	for _, item := range os.Environ() {
		key, val := getKeyVal(item)
		items[key] = val

	}
	return items
}

func getKeyVal(item string) (key, val string) {
	splits := strings.Split(item, "=")
	key = splits[0]
	val = strings.Join(splits[1:], "=")
	return
}
