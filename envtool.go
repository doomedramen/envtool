package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	importOption  bool
	exportOption  bool
	verboseOption = false
	pwd           string
)

const (
	fileName = "ENV.env"
)

type envObj struct {
	name  string
	value string
}

func main() {
	//set default path
	tmpPwd, err := os.Getwd()
	check(err)
	pwd = tmpPwd

	processArgs()

	if importOption {
		//check full path exists
		pwd = path.Join(pwd, fileName)
		checkForFile()

		importer()
	} else if exportOption {
		//check folder exists
		checkForFile()
		pwd = path.Join(pwd, fileName)

		exporter()
	}
}

func processArgs() {

	verboseFlag := flag.Bool("verbose", false, "prints more output")
	importFlag := flag.Bool("import", false, "import env")
	exportFlag := flag.Bool("export", false, "export env")
	folder := flag.String("path", "", "[optional] path to ENV.env file")

	flag.Parse()

	if *verboseFlag {
		verboseOption = *verboseFlag
	}
	if *importFlag {
		importOption = *importFlag
	} else if *exportFlag {
		exportOption = *exportFlag
	}

	if flag.NFlag() < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *folder != "" {
		pwd = *folder
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkForFile() {
	if verboseOption {
		fmt.Println("using", pwd, "for file location")
	}

	if _, err := os.Stat(pwd); os.IsNotExist(err) {
		fmt.Println("no such file or directory", pwd)
		os.Exit(1)
	}
}

func exporter() {

	envs := getEnv()
	mapB, _ := json.Marshal(envs)

	err := ioutil.WriteFile(pwd, mapB, 0644)
	check(err)

	fmt.Println("exported", len(envs), "variables to", pwd)

}

func importer() {

	filebyte, err := ioutil.ReadFile(pwd)
	check(err)

	// var deps []envObj

	var dat map[string]interface{}
	err = json.Unmarshal(filebyte, &dat)
	check(err)

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
	fmt.Println("imported", len(dat), "variables to", pwd)
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
