package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"morgansundqvist/todin/internalTypes"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

func CreateTask(tasks *[]internalTypes.Task) {

	var qs = []*survey.Question{
		{
			Name: "title",
			Prompt: &survey.Input{
				Message: "What is the task?",
			},
			Validate: survey.Required,
		},
		{
			Name: "priority",
			Prompt: &survey.Input{
				Message: "What is the priority?",
			},
			Validate: survey.Required,
		},
	}
	answers := struct {
		Title    string
		Priority int
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	*tasks = append(*tasks, internalTypes.Task{Id: FindHighestId(*tasks) + 1, Title: answers.Title, Priority: answers.Priority, IsDone: false, CreatedDateTime: time.Now()})
}

func ListTasksToDo(tasks *[]internalTypes.Task) {
	SortTasksByPriority(tasks)
	for _, data := range *tasks {
		if !data.IsDone {
			fmt.Printf("%d\t%d\t%s\n", data.Id, data.Priority, data.Title)
		}
	}
}

func ListTasksDone(tasks *[]internalTypes.Task) {
	for _, data := range *tasks {
		if data.IsDone {
			fmt.Printf("%d\t%d\t%s\n", data.Id, data.Priority, data.Title)
		}
	}
}

// Function that takes a slice of tasks, asks for which id to remove and removes it
func RemoveTask(tasks *[]internalTypes.Task) {
	var qs = []*survey.Question{
		{
			Name: "id",
			Prompt: &survey.Input{
				Message: "What is the id of the task you want to remove?",
			},
			Validate: survey.Required,
		},
	}
	answers := struct {
		Id int
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	index := FindTaskIndex(*tasks, answers.Id)
	if index == -1 {
		fmt.Println("No task with that id found")
		return
	}

	*tasks = append((*tasks)[:index], (*tasks)[index+1:]...)
}

func MarkTaskAsDone(tasks *[]internalTypes.Task) {
	var qs = []*survey.Question{
		{
			Name: "id",
			Prompt: &survey.Input{
				Message: "What is the id of the task you want to mark as done?",
			},
			Validate: survey.Required,
		},
	}
	answers := struct {
		Id int
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	index := FindTaskIndex(*tasks, answers.Id)
	if index == -1 {
		fmt.Println("No task with that id found")
		return
	}

	(*tasks)[index].IsDone = true
}

// Function that takes a slice of tasks and changes the priority of a task with a given id
func ChangePriority(tasks *[]internalTypes.Task) {
	var qs = []*survey.Question{
		{
			Name: "id",
			Prompt: &survey.Input{
				Message: "What is the id of the task you want to change the priority of?",
			},
			Validate: survey.Required,
		},
		{
			Name: "priority",
			Prompt: &survey.Input{
				Message: "What is the new priority?",
			},
			Validate: survey.Required,
		},
	}
	answers := struct {
		Id       int
		Priority int
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	index := FindTaskIndex(*tasks, answers.Id)
	if index == -1 {
		fmt.Println("No task with that id found")
		return
	}

	(*tasks)[index].Priority = answers.Priority
}

// Find the index of a task with a given id in a slice of tasks
func FindTaskIndex(tasks []internalTypes.Task, id int) int {
	for i, data := range tasks {
		if data.Id == id {
			return i
		}
	}
	return -1
}

// Function that takes tasks slice and finds the highest id and returns it
func FindHighestId(tasks []internalTypes.Task) int {
	var highestId int = 0
	for _, data := range tasks {
		if data.Id > highestId {
			highestId = data.Id
		}
	}
	return highestId
}

// Function that takes the slice and stores it in a json file
func SaveTasks(tasks *[]internalTypes.Task) {
	file, _ := json.MarshalIndent(tasks, "", " ")
	_ = ioutil.WriteFile("tasks.json", file, 0644)
}

// Function that resorts the slice by priority
func SortTasksByPriority(tasks *[]internalTypes.Task) {
	for i := 0; i < len(*tasks); i++ {
		for j := 0; j < len(*tasks)-1; j++ {
			if (*tasks)[j].Priority > (*tasks)[j+1].Priority {
				(*tasks)[j], (*tasks)[j+1] = (*tasks)[j+1], (*tasks)[j]
			}
		}
	}
}
