package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/hisamcode/todo"
)

var todoFileName = ".todo.json"

func main() {
	add := flag.Bool("add", false, "Add task to the ToDo list, without using parameter it show prompt stdin, in stdin enter on an empty line to finish.")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("del", 0, "Item to be deleted")
	verbose := flag.Bool("v", false, "verbose show date")
	hc := flag.Bool("hc", false, "hide complete task")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool.\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2025\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		printTW(l.StringWithOptions(*verbose, *hc))
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		lt := strings.Split(t, "\n")
		for _, t := range lt {
			l.Add(t)
		}

		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *del > 0:
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	fmt.Println("Enter multiline input. Press Enter on an empty line to finish:")
	s := bufio.NewScanner(r)
	var multiline string
	for s.Scan() {
		line := s.Text()
		if line == "" {
			break
		}
		multiline += line + "\n"
	}
	if err := s.Err(); err != nil {
		return "", err
	}

	multiline = strings.TrimSpace(multiline)
	if len(multiline) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}

	return multiline, nil
}

func printTW(s string) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprint(tw, s)
	tw.Flush()
}

func Main() int {
	main()
	return 0
}
