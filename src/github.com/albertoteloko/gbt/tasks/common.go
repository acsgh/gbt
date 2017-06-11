package tasks

import (
	"github.com/albertoteloko/gbt/log"
	"github.com/albertoteloko/gbt/file"
	"github.com/albertoteloko/gbt/utils"
	"io/ioutil"
	"os/exec"
	"strings"
)

func buildTask(folder string) bool {
	_, err := buildFolder(getFolderName(folder))

	if err != nil {
		log.Errorf("Error during folder building: %s", err)
	}
	return err == nil
}

func buildFolder(folder string) (string, error) {
	command := exec.Command("go", "build", "-o", file.GO_PATH+"/bin/"+getSourceFolderName(folder), folder)
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
		log.Errorf("Error during folder testing: %s", err)
	} else {
		namesLength := 0
		for _, result := range results {
			namesLength = utils.Max(namesLength, len(result.name))
		}
		for _, testResult := range results {
			testCorrect = testCorrect && (testResult.result == "PASS")
			log.Infof("%s %s [%s]", utils.FixWidth(testResult.name, namesLength+3), testResult.result, testResult.time)
		}
		log.Infof("Coverage: %s", coverage)
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
				parts := strings.Split(utils.ReplaceChars(utils.ReplaceChars(line, "", " ", "---", ")"), "\t", ":", "("), "\t")
				testOutputs = append(testOutputs, testResult{parts[1], parts[0], parts[2]})
			} else if strings.HasPrefix(line, "coverage: ") {
				coverage = utils.ReplaceChars(line, "", "coverage: ", " of statements")
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
		log.Errorf("Error during folder benchmarking: %s", err)
	} else {
		namesLength := 0
		for _, result := range results {
			namesLength = utils.Max(namesLength, len(result.name))
		}
		for _, result := range results {
			log.Infof("%s %s [%s]", utils.FixWidth(result.name, namesLength+3), result.rate, result.times)
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
				parts := strings.Split(utils.ReplaceChars(line, "", " ", "---", ")"), "\t")
				testOutputs = append(testOutputs, benchmarkResult{parts[0], parts[1], parts[2]})
			}
		}

		return testOutputs, nil
	}
}

func getFolderName(folder string) string {
	return file.GetFolderName(file.GO_PATH+"/src/", folder)
}
