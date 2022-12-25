package main

import (
	"fmt"
	"context"
	"io"
	"os"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
    "src/database_handling"
    "src/query_handling"
    "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)
//Global Object 
var CollectionHandler *Collection 

// The main function manages all the query handling and manages the database as well
const url = "mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=myRS"

// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		// .. check credentials ..
		setSession(name, response)
		redirectTarget = "/internal"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
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

// server main method

var router = mux.NewRouter()

func main() {
    //Initiate Mongo client
	mongo_client, err := mongo.NewClient(options.Client().ApplyURI(url)) // Give appropriate port here
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	err = mongo_client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mongo_client.Disconnect(ctx)



// Initiate Docker client
Cli_native, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
if err != nil {
    panic(err)
}
defer cli.Close()


InitiateMongoDB(mongo_client);

// Get userinfo-collection handler 
CollectionHandler = mongo_client.Database("private-cloud").Collection("userinfo")



	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/internal", internalPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/container/run/*", Container_Run)
	router.HandleFunc("/container/resume/*", Container_Resume)
	router.HandleFunc("/container/<regex>/*", Container_Stop_or_Remove)
	router.HandleFunc("/container/list/*", Container_List)
	router.HandleFunc("/upload_file/", upload_file)
	router.HandleFunc("/upload_folder/", upload_folder)

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}