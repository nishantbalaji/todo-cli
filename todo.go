package main

import (
	// "bufio"
	"encoding/json"
	"fmt"
	"os"
	"flag"
	// "strings"
)



type Lists struct {
	Todo      []string `json:"todo"`
	Completed []string `json:"completed"`
}

file, err := os.openFile("list.json", os.O_RDWR|os.O_CREATE, 0644)
check(err)

var (
	addFlag      = flag.String("a", "", "Add a todo")
	completeFlag = flag.Int("c", 0, "Complete a todo")
	helpFlag     = flag.Bool("h", false, "Help")
	verboseFlag  = flag.Bool("v", false, "Verbose")
)


func check(err error){
	if err != nil {
		panic(err)
	}
}


func Add(todo string) {
	todoList = append(todoList, todo)
	print("added '", todo, "'\n")
	List()
}

func Remove(index int) {
	index--
	if index < 0 || index >= len(todoList) {
		print("invalid index\n")
		return
	}
	todoList = append(todoList[:index], todoList[index+1:]...)
}

func Complete(index int) {
	index--
	if index < 0 || index >= len(todoList) {
		print("invalid index\n")
		return
	}
	completeList = append(completeList, todoList[index])
	todoList = append(todoList[:index], todoList[index+1:]...)
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

func Help() {
	fmt.Println("Usage: todo [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	// file, error := os.openFile("list.json", os.O_RDWR|os.O_CREATE, 0644)

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
		Add(*addFlag)
		return
	}

	if *completeFlag > 0 {
		Complete(*completeFlag)
		return
	}

	if flag.NArg() == 0 {
		List()
		return
	}

	print("invalid argument\n")
	print("Try 'todo -h' for more information\n")
}
