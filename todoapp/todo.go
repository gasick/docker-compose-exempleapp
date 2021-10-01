package main

import (
	"database/sql"
	"encoding/json"

	// "fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Todo struct {
	Id          int    `json: "id", db: "id"`
	Name        string `json:"name", db:"name"`
	Description string `json:"description", db:"description"`
}

func main() {
	initDB()
	// fmt.Printf("it worked!")

	initRouter()
}

func initRouter() {
	router := mux.NewRouter()

	// $ curl http://localhost:8000/ -v
	router.HandleFunc("/", Home)

	// $ curl -H "Content-Type: application/json" http://localhost:8000/todos -d '{"name":"Wash the garbage","description":"Be especially thorough"}' -v
	router.HandleFunc("/todos/", CreateTodo).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

var Home = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Всё работает!")
}

var CreateTodo = func(w http.ResponseWriter, r *http.Request) {
	todo := &Todo{}

	err := json.NewDecoder(r.Body).Decode(todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sqlStatement := `
        INSERT INTO todos (name, description)
        VALUES ($1, $2)
        RETURNING id
    `

	id := 0
	err = db.QueryRow(sqlStatement, todo.Name, todo.Description).Scan(&id)
	if err != nil {
		panic(err)
	}

	todo.Id = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// hook up to postgres db
func initDB() {
	var err error
	db, err = sql.Open("postgres", "dbname=gotodo sslmode=disable")

	if err != nil {
		panic(err)
	}
}
