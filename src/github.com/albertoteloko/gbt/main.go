package main

import (
	"github.com/albertoteloko/gbt/log"
	"flag"
	"os"
	//"strings"
)

var version string

var goInterface GoInterface = GoInterfaceFileSystem{}

var debug bool

var (
	GO_PATH    = os.Getenv("GOPATH")
	BIN_FOLDER = GO_PATH + "/bin"
)

var clean bool
var build bool
var test bool
var bench bool
var format bool
var all bool

func main() {
	flag.BoolVar(&debug, "v", false, "Vervose Output")

	//flag.BoolVar(&clean, "c", false, "Execute Clean Task")
	//flag.BoolVar(&build, "b", false, "Execute Build Task")
	//flag.BoolVar(&test, "t", false, "Execute Test Task")
	//flag.BoolVar(&bench, "bt", false, "Execute Benchmark Task")
	//flag.BoolVar(&format, "f", false, "Execute Format Task")
	//
	//flag.BoolVar(&all, "a", false, "Execute All Task")
	//
	//flag.StringVar(&directory, "dir", "", "Base Directory")
	flag.Parse()

	if debug {
		log.Level(log.DEBUG)
	} else {
		log.Level(log.INFO)
	}

	log.Debug("GUT %s", version)

	var definition, err = loadProjectDefinition()

	if err != nil {
		log.Error("Unable to load project definition: %v", err)
	} else {
		var validationErrors = validateDefinition(definition)

		if len(validationErrors) > 0 {
			log.Error("There are %v validation error in gbt.json", len(validationErrors))
			for _, validationError := range validationErrors {
				log.Error("\t* %v", validationError)
			}
		} else {
			log.Info("%v %v", definition.Name, definition.Version)

			err = goInterface.checkAndDownloadGo(definition)

			if err != nil {
				log.Error("Unable to load go version: %v", err)
			}
		}

	}

	//tasks := flag.Args()

	//dir := getBaseDir()
	//log.Debug("Base Dir: %s", dir)
	//
	//logTime("All Tasks", func() {
	//	folders, err := getFolders(dir, isGoFolder)
	//
	//	if err != nil {
	//		log.Error("Error during folder exploration: %s", err)
	//		return
	//	}
	//
	//	if clean {
	//		logTime("Clean", func() {
	//			cleanTask()
	//		})
	//	}
	//
	//	for _, folder := range folders {
	//		logTime("Process "+strings.Replace(folder, GO_PATH+"/src/", "", -1), func() {
	//			processFolder(folder)
	//		})
	//	}
	//})
}

func processFolder(folder string) {
	folderName := getFolderName(folder)

	if (format || all) && isGoFolder(folder) {
		logTime("Format "+folderName, func() {
			formatFolder(folder)
		})
	}

	buildCorrect := true
	if (build || all) && isGoMainFolder(folder) {
		logTime("Build "+folderName, func() {
			buildCorrect = buildTask(folder)
		})
	}

	if buildCorrect {
		testCorrect := true
		if (test || all) && isGoTestFolder(folder) {
			logTime("Test "+folderName, func() {
				testCorrect = testTask(folder)
			})
		}

		if testCorrect {
			if bench {
				logTime("Bench "+folderName, func() {
					benchmarkTask(folder)
				})
			}
		} else {
			log.Fatal("Tests fails")
		}
	} else {
		log.Fatal("Build fails")
	}
}
