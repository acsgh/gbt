package tasks

import (
	"github.com/albertoteloko/gbt/log"
	"strings"
	pd "github.com/albertoteloko/gbt/project-definition"
	"errors"
	"sort"
)

type Tasks []Task

type NotFoundTasks []string

type Task struct {
	Name         string
	Dependencies []string
	priority     uint
	run          func(pd.ProjectDefinition) error
}

type TaskCmd struct {
	Name         string
	Dependencies []string
	priority     uint
	cmd          string
}

func FindTasks(args []string, pd pd.ProjectDefinition) (Tasks, error) {
	tasks, notFoundTasks := getAllTasks(pd).filter(args)
	var err error

	if len(notFoundTasks) > 0 {
		err = errors.New("Tasks not found: " + strings.Join(notFoundTasks, ", "))
	}

	return tasks, err
}

func getAllTasks(pd pd.ProjectDefinition) Tasks {
	tasks := Tasks{
		Task{"clean", []string{}, 0, clean},
		Task{"format", []string{}, 10, format},
		Task{"compile", []string{}, 20, compile},
		Task{"test", []string{}, 30, test},
		Task{"benchmark", []string{}, 31, benchmark},
		Task{"build", []string{"compile", "test"}, 50, nil},
	}

	return tasks
}

func (this Tasks) filter(args []string) (Tasks, NotFoundTasks) {
	tasks := Tasks{}
	notFoundTasks := NotFoundTasks{}

	for _, taskName := range args {
		tasks, notFoundTasks = tasks.addTask(taskName, this, notFoundTasks)
	}

	sort.Sort(tasks)

	return tasks, notFoundTasks
}

func (this Tasks) addTask(taskName string, allTasks Tasks, notFoundTasks NotFoundTasks) (Tasks, []string) {
	tasks := this
	if !strings.HasPrefix(taskName, "-") {
		task := allTasks.find(taskName)

		if task != nil {
			tasks = append(tasks, *task)

			for _, childTaskName := range task.Dependencies {
				tasks, notFoundTasks = tasks.addTask(childTaskName, allTasks, notFoundTasks)
			}
		} else {
			notFoundTasks = notFoundTasks.add(taskName)
		}
	}
	return tasks, notFoundTasks
}

func (slice Tasks) Len() int {
	return len(slice)
}

func (slice Tasks) Less(i, j int) bool {
	return slice[i].priority < slice[j].priority;
}

func (slice Tasks) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (this NotFoundTasks) add(taskName string) NotFoundTasks {
	if !this.contains(taskName) {
		this = append(this, taskName)
	}
	return this
}

func (this NotFoundTasks) contains(taskName string) bool {
	for _, task := range this {
		if strings.ToLower(task) == strings.ToLower(taskName) {
			return true
		}

	}
	return false
}

func (this Tasks) find(taskName string) *Task {
	for _, task := range this {
		if strings.ToLower(task.Name) == strings.ToLower(taskName) {
			return &task
		}
	}
	return nil
}

func (this Tasks) Run(definition pd.ProjectDefinition) error {
	var err error
	for _, task := range this {

		if task.run != nil {
			err = log.LogTimeWithError("Task "+task.Name, func() error {
				return task.run(definition)
			})

			if err != nil {
				log.Errorf("Error in task %v: %v", task.Name, err)
				break
			}
		}
	}
	return err
}