package todo_test

import (
	"path/filepath"
	"testing"

	"github.com/hisamcode/todo"
)

// TestAdd tests the Add method of the list type
func TestAdd(t *testing.T) {
	t.Parallel()
	l := todo.List{}
	taskName := "New Task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l[0].Task)
	}
}

// TestComplete tests the Complete method of the list type
func TestComplete(t *testing.T) {
	t.Parallel()
	l := todo.List{}
	taskName := "New Task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead", taskName, l[0].Task)
	}

	if l[0].Done {
		t.Errorf("New task should not be completed")
	}

	l.Complete(1)

	if !l[0].Done {
		t.Errorf("New task should be completed")
	}
}

// TestDelete tests Delete method of the List
func TestDelete(t *testing.T) {
	t.Parallel()
	l := todo.List{}

	tasks := []string{
		"New Task 1",
		"New Task 2",
		"New Task 3",
	}

	for _, v := range tasks {
		l.Add(v)
	}

	if l[0].Task != tasks[0] {
		t.Errorf("Expected %q, got %q instead.", tasks[0], l[0].Task)
	}

	l.Delete(2)

	if len(l) != 2 {
		t.Errorf("Expected list length %d, got %d instead.", 2, len(l))
	}

	if l[1].Task != tasks[2] {
		t.Errorf("Expected %q, got %q instead", tasks[2], l[1].Task)
	}
}

// TestSaveGet tests the Save and Get method of the List
func TestSaveGet(t *testing.T) {
	t.Parallel()
	l1 := todo.List{}
	l2 := todo.List{}

	taskName := "New task"
	l1.Add(taskName)

	if l1[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead.", taskName, l1[0].Task)
	}

	td := t.TempDir()
	tf := filepath.Join(td, "haha")

	if err := l1.Save(tf); err != nil {
		t.Fatalf("Error saving list to file:%s", err)
	}

	if err := l2.Get(tf); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("Task %q should match %q task", l1[0].Task, l2[0].Task)
	}
}
