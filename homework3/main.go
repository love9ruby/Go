package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Todo struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type TodoList struct {
	Todos []Todo `json:"todos"`
}

func (list *TodoList) toJSON() string {
	jsonData, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		fmt.Println("JSON marshaling error:", err)
		return ""
	}
	return string(jsonData)
}

func (list *TodoList) toCSV() string {
	csv := "text,done\n"
	for _, todo := range list.Todos {
		csv += fmt.Sprintf("%s,%t\n", todo.Text, todo.Done)
	}
	return csv
}

func (list *TodoList) toHTML() string {
	html := "<html><head><title>Todo List</title></head><body><h1>Todo List</h1><ul>"
	for _, todo := range list.Todos {
		html += fmt.Sprintf("<li>%s</li>", todo.Text)
	}
	html += "</ul></body></html>"
	return html
}

func (list *TodoList) toMarkdown() string {
	markdown := "# Todo List\n\n"
	for _, todo := range list.Todos {
		markdown += fmt.Sprintf("- [%s]\n", todo.Text)
	}
	return markdown
}

func (list *TodoList) toText() string {
	text := "Todo List\n\n"
	for _, todo := range list.Todos {
		text += fmt.Sprintf("- %s\n", todo.Text)
	}
	return text
}

func writeToFile(filename, content string) error {
	// string to []byte
	data := []byte(content)
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("File writing error:", err)
		return err
	}
	return nil
}

type AsyncWriteFileResult interface {
	Write(filename, content string)
	Wait()
}

type wgImpl struct {
	wg *sync.WaitGroup
}

type atomicImpl struct {
	an *atomic.Int32
}

type chImpl struct {
	ch chan bool
}

func NewAsyncWriteFileResult(im int) AsyncWriteFileResult {
	switch im {
	case 1:
		return &wgImpl{wg: &sync.WaitGroup{}}
	case 2:
		return &atomicImpl{an: &atomic.Int32{}}
	case 3:
		return &chImpl{ch: make(chan bool)}
	}
	return nil
}

func (w *wgImpl) Write(filename, content string) {
	w.wg.Add(1)
	defer w.wg.Done()
	err := writeToFile(filename, content)
	if err != nil {
		return
	}
}

func (w *wgImpl) Wait() {
	w.wg.Wait()
}

func (a *atomicImpl) Write(filename, content string) {
	err := writeToFile(filename, content)
	if err != nil {
		return
	}
	a.an.Add(1)
}

func (a *atomicImpl) Wait() {
	for a.an.Load() < 5 {
		time.Sleep(time.Millisecond)
	}
}

func (c *chImpl) Write(filename, content string) {
	err := writeToFile(filename, content)
	if err != nil {
		c.ch <- false
		return
	}
	c.ch <- true
}

func (c *chImpl) Wait() {
	for i := 0; i < 5; i++ {
		select {
		case result := <-c.ch:
			if result {
				fmt.Println("File written successfully")
			} else {
				fmt.Println("File writing error")
			}
		}
	}
}

func exportTodos(todos TodoList) {
	awft := NewAsyncWriteFileResult(1)
	nameArray := []string{"todo.json", "todo.csv", "todo.html", "todo.md", "todo.txt"}
	contentArray := []string{todos.toJSON(), todos.toCSV(), todos.toHTML(), todos.toMarkdown(), todos.toText()}
	for i, name := range nameArray {
		go awft.Write(name, contentArray[i])
	}
	awft.Wait()
}

func main() {
	//data
	todoList := &TodoList{}

	a := app.New()
	w := a.NewWindow("Hello")
	w.Resize(fyne.NewSize(400, 600))

	input := widget.NewEntry()
	input.SetPlaceHolder("Add a todo")
	title := widget.NewLabel("My Todo List")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	list := widget.NewList(
		func() int {
			return len(todoList.Todos)
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("template")
			label.TextStyle = fyne.TextStyle{Monospace: true, Italic: true}
			label.Alignment = fyne.TextAlignCenter
			return label
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if todoList.Todos[id].Done {
				item.(*widget.Label).SetText("V " + todoList.Todos[id].Text)
			} else {
				item.(*widget.Label).SetText("X " + todoList.Todos[id].Text)
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		todoList.Todos[id].Done = !todoList.Todos[id].Done
		refreshList(list)
	}
	// 在標題下面添加兩個空白行
	space1 := widget.NewLabel("")
	space2 := widget.NewLabel("")
	hello := widget.NewLabel("Todo List!")
	// now clock
	now := time.Now()
	nowLabel := widget.NewLabel(now.Format("2006-01-02 15:04:05"))
	go func() {
		for {
			time.Sleep(time.Second)
			now = time.Now()
			nowLabel.SetText(now.Format("2006-01-02 15:04:05"))
		}
	}()

	addButton := widget.NewButton("Add", func() {
		if input.Text == "" {
			return
		}
		todoList.Todos = append(todoList.Todos, Todo{Text: input.Text})
		input.SetText("")
		refreshList(list)
	})
	deleteAllButton := widget.NewButton("Delete All", func() {
		todoList.Todos = nil
		refreshList(list)
	})
	exportAllButton := widget.NewButton("Export All", func() {
		exportTodos(*todoList)
	})
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))
	content := container.NewBorder(
		container.NewVBox(title, space1, space2, input, addButton, deleteAllButton, exportAllButton),
		nowLabel, nil, nil, list,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
func refreshList(list *widget.List) {
	list.Refresh()
}











