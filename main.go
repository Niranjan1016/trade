package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type User struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

const (
	host     = "localhost"
	port     = 5433
	dbuser   = "postgres"
	password = "password"
	dbname   = "finnhub"
)

func health(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Healthy server")
	w.Header().Add("content-type", "application/json")
	w.Write([]byte("HEALTHY"))
}

func (app *App) Initialize() {
	fmt.Println("Initializaling DB connection")
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, dbuser, password, dbname)
	var err error
	app.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to PostgresDB", err)
	}
	err = app.DB.Ping()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", health).Methods("GET")
	r.HandleFunc("/register", app.SignUp).Methods("POST")
	app.Router = r

}

func (app *App) SignUp(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	w.Header().Add("Content-Type", "application/json")
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error Occured while parsing request body"))
	}

	if user.UserName == "" {
		w.WriteHeader(400)
		w.Write([]byte("UserName cannot be empty"))
	}

	if user.Password == "" {
		w.WriteHeader(400)
		w.Write([]byte("Password cannot be empty"))
	}

	if user.Email == "" {
		w.WriteHeader(400)
		w.Write([]byte("Email cannot be empty"))
	}

	queryErr := app.DB.QueryRow("insert into finnhub.\"Users\"(username, password, email , created_on) VALUES($1, $2, $3, $4) returning user_id", user.UserName, user.Password, user.Email, time.Now()).Scan(&user.UserID)

	if queryErr != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error inserting record into database :" + queryErr.Error()))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("User " + user.UserName + " created successfuly"))

}

func main() {
	var app = App{}
	fmt.Println("This is test")
	app.Initialize()
	defer app.DB.Close()
	http.ListenAndServe("localhost:8080", app.Router)
}
