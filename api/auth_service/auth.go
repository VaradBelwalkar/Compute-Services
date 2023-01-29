package auth_service
import(
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"github.com/VaradBelwalkar/database_handling"
	"github.com/VaradBelwalkar/auth_service"
)



var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
//RECHECK
func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}





// login handler

func LoginHandler(response http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		// .. check credentials against db entry
        check,err:=Authenticate_user(username,pass)
		if err!=nil{
			return err;
		}
		if check == false{
			fmt.Fprintf(response,"Invalid credentials")
			return 
		}
		//Handle JWT signing and header creation 
		token,err:=SignHandler(username)
		if err!=nil{
			return 
		}
		tokenString = "Bearer "+token
		w.Header().Set("Authorization",tokenString)
		fmt.Fprintf(w,"Authentication Successful")

        //setting cookie based session
		CreateSession(response,username)
		//redirectTarget = "/internal"
	}
	//http.Redirect(response, request, redirectTarget, 302)
}

// logout handler

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
	sessionID,session_username,err:= RetrieveSession(request)
	if err!=nil || session_username == ""{
		fmt.Printf(w,"Logout Failed")
		return
	}
	username,err:=VerifyHandler(request)
	if err!=nil || username == ""{
		fmt.Fprintf(w,"Invalid JWT")
		return
	}
	if session_username!=username{
		fmt.Fprintf(w,"Invalid Cookie!")
	}

	DeleteSession(sessionID)
	//http.Redirect(response, request, "/", 302)
}

// index page

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
    <label for="name">User name</label>
    <input type="text" id="name" name="name">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	// Add CSRF

	fmt.Fprintf(response, indexPage)
}

// internal page

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		fmt.Fprintf(response, internalPage, userName)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}