package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"morgansundqvist/todin/handler"
	"morgansundqvist/todin/internalTypes"
	"os"
	"os/exec"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
)

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func main() {
	// Path to the JSON file
	filePath := "./tasks.json"

	// Try to open the file
	file, err := os.Open(filePath)
	if err != nil {
		// If the file doesn't exist, create it
		if os.IsNotExist(err) {
			file, err = os.Create(filePath)
			if err != nil {
				panic(err)
			}
			// Write an empty JSON object to the file
			emptyJSON := []byte("[]")
			err = ioutil.WriteFile(filePath, emptyJSON, 0644)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	defer file.Close()

	// Read the contents of the file
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON into a struct
	tasks := []internalTypes.Task{}
	err = json.Unmarshal(fileContents, &tasks)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	var mainQs = []*survey.Question{
		{
			Name: "actionToExecute",
			Prompt: &survey.Select{
				Message: "Choose an action:",
				Options: []string{"create task", "list tasks to do", "mark task as done", "change priority", "list tasks done", "remove task", "clear", "exit"},
				Default: "create task",
			},
		},
	}

	answers := struct {
		ActionToExecute string
	}{}

	for {
		fmt.Println("")
		err = survey.Ask(mainQs, &answers, survey.WithPageSize(10))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if answers.ActionToExecute == "create task" {
			//call createTask in handler package
			handler.CreateTask(&tasks)
		} else if answers.ActionToExecute == "list tasks to do" {
			handler.ListTasksToDo(&tasks)
			fmt.Println()
		} else if answers.ActionToExecute == "mark task as done" {
			handler.MarkTaskAsDone(&tasks)
		} else if answers.ActionToExecute == "list tasks done" {
			handler.ListTasksDone(&tasks)
		} else if answers.ActionToExecute == "remove task" {
			handler.RemoveTask(&tasks)
		} else if answers.ActionToExecute == "change priority" {
			handler.ChangePriority(&tasks)
		} else if answers.ActionToExecute == "clear" {
			CallClear()
		} else if answers.ActionToExecute == "exit" {
			break
		} else {
			fmt.Println("Invalid action")
		}
		handler.SaveTasks(&tasks)

	}

}
