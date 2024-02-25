package CLI_TODO_APP

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

// Add a new task to the list
func (t *Todos) Add(task string) {
	*t = append(*t, item{Task: task, Done: false, CreatedAt: time.Now(), CompletedAt: time.Time{}})
}

// ValidateIndex checks if the index is within the range of the list
func (t *Todos) ValidateIndex(index int) error {
	if index <= 0 || index > len(*t) {
		return errors.New("index out of range")
	}
	return nil
}

// Complete a task from the list
func (t *Todos) Complete(index int) error {
	ls := *t
	if err := ls.ValidateIndex(index); err != nil {
		return err
	}
	ls[index-1].Done = true
	ls[index-1].CompletedAt = time.Now()
	return nil
}

// Delete a task from the list
func (t *Todos) Delete(index int) error {
	ls := *t
	if err := t.ValidateIndex(index); err != nil {
		return err
	}
	*t = append(ls[:index-1], ls[index:]...)
	return nil
}

func (t *Todos) Load(filename string) error {
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
	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (t *Todos) Print() {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "?Done"},
			{Align: simpletable.AlignCenter, Text: "Created At"},
			{Align: simpletable.AlignCenter, Text: "Completed At"},
		},
	}
	var cells [][]*simpletable.Cell
	for i, item := range *t {
		i++
		task := blue(fmt.Sprintf("%s", item.Task))
		done := gray(fmt.Sprintf("%s  ", "❌"))
		if item.Done {
			task = green(fmt.Sprintf("%s", item.Task))
			done = green(fmt.Sprintf("%s   ", "✅"))
		}

		cells = append(cells, *&[]*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", i)},
			{Text: task},
			{Align: simpletable.AlignCenter, Text: done},
			{Align: simpletable.AlignCenter, Text: item.CreatedAt.Format(time.RFC850)},
			{Align: simpletable.AlignCenter, Text: item.CompletedAt.Format(time.RFC850)},
		})
	}
	table.Body = &simpletable.Body{Cells: cells}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 3, Text: red(fmt.Sprintf("Pending: %d", t.CountPending()))},
			{Align: simpletable.AlignRight, Span: 2, Text: green(fmt.Sprintf("Total: %d", len(*t)))},
		},
	}
	table.SetStyle(simpletable.StyleUnicode)

	table.Print()
}

func (t *Todos) CountPending() int {
	total := 0
	for _, item := range *t {
		if !item.Done {
			total++
		}
	}
	return total
}
