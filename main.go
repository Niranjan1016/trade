package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"
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

func health(w http.ResponseWriter, r *http.Request) {
        tag :=  os.Getenv("DEPLOYMENT_TAG")
	fmt.Println("Healthy server")
	w.Header().Add("content-type", "application/json")
	w.Write([]byte("HEALTHY " + tag))
}

func (app *App) Initialize() {
	fmt.Println("Initializaling DB connection")
        host := os.Getenv("POSTGRES_SERVICE_HOST")
	port := os.Getenv("POSTGRES_SERVICE_PORT")
	dbuser := os.Getenv("POSTGRES_DB_USERNAME")
	password := os.Getenv("POSTGRES_DB_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
        tag :=  os.Getenv("DEPLOYMENT_TAG")
        fmt.Printf("Postgres Service Host : %v\n", host)
	fmt.Printf("Postgres Service Port : %v\n", port)
	fmt.Printf("Postgres User : %v\n", dbuser)
	fmt.Printf("Postgres DB : %v\n", dbname)
        fmt.Printf("Deployment Name: %v\n", tag)
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, dbuser, password, dbname)
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
	r.HandleFunc("/trade/health", health).Methods("GET")
	r.HandleFunc("/trade/register", app.SignUp).Methods("POST")
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
	http.ListenAndServe(":8080", app.Router)
}
