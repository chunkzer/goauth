// This is the name of our package
// Everything with this package name can see everything
// else inside the same package, regardless of the file they are in
package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	// "fmt" has methods for formatted I/O operations (like printing to the console)
	"fmt"
	// The "net/http" library has methods to implement HTTP clients and servers
	"net/http"
  "github.com/gorilla/mux"
	"database/sql"
	"encoding/json"
	"log"
	_ "github.com/lib/pq"
	"io/ioutil"

	"models"
)

var db *sql.DB

// The new router function creates the router and
// returns it to us. We can now use this function
// to instantiate and test the router outside of the main function
func newRouter() *mux.Router {
	// "Signin" and "Signup" are handler that we will implement
	// initialize our database connection
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("GET")
	r.HandleFunc("/signup", signup).Methods("POST")
	return r
}

func main() {
	// The router is now formed by calling the `newRouter` constructor function
	// that we defined above. The rest of the code stays the same
	initDB()
	r := newRouter()
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", r))
}

func initDB(){
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	db, err = sql.Open("postgres", "dbname=goauthdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database running successfully")
}

func signup(w http.ResponseWriter	, r *http.Request) {
	fmt.Fprintf(w, "Hello Signup World!")
}

func login(w http.ResponseWriter	, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	user := &models.User{}

	err = json.Unmarshal(body, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.Username, user.Password)

	rows, err := db.Query("SELECT * from USERS WHERE username = $1", user.Username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			username   string
			password string
		)
		if err := rows.Scan(&username, &password); err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", http.StatusOK)
	}
	http.Redirect(w, r, "/", http.StatusInternalServerError)

}
