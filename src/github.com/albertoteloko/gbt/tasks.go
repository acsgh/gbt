package main

import (
	"github.com/albertoteloko/gbt/log"
	"github.com/albertoteloko/gbt/file"
	"io/ioutil"
	"os/exec"
	"strings"
)

func cleanTask() {
	//os.RemoveAll(BIN_FOLDER)
	//os.MkdirAll(BIN_FOLDER, 0777)
}

func formatFolder(dir string) {
	files, err := file.GetFiles(dir, file.And(file.IsGoFile, file.IsGitHubPath))

	if err != nil {
		log.Error("Error during folder read: %s", err)
	}

	for _, f := range files {
		log.Debug("Formating file: %s", f)
		if err := formatFile(f); err != nil {
			log.Error("Error formating file %s, %s", f, err)
		} else {
			log.Debug("Formated file: %s", f)
		}
	}
}

func formatFile(file string) error {
	return exec.Command("gofmt", "-w", file).Run()
}

func buildTask(folder string) bool {
	_, err := buildFolder(getFolderName(folder))

	if err != nil {
		log.Error("Error during folder building: %s", err)
	}
	return (err == nil)
}

func buildFolder(folder string) (string, error) {
	command := exec.Command("go", "build", "-o", GO_PATH+"/bin/"+getSourceFolderName(folder), folder)
	bytes, err := command.Output()
	return string(bytes), err
}

func getSourceFolderName(folderPath string) string {
	if strings.Contains(folderPath, "/") {
		return folderPath[strings.LastIndex(folderPath, "/")+1:]
	} else {
		return folderPath
	}
}

type testResult struct {
	name   string
	result string
	time   string
}

func testTask(folder string) bool {
	testCorrect := true

	folderName := getFolderName(folder)
	results, coverage, err := testFolder(folderName)

	if err != nil {
		log.Error("Error during folder testing: %s", err)
	} else {
		namesLength := 0
		for _, result := range results {
			namesLength = max(namesLength, len(result.name))
		}
		for _, testResult := range results {
			testCorrect = testCorrect && (testResult.result == "PASS")
			log.Info("%s %s [%s]", fixWidth(testResult.name, namesLength+3), testResult.result, testResult.time)
		}
		log.Info("Coverage: %s", coverage)
	}
	return testCorrect
}

func testFolder(folder string) ([]testResult, string, error) {
	command := exec.Command("go", "test", folder, "-v", "-cover")
	bytes, err := command.Output()

	output := string(bytes)

	if !strings.Contains(output, "FAIL") && err != nil {
		return []testResult{}, "", err
	} else {
		testOutputs := []testResult{}
		coverage := "0.0%"
		for _, line := range strings.Split(output, "\n") {
			if strings.HasPrefix(line, "--- ") {
				parts := strings.Split(replaceChars(replaceChars(line, "", " ", "---", ")"), "\t", ":", "("), "\t")
				testOutputs = append(testOutputs, testResult{parts[1], parts[0], parts[2]})
			} else if strings.HasPrefix(line, "coverage: ") {
				coverage = replaceChars(line, "", "coverage: ", " of statements")
			}

		}

		return testOutputs, coverage, nil
	}
}

func isGoTestFolder(name string) bool {
	files, err := ioutil.ReadDir(name)

	if err != nil {
		return false
	}

	for _, f := range files {
		fileName := name + "/" + f.Name()
		if !f.IsDir() && file.IsGoTestFile(fileName) {
			return true
		}
	}
	return false
}

type benchmarkResult struct {
	name  string
	times string
	rate  string
}

func benchmarkTask(folder string) {
	folderName := getFolderName(folder)

	results, err := benchmarkFolder(folderName)

	if err != nil {
		log.Error("Error during folder benchmarking: %s", err)
	} else {
		namesLength := 0
		for _, result := range results {
			namesLength = max(namesLength, len(result.name))
		}
		for _, result := range results {
			log.Info("%s %s [%s]", fixWidth(result.name, namesLength+3), result.rate, result.times)
		}
	}
}

func benchmarkFolder(folder string) ([]benchmarkResult, error) {
	command := exec.Command("go", "test", folder, "-bench", ".", "-run", "NoOneWillFitThisDescription")
	bytes, err := command.Output()

	output := string(bytes)

	if !strings.Contains(output, "FAIL") && err != nil {
		return []benchmarkResult{}, err
	} else {
		testOutputs := []benchmarkResult{}

		for _, line := range strings.Split(output, "\n") {
			if strings.HasPrefix(line, "Benchmark") {
				parts := strings.Split(replaceChars(line, "", " ", "---", ")"), "\t")
				testOutputs = append(testOutputs, benchmarkResult{parts[0], parts[1], parts[2]})
			}
		}

		return testOutputs, nil
	}
}

func getFolderName(folder string) string {
	return file.GetFolderName(GO_PATH+"/src/", folder)
}