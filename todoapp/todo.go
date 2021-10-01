package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

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
	migrateDB()

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
	hostname := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DB")
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	initDBParams := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		hostname, port, username, password, database)
	var err error
	db, err = sql.Open("postgres", initDBParams)

	if err != nil {
		panic(err)
	}
}

func migrateDB() {
	sql := `
		create table todos (
			id serial primary key,
			name text,
			description text
		);
	`

	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
}
