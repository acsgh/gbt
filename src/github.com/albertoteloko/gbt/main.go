package main

import (
	"github.com/albertoteloko/gbt/log"
	"github.com/albertoteloko/gbt/file"
	"os"
	pd "github.com/albertoteloko/gbt/project-definition"
	gi "github.com/albertoteloko/gbt/go-interface"
)

var version string

var goInterface gi.GoInterface = gi.GoInterfaceFileSystem{}
var projectDefinitionLoader pd.ProjectDefinitionLoader = pd.ProjectDefinitionLoaderFileSystem{}

func main() {
	configureLogs()

	if isHelp(){
		printHelp()
		os.Exit(1)
	}

	log.LogTime("GBT "+version, func() {
		initTool()

		var definition, err = projectDefinitionLoader.Load()

		if err != nil {
			log.Error("Unable to load project definition: %v", err)
		} else {
			var validationErrors = definition.Validate()

			if len(validationErrors) > 0 {
				log.Error("There are %v validation error in gbt.json", len(validationErrors))
				for _, validationError := range validationErrors {
					log.Error("\t* %v", validationError)
				}
			} else {
				run(definition)
			}

		}
	})
	//tasks := flag.Args()

	//dir := getBaseDir()
	//log.Debug("Base Dir: %s", dir)
	//
	//log.LogTime("All Tasks", func() {
	//	folders, err := getFolders(dir, isGoFolder)
	//
	//	if err != nil {
	//		log.Error("Error during folder exploration: %s", err)
	//		return
	//	}
	//
	//	if clean {
	//		log.LogTime("Clean", func() {
	//			cleanTask()
	//		})
	//	}
	//
	//	for _, folder := range folders {
	//		log.LogTime("Process "+strings.Replace(folder, GO_PATH+"/src/", "", -1), func() {
	//			processFolder(folder)
	//		})
	//	}
	//})
}

func initTool() {
	//os.RemoveAll(BIN_FOLDER)
	os.MkdirAll(file.GBT_PATH, 0775)
	os.MkdirAll(file.GO_DISTRO_PATH, 0775)
	os.MkdirAll(file.GBT_DISTROS_PATH, 0775)
	os.MkdirAll(file.DEPENDENCIES_PATH, 0775)
	os.MkdirAll(file.TMP_PATH, 0775)
}

func run(definition pd.ProjectDefinition) {
	log.Info("Project: %v %v", definition.Name, definition.Version)
	var err = goInterface.CheckAndDownloadGo(definition)
	if err != nil {
		log.Error("Unable to load go version: %v", err)
	}
}
//
//func processFolder(folder string) {
//	folderName := getFolderName(folder)
//
//	if (format || all) && file.IsGoFolder(folder) {
//		log.LogTime("Format "+folderName, func() {
//			formatFolder(folder)
//		})
//	}
//
//	buildCorrect := true
//	if (build || all) && file.IsGoMainFolder(folder) {
//		log.LogTime("Build "+folderName, func() {
//			buildCorrect = buildTask(folder)
//		})
//	}
//
//	if buildCorrect {
//		testCorrect := true
//		if (test || all) && isGoTestFolder(folder) {
//			log.LogTime("Test "+folderName, func() {
//				testCorrect = testTask(folder)
//			})
//		}
//
//		if testCorrect {
//			if bench {
//				log.LogTime("Bench "+folderName, func() {
//					benchmarkTask(folder)
//				})
//			}
//		} else {
//			log.Fatal("Tests fails")
//		}
//	} else {
//		log.Fatal("Build fails")
//	}
//}
