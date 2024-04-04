package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

var port = 80

var PageTitle = "Todo List"

var id uint = 1

type Todo struct {
	Id    uint
	Done  bool
	Title string
}

func createTodo(title string, done bool) Todo {
	id++
	return Todo{
		Id:    id,
		Title: title,
		Done:  done,
	}
}

type Data struct {
	PageTitle string
	Todos     []Todo
}

func newData() Data {
	return Data{
		PageTitle: PageTitle,
		Todos: []Todo{
			createTodo("Cook dinner", false),
			createTodo("Clean up", false),
			createTodo("Chill", false),
			createTodo("Exercise", true),
		},
	}
}

func main() {
	router := mux.NewRouter()

	data := newData()

	tmpl := template.Must(template.ParseFiles("index.html"))

	// Delete todo
	router.HandleFunc("/todo/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		idInt, errInt := strconv.Atoi(idStr)
		id := uint(idInt)

		if errInt != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]string)
			resp["message"] = "Id not valid"
			jsonResp, err := json.Marshal(resp)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}

			w.Write(jsonResp)
			return
		}

		var nextTodos []Todo
		for _, todo := range data.Todos {
			if todo.Id != id {
				nextTodos = append(nextTodos, todo)
			}
		}

		data.Todos = nextTodos

		http.Redirect(w, r, "/", 303)
	}).Methods("POST")

	// Update todo
	router.HandleFunc("/todo/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		idInt, errInt := strconv.Atoi(idStr)
		id := uint(idInt)

		if errInt != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]string)
			resp["message"] = "Id not valid"
			jsonResp, err := json.Marshal(resp)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}

		isDoneStr := r.PostFormValue("done")
		isDone, errBool := strconv.ParseBool(isDoneStr)

		if errBool != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]string)
			resp["message"] = "Id not valid"
			jsonResp, err := json.Marshal(resp)

			if err != nil {
				log.Fatalf("Error happened in JSON marshal. Err: %s", err)
			}
			w.Write(jsonResp)
			return
		}

		var nextTodos []Todo
		for _, todo := range data.Todos {
			if todo.Id == id {
				todo.Done = isDone
			}
			nextTodos = append(nextTodos, todo)
		}

		data.Todos = nextTodos

		http.Redirect(w, r, "/", 303)
	}).Methods("POST")

	// Create todo
	router.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		newTodo := r.PostFormValue("new-todo")
		fmt.Println(newTodo)
		fmt.Println(data)

		data.Todos = append(data.Todos, createTodo(newTodo, false))

		http.Redirect(w, r, "/", 303)
	}).Methods("POST")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data)
	})

	portString := fmt.Sprintf(":%v", port)

	// Static
	fs := http.FileServer(http.Dir("assets/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	fmt.Printf("Listening on port %v\n", port)
	http.ListenAndServe(portString, router)
}
