package main

import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
"log"
"fmt"
"time"
"os"
)

type Graph []struct {
  Date      time.Time
  User      string
  Station   string
  items     int
  hours     float64
}

var db *sql.DB

func opendb() (db *sql.DB, messagebox Message) {
  // Get a database handle.
  var err error
  // var user string
  fmt.Println("Connecting to DB...")
  fmt.Println("user:",os.Getenv("USER"))
  fmt.Println("pass:",os.Getenv("PASS"))
  fmt.Println("server:",os.Getenv("SERVER"))
  fmt.Println("port:",os.Getenv("PORT"))
  fmt.Println("Opening Database...")
  connectstring := os.Getenv("USER")+":"+os.Getenv("PASS")+"@tcp("+os.Getenv("SERVER")+":"+os.Getenv("PORT")+")/orders"
  fmt.Println("Connection: ",connectstring)
  db, err = sql.Open("mysql",
  connectstring)
  if err != nil {
    messagebox.Success = false
    messagebox.Body = err.Error()
    fmt.Println("Message: ",messagebox.Body)
    return nil,messagebox
  }

  //Test Connection
  pingErr := db.Ping()
  if pingErr != nil {return nil,handleerror(pingErr)}

  //Success!
    fmt.Println("Returning Open DB...")
    messagebox.Success = true
    messagebox.Body = "Success"
  return db,messagebox
}

func efficiency() (graph Graph, message Message){

  //Test Connection
  pingErr := db.Ping()
  if pingErr != nil {return nil,handleerror(pingErr)}
  return
}

//Authenticate user from DB
func userauth(user string, pass string) (permission string, message Message){
    // Get a database handle.
    var err error
    //Test Connection
    pingErr := db.Ping()
    if pingErr != nil {
        return "notfound",handleerror(pingErr)
    }
    //set Variables
    //Query
    var newquery string = "select permissions from users where username = ? and password = ?"
		// fmt.Println(newquery)
    rows, err := db.Query(newquery,user,pass)
    if err != nil {
        return "notfound",handleerror(err)
    }
    defer rows.Close()
    //Pull Data
    for rows.Next() {
    	err := rows.Scan(&permission)
    	if err != nil {
          return "notfound",handleerror(err)
    	}
    }
    err = rows.Err()
    if err != nil {
        return "notfound",handleerror(err)
    }
	if permission == "" {
    message.Title = "User not found"
    message.Body = "User/password not found in database. Please try again."
		return "notfound",message
	}
  message.Title = "Success"
  message.Body = "Successfully logged in"
  return permission, message
}

//Authenticate user from DB
func userdata(user string) (permission string){
    // Get a database handle.
    var err error
    //Test Connection
    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    //set Variables
    //Query
    var newquery string = "select permissions from users where username = ?"
		// fmt.Println(newquery)
    rows, err := db.Query(newquery,user)
    if err != nil {
    	log.Fatal(err)
    }
    defer rows.Close()
    //Pull Data
    for rows.Next() {
    	err := rows.Scan(&permission)
    	if err != nil {
    		log.Fatal(err)
    	}
    }
    err = rows.Err()
    if err != nil {
    	log.Fatal(err)
    }
	if permission == "" {
		return "notfound"
	}

  return permission
}
