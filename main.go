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
	Graph []Graph
}

type Message struct {
	Success bool
	// User string
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
	db, messagebox = opendb()
	fmt.Println(messagebox.Body)
	http.HandleFunc("/", report)
  http.HandleFunc("/login", login)
  http.HandleFunc("/logout", Logout)
  http.HandleFunc("/signin", Signin)
	http.HandleFunc("/order", Order)
	http.HandleFunc("/dashboard", Dashboard)
	http.ListenAndServe(":8081", nil)
}

func report(w http.ResponseWriter, r *http.Request) {
  tmpl, err := template.ParseFiles("layout.html")
  fmt.Println(err)
  tmpl.Execute(w,"")
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	fmt.Println(auth(w,r))
    t, _ := template.ParseFiles("dashboard.html","header.html","login.js")
		fmt.Println("Loading Dashboard...")
    var page Page
    page.Title = "Dashboard"
		// if ordernum == "" {
		// 	page.Message.Body = ""
		// }
		page.Message,page.Graph  = Efficiency()
		// page.Order.ID=67099
    fmt.Println(page)
    t.Execute(w, page)
}

func Order(w http.ResponseWriter, r *http.Request) {
	fmt.Println(auth(w,r))
    t, _ := template.ParseFiles("order.html","header.html","login.js")
		fmt.Println("Looking up order ",r.FormValue("ordernum"))
    var page Page
    page.Title = "Order Lookup"
    ordernum,err := strconv.Atoi(r.FormValue("ordernum"))
		if err != nil {
			page.Message.Body = err.Error()
		}
		// if ordernum == "" {
		// 	page.Message.Body = ""
		// }
		page.Message,page.Order  = Orderlookup(ordernum)
		// page.Order.ID=67099
    fmt.Println(page)
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
