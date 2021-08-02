package main

import (
	// "bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"log"
	"os"
	// "strings"
)

type Lists struct {
	Todo      []string `json:"todo"`
	Completed []string `json:"completed"`
}

var (
	lists        Lists
	todoList     []string
	completeList []string
	addFlag      = flag.String("a", "", "Add a todo")
	completeFlag = flag.Int("c", 0, "Complete a todo")
	helpFlag     = flag.Bool("h", false, "Help")
	verboseFlag  = flag.Bool("v", false, "Verbose")
	removeFlag   = flag.Int("r", 0, "Remove a todo")
	clearFlag    = flag.Bool("clear", false, "Clear all todos")
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func valid(index int) bool {
	if index < 0 || index >= len(todoList) {
		print("invalid index\n")
		return false
	}
	return true
}

func Add(todo string, file *os.File) {
	todoList = append(todoList, todo)
	print("added '", todo, "'\n")
	List()
	update(file)
}

func Remove(index int, file *os.File) {
	index--
	if valid(index) {
		todoList = append(todoList[:index], todoList[index+1:]...)
		update(file)
	}
}

func Complete(index int, file *os.File) {
	index--
	if valid(index) {
		completeList = append(completeList, todoList[index])
		todoList = append(todoList[:index], todoList[index+1:]...)
		update(file)
		ListFull()
	}
}

func List() {
	print("Todo:\n")
	for i := 0; i < len(todoList); i++ {
		print(i+1, ". "+todoList[i]+"\n")
	}
}

func ListFull() {
	List()
	print("\nComplete:\n")
	for i := 0; i < len(completeList); i++ {
		print(i+1, ". "+completeList[i]+"\n")
	}
}

func Clear(file *os.File) {
	label := "Are you sure you want to delete all list items?\n"
	clear := yesNo(label)
	if !clear {
		print("No items cleared.\n")
		return
	}
	todoList = []string{}
	completeList = []string{}
	update(file)
	print("Cleared all items.\n")
}

func Help() {
	fmt.Println("Usage: todo [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func yesNo(label string) bool {
	prompt := promptui.Select{
		Label: label,
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

func update(file *os.File) {
	os.Truncate(file.Name(), 0)
	file.Seek(0, 0)
	lists.Todo = todoList
	lists.Completed = completeList

	byteValue, _ := json.Marshal(lists)
	_, err := file.Write(byteValue)
	check(err)
}

func main() {

	// Open file if it exists or create a new one
	file, err := os.OpenFile("list.json", os.O_RDWR|os.O_CREATE, 0644)
	check(err)

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &lists)

	todoList = lists.Todo
	completeList = lists.Completed

	flag.Parse()

	if flag.NArg() > 1 {
		print("Too many arguments\n")
		print("Try 'todo -h' for more information\n")
		return
	}

	if *clearFlag {
		Clear(file)
		return
	}

	if *helpFlag {
		Help()
		return
	}

	if *verboseFlag {
		ListFull()
		return
	}

	if *addFlag != "" {
		Add(*addFlag, file)
		return
	}

	if *removeFlag > 0 {
		Remove(*completeFlag, file)
		return
	}

	if *completeFlag > 0 {
		Complete(*completeFlag, file)
		return
	}

	if flag.NArg() == 0 {
		List()
		return
	}

	print("invalid argument\n")
	print("Try 'todo -h' for more information\n")

	defer file.Close()
}
