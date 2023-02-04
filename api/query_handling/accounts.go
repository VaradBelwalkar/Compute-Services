package query_handling

import (

	"context"
	"net/http"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
	as "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
)
// Here the user will be authenticated first and then request will be fulfilled


func RegisterUser(w http.ResponseWriter, r *http.Request) {

	//CSRF handling
	check:=as.HandleSubmit(w,r)
	if check!=true{
		return
	}

	// Parse the POST request body and retrieve the form values
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)  // Here status code is 400
		return
	}


	// Validate the form values
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check if a document with the given username already exists
	var result bson.M
	err = db.CollectionHandler.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err == nil {
		w.WriteHeader( http.StatusConflict)
	} else{			// Here if error is not nil, means document is not found, so free to create new document for the user

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Insert the new user into the database
	_, err = db.CollectionHandler.InsertOne(context.TODO(), bson.M{"username": username, "hashed_password": hashedPassword})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)  
}
}

// Parse and authenticate the user and then remove the account from the database
func RemoveAccount(w http.ResponseWriter, r *http.Request) {

	//CSRF handling
	check:=as.HandleSubmit(w,r)
	if check!=true{
		return
	}
	
	// Parse the POST request body and retrieve the form values
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	// Validate the form values
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check if a document with the given username already exists
	var result bson.M
	err = db.CollectionHandler.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err == nil {
		
			// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Check if the password matches
		if string(hashedPassword[:]) == result["hashed_password"]{
		// Remove the user document from the user_details

		_,err = db.CollectionHandler.DeleteOne(context.TODO(),bson.M{"username": username})
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)} else{
			w.WriteHeader(http.StatusForbidden)
		}	

	} else{

		w.WriteHeader(http.StatusNotFound)
}
}
