package main

import (
	"fmt"
	"net/http"
  "html/template"
)

type Page struct{
  Title string
  Message Message
}

type Message struct {
	Success bool
	// User string
	Title string
  Body string
}

func message(r *http.Request) (messagebox Message){
  if r.URL.Query().Get("messagetitle") != ""{
    messagebox := Message{false,r.URL.Query().Get("messagetitle"),r.URL.Query().Get("messagebody")}
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
  http.HandleFunc("/signin", Signin)
	http.ListenAndServe(":8081", nil)
}

func report(w http.ResponseWriter, r *http.Request) {
  tmpl, err := template.ParseFiles("layout.html")
  fmt.Println(err)
  tmpl.Execute(w,"")
}

func login(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("login.html","header.html","login.js")
    var page Page
    page.Title = "testTitle"
    page.Message = message(r)
    fmt.Println(page)
    t.Execute(w, page)
}
