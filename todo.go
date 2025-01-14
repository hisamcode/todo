package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type List []item

func (l *List) String() string {
	var out string
	out += "No\tTask\tComplete\n"
	for k, t := range *l {
		out += fmt.Sprintf("%d\t%s\t%t\n", k+1, t.Task, t.Done)
	}
	return out
}

func (l *List) StringWithOptions(verbose, hideComplete bool) string {
	var out string
	out += "No\tTask\tComplete"
	if verbose {
		out += "\tCreated At\tCompleted At"
	}
	out += "\n"
	var completedAt string
	var num int
	for _, t := range *l {
		num++
		if hideComplete {
			if t.Done {
				continue
			}
		}
		out += fmt.Sprintf("%d\t%s\t%t", num, t.Task, t.Done)
		if verbose {
			completedAt = "-"
			if t.CompletedAt != nil {
				completedAt = t.CompletedAt.Format(time.DateTime)
			}
			out += fmt.Sprintf("\t%s\t%s", t.CreatedAt.Format(time.DateTime), completedAt)
		}
		out += "\n"
	}
	return out
}

// Add method creates a new todo item and appends it to the list
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}

	*l = append(*l, t)
}

// Complete method marks a ToDo item as completed
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}
	// adjusting index for 0 based index
	now := time.Now()
	ls[i-1].Done = true
	ls[i-1].CompletedAt = &now
	return nil
}

// Delete method deletes a ToDo item from the list
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}
	*l = append(ls[:i-1], ls[i:]...)
	return nil
}

const (
	FilePermissionReadWrite = 0644
)

// Save method encodes the List as JSON and saves it using the provided file name
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, js, FilePermissionReadWrite)
}

// Get method opens the provided file name, decodes
// the JSON data and parses it into List
func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)
}
