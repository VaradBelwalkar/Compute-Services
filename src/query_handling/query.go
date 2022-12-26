package query_handling

import (

	"fmt"
	"net/http"
	"github.com/VaradBelwalkar/Private-Cloud/src/container_apis"
	"github.com/VaradBelwalkar/Private-Cloud"
	"path/filepath"

)
// Here the user will be authenticated first and then request will be fulfilled


// HandlerFunc to be registered
func Container_Run(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := _.URL.Path // Suppose "/foo/bar/something"
	imageName:=filepath.Base(URLPath) // will output "something"
	//Get the requested Image from from the request-URL and pass it to the Container handler
	ContainerCreate(ctx,Cli,imageName)


  }
  
  func Container_Resume(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := _.URL.Path // Suppose "/foo/bar/something"
	containerName:=filepath.Base(URLPath) // will output "something"	
	ContainerStart(ctx,Cli,containerName)
  }
  
  func Container_List(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	OwnedContainerInfo(ctx,)
  
  }

  func Container_Stop_or_Remove(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := _.URL.Path // Suppose "/foo/bar/something"
	imageName:=filepath.Base(URLPath) // will output "something" 
  
  }

  func Container_List(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
  
  }

  func upload_file(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()


  }

  func upload_folder(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()

  }

  func registerUser(w http.ResponseWriter, r *http.Request) {
	// Parse the POST request body and retrieve the form values
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	// Validate the form values
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	// Check if a document with the given username already exists
	var result bson.M
	err = CollectionHandler.FindOne(ctx, bson.M{"username": username}).Decode(&result)
	if err == nil {
		http.Error(w, "Username already exists!", http.StatusBadRequest)
	} else{

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the new user into the database
	_, err = CollectionHandler.InsertOne(ctx, bson.M{"username": username, "hashed_password": hashedPassword})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}
}

// Parse and authenticate the user and then remove the account from the database
func remove_account(w http.ResponseWriter, r *http.Request) {
	// Parse the POST request body and retrieve the form values
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	// Validate the form values
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	// Check if a document with the given username already exists
	var result bson.M
	err = CollectionHandler.FindOne(ctx, bson.M{"username": username}).Decode(&result)
	if err == nil {
		
			// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else{

		http.Error(w, "Username doesn't exist!", http.StatusBadRequest)

	// Return a success response
	w.WriteHeader(http.StatusOK)
}
}
