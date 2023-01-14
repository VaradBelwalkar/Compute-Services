package main

import (
	"fmt"
	"net/http"
)



// ******************** Here understand how sessions work in Golang ********************




//This file doesn't contribute any code to the package, just for learning purpose




// A struct that represents a user session
type Session struct {
	ID     string
	Values map[string]interface{}
}

// A map to store the sessions
var sessions = map[string]*Session{}

// Generate a new session ID
func generateSessionID() string {
	// Generate a random session ID
	// You can use a library like "crypto/rand" to generate a random ID
	return "12345"
}

// Create a new session
func newSession() *Session {
	return &Session{
		ID:     generateSessionID(),
		Values: map[string]interface{}{},
	}
}

// Save a session
func saveSession(session *Session) {
	sessions[session.ID] = session
}

// Get a session by ID
func getSession(id string) (*Session, error) {
	session, ok := sessions[id]
	if !ok {
		return nil, fmt.Errorf("Session not found")
	}
	return session, nil
}

// A handler function that creates a new session and saves it
func createSession(w http.ResponseWriter, r *http.Request) {
	// Create a new session
	session := newSession()

	// Save the session
	saveSession(session)

	// Set the session ID as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: session.ID,
	})

	// Write a response
	w.Write([]byte("Session created"))
}

// A handler function that retrieves a session by ID
func retrieveSession(w http.ResponseWriter, r *http.Request) {
	// Get the session ID from the request
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the session by ID
	session, err := getSession(sessionID.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write a response
	w.Write([]byte("Session retrieved"))
}

func main() {
	http.HandleFunc("/create", createSession)
	http.HandleFunc("/retrieve", retrieveSession)
	http.ListenAndServe(":8080", nil)
}





//Explaination ->
//Sessions are nothing but some data structures stored in the memory by the server to retrieve session info
//Like in above scenario, the session is nothing but,
//Whenever a user visits the website first, a session is created i.e a session ID is generated with empty session values 
//Then this session is stored in a Map interface to retrieve it the next time a request with the SAME session.ID comes in
//The above code explains everything

//Here we are passing the session ID to the frontend as a cookie, so frontend needs to save that cookie and 
//need to use that sessionID cookie for further requests so that server can retrieve session info 
//This way csrf is prevented like follow

//First time user visits the server, server creates a unique sessionID and a CSRF token, and saves that 
//session in a session Map and embeds the CSRF in hidden form field. Also it passes the sessionID as cookie to frontend

//The frontend uses this sessionID cookie to make the post request after filling the form, 
//When server receives the POST request, it first tries to get the session cookie, if it doesn't found, it rejects the form 
//Else when it finds a sessionID cookie, it then searches for the sessionID in session Map, and if found,
//Checks the CSRF token value from the stored session to the CSRF received from the form,
//If match is found, it processes the form, else rejects it
