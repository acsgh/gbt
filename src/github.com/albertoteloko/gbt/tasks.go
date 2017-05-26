package main

import (
	"github.com/albertoteloko/gbt/log"
	"github.com/albertoteloko/gbt/file"
	"github.com/albertoteloko/gbt/utils"
	"io/ioutil"
	"os/exec"
	"strings"
	pd "github.com/albertoteloko/gbt/project-definition"
)

type Tasks [] Task


type Task struct {
	Name         string
	Dependencies []string
	priority     uint
	run       func(pd.ProjectDefinition) error
}


type TaskCmd struct {
	Name         string
	Dependencies []string
	priority     uint
	cmd          string
}

var tasks Tasks = Tasks{
	Task{"clean", []string{}, 0, clean},
}

func (this Tasks) findTasks(args Args, pd pd.ProjectDefinition) (Tasks, error) {
	return this, nil
}

func (this Tasks) run(definition pd.ProjectDefinition) error {
	var err error
	for _, task := range this {
		err = log.LogTimeWithError("Task " + task.Name, func() error {
			return task.run(definition)
		})

		if err != nil {
			log.Errorf("Error in task %v: %v", task.Name, err)
			break
		}
	}
	return err
}

func clean(pd pd.ProjectDefinition) error {
	//os.RemoveAll(BIN_FOLDER)
	//os.MkdirAll(BIN_FOLDER, 0777)

	return nil
}

func formatFolder(dir string) {
	files, err := file.GetFiles(dir, file.And(file.IsGoFile, file.IsGitHubPath))

	if err != nil {
		log.Errorf("Errorf during folder read: %s", err)
	}

	for _, f := range files {
		log.Debugf("Formating file: %s", f)
		if err := formatFile(f); err != nil {
			log.Errorf("Errorf formating file %s, %s", f, err)
		} else {
			log.Debugf("Formated file: %s", f)
		}
	}
}

func formatFile(file string) error {
	return exec.Command("gofmt", "-w", file).Run()
}

func buildTask(folder string) bool {
	_, err := buildFolder(getFolderName(folder))

	if err != nil {
		log.Errorf("Errorf during folder building: %s", err)
	}
	return (err == nil)
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
		log.Errorf("Errorf during folder testing: %s", err)
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
		log.Errorf("Errorf during folder benchmarking: %s", err)
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
