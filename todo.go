package main

import (
	// "bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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
	if index < 0 || index >= len(todoList) {
		print("invalid index\n")
		return
	}
	completeList = append(completeList, todoList[index])
	todoList = append(todoList[:index], todoList[index+1:]...)
	update(file)
	ListFull()
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
		fmt.Sprint("%v. %v\n", i+1, completeList[i])
	}
}

func Clear() {

}

func Help() {
	fmt.Println("Usage: todo [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func update(file *os.File) {
	os.Truncate(file.Name(), 0)
	lists.Todo = todoList
	lists.Completed = completeList

	byteValue, _ := json.Marshal(lists)

	_, err := file.Write(byteValue)
	check(err)
}

func main() {
	// file, error := os.openFile("list.json", os.O_RDWR|os.O_CREATE, 0644)

	file, err := os.OpenFile("list.json", os.O_RDWR|os.O_CREATE, 0644)
	check(err)

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &lists)

	fmt.Printf("%+v\n", lists)
	todoList = lists.Todo
	completeList = lists.Completed

	flag.Parse()

	if flag.NArg() > 1 {
		print("Too many arguments\n")
		print("Try 'todo -h' for more information\n")
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
