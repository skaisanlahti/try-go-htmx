package todos

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type todo struct {
	Id   string
	Task string
	Done bool
}

type page struct {
	Todos []todo
	Task  string
	Error string
	Key   int64 // use as cache key or name attribute to bypass caching in browser
}

type service struct {
	page  *template.Template
	todos map[string]todo
	mutex sync.RWMutex
}

//go:embed todo-page.html
var templateFS embed.FS

func NewService() *service {
	return &service{
		page:  template.Must(template.ParseFS(templateFS, "todo-page.html")),
		todos: make(map[string]todo),
		mutex: sync.RWMutex{},
	}
}

const (
	HomePath       = "/todos"
	GetTodosPath   = "/todos/list"
	AddTodoPath    = "/todos/add"
	ToggleTodoPath = "/todos/toggle/"
	RemoveTodoPath = "/todos/remove/"
)

func (s *service) newPageData() page {
	page := page{Key: time.Now().UnixMilli()}
	for _, t := range s.todos {
		page.Todos = append(page.Todos, t)
	}

	sort.Slice(page.Todos, func(i, j int) bool {
		return page.Todos[i].Task < page.Todos[j].Task
	})

	return page
}

func (s *service) Home(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	page := s.newPageData()
	log.Printf("Rendered page")
	s.page.Execute(w, page)
}

func (s *service) GetTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()
	page := s.newPageData()
	log.Printf("Got %d todos", len(page.Todos))
	s.page.ExecuteTemplate(w, "list", page)
}

func (s *service) AddTodo(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	task := r.FormValue("task")
	if task == "" {
		page := s.newPageData()
		page.Error = "Task can't be empty."
		s.page.ExecuteTemplate(w, "page", page)
		return
	}

	todo := todo{Id: uuid.NewString(), Task: task, Done: false}
	s.todos[todo.Id] = todo
	log.Printf("Added todo %s", todo.Id)
	page := s.newPageData()
	s.page.ExecuteTemplate(w, "page", page)
}

func (s *service) ToggleTodo(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pathParts := getPathParts(r.URL.Path)
	id := pathParts[2]
	s.mutex.Lock()
	defer s.mutex.Unlock()
	t, found := s.todos[id]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t.Done = !t.Done
	s.todos[id] = t
	log.Printf("Toggled todo %s: %t", t.Id, t.Done)
	page := s.newPageData()
	s.page.ExecuteTemplate(w, "list", page)
}

func (s *service) RemoveTodo(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	pathParts := getPathParts(r.URL.Path)
	id := pathParts[2]
	s.mutex.Lock()
	defer s.mutex.Unlock()
	t, found := s.todos[id]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	delete(s.todos, t.Id)
	log.Printf("Deleted todo %s", t.Id)
	page := s.newPageData()
	s.page.ExecuteTemplate(w, "list", page)
}

func getPathParts(path string) []string {
	return strings.Split(strings.Trim(path, "/"), "/")
}
