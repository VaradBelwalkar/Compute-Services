package query_handling

import (

	"fmt"
	"context"
	"net/http"
	"github.com/VaradBelwalkar/Private-Cloud/api/container_apis"
	"github.com/VaradBelwalkar/Private-Cloud/api/database_handling"
	"path/filepath"
	"golang.org/x/crypto/bcrypt"
	"github.com/docker/docker/client"
	"go.mongodb.org/mongo-driver/bson"
)
// Here the user will be authenticated first and then request will be fulfilled


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
	err = database_handling.CollectionHandler.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err == nil {
		http.Error(w, "Username already exists!", http.StatusBadRequest)
	} else{			// Here if error is not nil, means document is not found, so free to create new document for the user

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the new user into the database
	_, err = database_handling.CollectionHandler.InsertOne(context.TODO(), bson.M{"username": username, "hashed_password": hashedPassword})
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
	err = database_handling.CollectionHandler.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err == nil {
		
			// Hash the password
		_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Remove the user document from the user_details

		_,err = database_handling.CollectionHandler.DeleteOne(context.TODO(),bson.M{"username": username})
		if err!=nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)	


	} else{

		http.Error(w, "Username doesn't exist!", http.StatusBadRequest)
}
}
