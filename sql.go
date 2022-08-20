package main

import (
"database/sql"
_ "github.com/go-sql-driver/mysql"
"log"
"fmt"
"time"
"os"
)

type Graph struct {
  User        *string
  Efficiency  *float64
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
  connectstring := os.Getenv("USER")+":"+os.Getenv("PASS")+"@tcp("+os.Getenv("SERVER")+":"+os.Getenv("PORT")+")/orders?parseTime=true"
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

func Orderlookup(ordernum int) (message Message,orderdetail OrderDetail) {
  // Get a database handle.
  var err error

  //Test Connection
  pingErr := db.Ping()
  if pingErr != nil {
    return handleerror(pingErr),orderdetail
  }

  //Query
  var newquery string = "select a.id,b.user,b.time,c.user,c.time from orders a LEFT JOIN (select * FROM scans where station='pick') b ON a.id = b.ordernum LEFT JOIN (select * FROM scans where station='ship') c ON a.id = c.ordernum  WHERE a.statusid not in (0) and a.id = ? order by 1,5;"

  //Run Query
  fmt.Println("Looking up order: ",ordernum)
  location, err := time.LoadLocation("America/Chicago")
  rows, err := db.Query(newquery,ordernum)
  if err != nil {
    return handleerror(err),orderdetail
  }
  defer rows.Close()

  //Pull Data
  for rows.Next() {
    err := rows.Scan(&orderdetail.ID,&orderdetail.Picker,&orderdetail.Picktime,&orderdetail.Shipper,&orderdetail.Shiptime)
    if err != nil {
      return handleerror(err),orderdetail
    }
  }
  if orderdetail.ID == 0 {message.Body = "Order not found"}
  orderdetail.Picktime = orderdetail.Picktime.In(location)
  orderdetail.Shiptime = orderdetail.Shiptime.In(location)
  return message, orderdetail
}

func Efficiency() (message Message, graph []Graph){

  //Test Connection
  pingErr := db.Ping()
  if pingErr != nil {return handleerror(pingErr),graph}

  var newquery string = "SELECT d.user,sum(d.items)/sum(e.hours) FROM (SELECT a.date,a.user,c.usercode,sum(b.items_total) items FROM (SELECT ordernum, station, user, DATE(scans.time) as date from scans where station='pick' group by ordernum, station, user, DATE(scans.time)) a INNER JOIN (SELECT id, items_total from orders) b on a.ordernum = b.id LEFT JOIN (SELECT usercode,username from users) c on a.user = c.username GROUP BY a.date,a.user,c.usercode) d LEFT JOIN (SELECT DATE(clock_in) clockin,payroll_id, sum(paid_hours) hours from shifts where role='Shipping' group by DATE(clock_in),payroll_id) e on d.date = e.clockin and d.usercode = e.payroll_id GROUP BY d.user ORDER BY 1,2;"

  //Run Query
  fmt.Println("Running Report")
  // location, err := time.LoadLocation("America/Chicago")
  rows, err := db.Query(newquery)
  if err != nil {
    return handleerror(err),graph
  }
  defer rows.Close()

  //Pull Data
  for rows.Next() {
    var r Graph
    err := rows.Scan(&r.User,&r.Efficiency)
    if err != nil {
      return handleerror(err),graph
    }
    graph = append(graph,r)
  }
  return message,graph
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
