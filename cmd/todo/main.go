package main

import (
	"bufio"
	"flag"
	"fmt"
	todo "github.com/allrivenjs/CLI-TODO-APP"
	"io"
	"os"
	"strings"
)

const (
	todoFile = ".todos.json"
)

func main() {
	add := flag.Bool("add", false, "Add a new task")
	complete := flag.Int("complete", 0, "Complete a todo as completed")
	del := flag.Int("delete", 0, "Delete a todo from the list")
	list := flag.Bool("list", false, "List all todos")
	flag.Parse()

	todos := &todo.Todos{}
	if err := todos.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		panic(err)
	}

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		todos.Add(task)
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		err = todos.Store(todoFile)
	case *del > 0:
		err := todos.Delete(*del)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		err = todos.Store(todoFile)
	case *list:
		todos.Print()
	default:
		fmt.Println("No command specified")
		os.Exit(0)
	}
}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	text := scanner.Text()
	if len(text) == 0 {
		return "", fmt.Errorf("input is required")
	}
	return text, nil
}
