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
  X        *string
  Y        *float64
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
  if pingErr != nil {
    return nil,handleerror(pingErr)
  }

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
    db, message = opendb()
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

  var newquery string = "SELECT d.user,sum(d.items)/sum(e.hours) FROM (SELECT a.date,a.user,c.usercode,sum(b.items_total) items FROM (SELECT ordernum, station, user, DATE(scans.time) as date from scans where station='pick' group by ordernum, station, user, DATE(scans.time)) a INNER JOIN (SELECT id, items_total from orders) b on a.ordernum = b.id LEFT JOIN (SELECT usercode,username from users) c on a.user = c.username GROUP BY a.date,a.user,c.usercode) d LEFT JOIN (SELECT DATE(clock_in) clockin,payroll_id, sum(paid_hours) hours from shifts where role='Shipping' group by DATE(clock_in),payroll_id) e on d.date = e.clockin and d.usercode = e.payroll_id WHERE d.items IS NOT NULL and e.hours IS NOT NULL GROUP BY d.user ORDER BY 1,2;"

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
    err := rows.Scan(&r.X,&r.Y)
    if err != nil {
      return handleerror(err),graph
    }
    graph = append(graph,r)
  }
  return message,graph
}

func Groupefficiency() (message Message, graph []Graph){

  //Test Connection
  pingErr := db.Ping()
  if pingErr != nil {return handleerror(pingErr),graph}

  var newquery string = "SELECT shipments.date, items/hours efficiency FROM (select CAST(c.time as date) date, sum(a.items_total) items from orders a  LEFT JOIN (select * FROM scans where station='ship') c ON a.id = c.ordernum  WHERE a.statusid not in (0) and c.time is not null  GROUP BY CAST(c.time as date)  ) shipments  LEFT JOIN (select cast(clock_in as date) date,sum(paid_hours) hours FROM shifts WHERE role = 'Shipping' group by cast(clock_in as date)) d on d.date = shipments.date  WHERE items is not null and hours is not null order by 1;"

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
    err := rows.Scan(&r.X,&r.Y)
    if err != nil {
      return handleerror(err),graph
    }
    graph = append(graph,r)
  }
  return message,graph
}

func Updatepass (user string, pass string, secret string) (message Message, success bool){
  pingErr := db.Ping()
  if pingErr != nil {
      return handleerror(pingErr),false
  }

  //Check for secret
  if secret != os.Getenv("SECRET") {
    message.Title = "Secret Auth Failed"
    message.Body = "Secret Auth Failed"
    return message,false
  }

  hashpass := hashAndSalt([]byte(pass))
  fmt.Println("Creating password hash of length ",len(hashpass),": ", hashpass)
  var newquery string = "update users set password = ? where username = ? and password = ''"
  rows, err := db.Query(newquery,hashpass,user)
  if err != nil {
      return handleerror(err),false
  }
  defer rows.Close()
  message.Title = "Success"
  message.Body = "Success"
  return message,true
}

//Authenticate user from DB
func userauth(user string, pass string) (permission string, message Message){
    // Get a database handle.
    var err error
    var dbpass string
    //Test Connection
    pingErr := db.Ping()
    if pingErr != nil {
        return "notfound",handleerror(pingErr)
    }
    //set Variables
    //Query
    var newquery string = "select password, permissions from users where username = ?"
		// fmt.Println(newquery)
    rows, err := db.Query(newquery,user)
    if err != nil {
        return "notfound",handleerror(err)
    }
    defer rows.Close()
    //Pull Data
    for rows.Next() {
    	err := rows.Scan(&dbpass,&permission)
    	if err != nil {
          return "notfound",handleerror(err)
    	}
    }
    err = rows.Err()
    if err != nil {
        return "notfound",handleerror(err)
    }

    fmt.Println("Checking Permissions: ", permission)
  //If user has not set a password
  if dbpass == "" {
    message.Title = "Set Password"
    message.Body = "Password not set, please create password"
		return "newuser",message
  }

  //If Permissions do not exist for user
	if permission == "" {
    message.Title = "Permission not found"
    message.Body = "Permissions not set for user. Please contact your system administrator."
		return "notfound",message
	}

  if comparePasswords(dbpass,[]byte(pass)) {
    message.Title = "Success"
    message.Body = "Successfully logged in"
    // permission = "notfound"
    return permission, message
  }
  message.Title = "Login Failed"
  message.Body = "Login Failed"
  permission = "notfound"
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
