package main

import (
	"fmt"
	"net/http"
  "html/template"
	"time"
	"strconv"
)

type OrderDetail struct {
	ID int
	Picker *string
	Shipper *string
	Picktime time.Time
	Shiptime time.Time
}

type Page struct{
  Title string
  Message Message
	Order OrderDetail
	Permission string
	Startdate string
	Enddate string
	Graph1 []Graph
	Graph2 []Graph
	Graph3 []Graph
	Graph4 []Graph
	Graph5 []Graph
	Graph6 []Graph
	Table1 []Table
}

type Message struct {
	Success bool
	Title string
  Body string
}

func message(r *http.Request) (messagebox Message){
  if r.URL.Query().Get("messagetitle") != ""{
    messagebox.Body = r.URL.Query().Get("messagebody")
    fmt.Println("Message: ",messagebox)
  }
  return messagebox
}

func handleerror(err error) (message Message){
  if err != nil {
    message.Title = "Error"
    message.Success = false
    message.Body = err.Error()
    fmt.Println("Message: ",message.Body)
    return message
  }
  message.Success = true
  message.Body = "Success"
  return message
}

func main() {
  	fmt.Println("Starting Server...")
	var messagebox Message
	// db, messagebox = opendb()
	fmt.Println(messagebox.Body)
	http.HandleFunc("/", login)
  	http.HandleFunc("/login", login)
  	http.HandleFunc("/signup", signup)
  	http.HandleFunc("/logout", Logout)
  	http.HandleFunc("/signin", Signin)
  	http.HandleFunc("/usercreate", Usercreate)
	http.HandleFunc("/order", Order)
	http.HandleFunc("/error", Error)
	http.HandleFunc("/dashboard", Dashboard)
	http.HandleFunc("/products", Products)
	http.ListenAndServe(":8082", nil)
}

func report(w http.ResponseWriter, r *http.Request) {
  tmpl, err := template.ParseFiles("layout.html")
  fmt.Println(err)
  tmpl.Execute(w,"")
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	var page Page
	// page.Permission = auth(w,r)
    t, _ := template.ParseFiles("dashboard.html","header.html","login.js")
	fmt.Println("Loading Dashboard...")
    page.Title = "Dashboard"
	var startdate time.Time
	var enddate time.Time
	if r.FormValue("startdate") != "" && r.FormValue("enddate") != "" {
		startdate,_ = time.Parse("2006-01-02",r.FormValue("startdate"))
		enddate,_ = time.Parse("2006-01-02",r.FormValue("enddate"))
	} else {
		startdate = time.Now().AddDate(0,0,-21)
		enddate = time.Now()
	}
	fmt.Println("Start:",r.FormValue("startdate")," End:",enddate)
	// page.Message,page.Graph1 = Efficiency(startdate,enddate)
	// page.Message,page.Graph2 = Groupefficiency(startdate,enddate)
	// page.Message,page.Graph3 = ErrorLookup(startdate,enddate)
	// page.Message,page.Graph4 = Servicelevel(time.Now().AddDate(0,0,-63),time.Now())
	// page.Message,page.Table1 = ErrorList(startdate,enddate,30)
	page.Startdate = startdate.Format("2006-01-02")
	page.Enddate = enddate.Format("2006-01-02")
    fmt.Println(page)
    t.Execute(w, page)
}

func Products(w http.ResponseWriter, r *http.Request) {
	var page Page
	// page.Permission = auth(w,r)
    t, _ := template.ParseFiles("products.html","header.html","login.js")
	fmt.Println("Loading Products...")
    page.Title = "Products"
	var startdate time.Time
	var enddate time.Time
	if r.FormValue("startdate") != "" && r.FormValue("enddate") != "" {
		startdate,_ = time.Parse("2006-01-02",r.FormValue("startdate"))
		enddate,_ = time.Parse("2006-01-02",r.FormValue("enddate"))
	} else {
		startdate = time.Now().AddDate(0,0,-21)
		enddate = time.Now()
	}
	fmt.Println("Start:",r.FormValue("startdate")," End:",enddate)
	page.Startdate = startdate.Format("2006-01-02")
	page.Enddate = enddate.Format("2006-01-02")
    fmt.Println(page)
    t.Execute(w, page)
}

func Order(w http.ResponseWriter, r *http.Request) {
		var page Page
		page.Permission = auth(w,r)
    t, _ := template.ParseFiles("order.html","header.html","login.js")
		fmt.Println("Looking up order ",r.FormValue("ordernum"))
    page.Title = "Order Lookup"
    ordernum,err := strconv.Atoi(r.FormValue("ordernum"))
		if err != nil {
			page.Message.Body = err.Error()
		}
		page.Message,page.Order  = Orderlookup(ordernum)
		// page.Order.ID=67099
    fmt.Println(page)
    t.Execute(w, page)
}

func Error(w http.ResponseWriter, r *http.Request) {
		var page Page
	  fmt.Println("Comment: ",r.FormValue("comment"))
	  fmt.Println("Issue: ",r.FormValue("issue"))
	  fmt.Println("orderid: ",r.FormValue("orderid"))
		t, _ := template.ParseFiles("error.html","header.html","login.js")
		page.Permission = auth(w,r)
		page.Message = message(r)
    page.Title = "Error Entry"
		orderid,err := strconv.Atoi(r.FormValue("orderid"))
		page.Message = handleerror(err)
		page.Message = ErrorEnter(r.FormValue("comment"),r.FormValue("issue"),orderid)
    t.Execute(w, page)
}

func login(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("login.html","header.html","login.js")
    var page Page
    page.Title = "Login"
    page.Message = message(r)
    fmt.Println(page)
    t.Execute(w, page)
}

func signup(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("signup.html","header.html","login.js")
    var page Page
    page.Title = "Sign Up"
    page.Message = message(r)
    fmt.Println(page)
    t.Execute(w, page)
}
