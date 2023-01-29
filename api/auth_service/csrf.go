package auth_service

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"text/template"
)
//Here the csrf token is hidden in the form which will be submitted by user in its session
const csrfTemplate = `
<html>
<body>
<form action="/login" method="POST">
  <label for="name">User name</label>
  <input type="text" id="name" name="name">
  <label for="password">Password</label>
  <input type="password" id="password" name="password">
  <input type="hidden" name="csrf" value="{{.csrf}}">
  <input type="submit" value="Submit">
</form>
</body>
</html>
`

// Generate a random CSRF token
func generateCSRFToken() (string, error) {
	// Generate a random token
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}


// Render the form
func RenderForm(w http.ResponseWriter, r *http.Request) {
	// Generate a CSRF token
	csrfToken, err := generateCSRFToken()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the CSRF token in the user's session (will will use this to check against token in the submitted form)
	session, err := r.Cookie("session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["csrf"] = csrfToken
	http.SetCookie(w, session)

	// Render the form
	tmpl, err := template.New("form").Parse(csrfTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//tmpl object expects object  of type map[string]string to update values in the provided form tempate
	data := map[string]string{
		"csrf": csrfToken,
	}
	//Here The Execute command is first updating the form with data object  i.e updating the fields the tmpl can 
	//understand i.e {{.csrf}} field and then updating their values and sending the template to frontend
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}



func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	// Check the request method
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get the CSRF token from the request
	csrfToken := r.FormValue("csrf")

	// Get the CSRF token from the user's session
	session, err := r.Cookie("session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sessionCSRFToken, ok := session["csrf"]
	if !ok {
		http.Error(w, "Invalid CSRF token", http.StatusBadRequest)
		return
	}

	// Compare the CSRF tokens
	if csrfToken != sessionCSRFToken {
		http.Error(w, "Invalid CSRF token", http.StatusBadRequest)
		return
	}
	//Returns boolean value whether user is credible or not
	check,err:=LoginHandler(w,r)


	// If the CSRF tokens match, process the form submission
	// ...
}
