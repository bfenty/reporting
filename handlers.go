package main

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/google/uuid"
	// "io/ioutil"
	//"html/template"
)

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var sessions = map[string]session{}

// each session contains the username of the user and the time at which it expires
type session struct {
	username string
	expiry   time.Time
}

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

// Create a struct that models the structure of a user in the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func Usercreate(w http.ResponseWriter, r *http.Request){
	var creds Credentials
	var message Message
	var success bool
				// fmt.Println("method:", r.Method) //get request method
				r.ParseForm()
				// logic part of log in
				creds.Username = r.FormValue("username")
				creds.Password = r.FormValue("password")
				if creds.Password != r.FormValue("password2") {
					message.Title = "Non-matching passwords"
					message.Body = "Passwords do not match"
					http.Redirect(w, r, "/signup?messagetitle="+message.Title+"&messagebody="+message.Body, http.StatusSeeOther)
					return
				}
				fmt.Println("Creating user ",creds.Username,"...")
				message, success = Updatepass(creds.Username,creds.Password,r.FormValue("secret"))
				if success {
					http.Redirect(w, r, "/dashboard?messagetitle="+message.Title+"&messagebody="+message.Body, http.StatusSeeOther)
					return
				}
				http.Redirect(w, r, "/signup?messagetitle="+message.Title+"&messagebody="+message.Body, http.StatusSeeOther)
				return
}

func Signin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging in...")
	var creds Credentials
				// fmt.Println("method:", r.Method) //get request method
				r.ParseForm()
				// logic part of log in
					creds.Username = r.FormValue("username")
					creds.Password = r.FormValue("password")
	permission,message := userauth(creds.Username,creds.Password)
	fmt.Println(message.Body)

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if permission == "notfound" {
		fmt.Println(message.Body)
		http.Redirect(w, r, "/login?messagetitle="+message.Title+"&messagebody="+message.Body, http.StatusSeeOther)
		return
	}

	if permission == "newuser" {
		fmt.Println(message.Body)
		http.Redirect(w, r, "/signup?messagetitle="+message.Title+"&messagebody="+message.Body, http.StatusSeeOther)
		return
	}

	// Create a new random session token
	// we use the "github.com/google/uuid" library to generate UUIDs
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(300 * time.Second)
	// fmt.Println("Authorized")

	// Set the token in the session map, along with the session information
	sessions[sessionToken] = session{
		username: creds.Username,
		expiry:   expiresAt,
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	// fmt.Println(sessions)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func auth(w http.ResponseWriter, r *http.Request) (permission string){
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			//Redirect Login
			http.Redirect(w, r, "/login?messagetitle=User Unauthorized&messagebody=Please login to use this site", http.StatusSeeOther)
			// If the cookie is not set, return an unauthorized status
			return "Unauthorized"
		}
	}
	sessionToken := c.Value

	// We then get the name of the user from our session map, where we set the session token
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		// w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Unauthorized")
		http.Redirect(w, r, "/login?messagetitle=User Unauthorized&messagebody=Please login to use this site", http.StatusSeeOther)
		return "Unauthorized"
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		// w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Unauthorized")
		http.Redirect(w, r, "/login?messagetitle=User Unauthorized&messagebody=Please login to use this site", http.StatusSeeOther)
		return "Unauthorized"
	}
	// Finally, return the welcome message to the user
	// w.Write([]byte(fmt.Sprintf("Welcome %s!", userSession.username)))
	fmt.Println("Authorized")
	// If the previous session is valid, create a new session token for the current user
	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(300 * time.Second)

	// Set the token in the session map, along with the user whom it represents
	sessions[newSessionToken] = session{
		username: userSession.username,
		expiry:   expiresAt,
	}

	// Delete the older session token
	delete(sessions, sessionToken)

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(300 * time.Second),
	})
	return userdata(userSession.username)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			http.Redirect(w, r, "/login?messagetitle=User Unauthorized&messagebody=Please login to use this site", http.StatusSeeOther)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, "/login?messagetitle=User Unauthorized&messagebody=Please login to use this site", http.StatusSeeOther)
		return
	}
	sessionToken := c.Value

	// remove the users session from the session map
	delete(sessions, sessionToken)

	// We need to let the client know that the cookie is expired
	// In the response, we set the session token to an empty
	// value and set its expiry as the current time
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
			//Redirect Login
			http.Redirect(w, r, "/login?messagetitle=Logout Successful&messagebody=You have successfully been logged out", http.StatusSeeOther)
}
