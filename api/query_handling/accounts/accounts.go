package accounts

import (

	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	mndb "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"
	auth "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/auth"
	csrf "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/csrf"
	jwt "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/jwt"
	twofa "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/twofa"
	session "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/sessions"
	containers "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/query_handling/containers"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/docker/docker/api/types/volume"
)
// Here the user will be authenticated first and then request will be fulfilled

type Config struct {
	Emails []string `json:"emails"`
}

func isAllowed(email string) bool{
		// read the config file
		data, err := ioutil.ReadFile("config.json")
		if err != nil {
			return false
		}
	
		// unmarshal the JSON data into a Config struct
		var config Config
		if err := json.Unmarshal(data, &config); err != nil {
			return false
		}
	
		found := false
		for _, e := range config.Emails {
			if e == email {
				found = true
				break
			}
		}
	
		if found {
			return true
		} else {
			return false
		}
}





func RegisterUser(w http.ResponseWriter, r *http.Request) {

	//CSRF handling
	check:=csrf.HandleSubmit(w,r)
	if check!=true{
		return
	}
	
	// Parse the POST request body and retrieve the form values
	err := r.ParseForm()
	if err != nil {		
		w.WriteHeader(http.StatusBadRequest)  // Here status code is 400 something went wrong at your side!
		return
	}


	// Validate the form values
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	EMAIL:=r.Form.Get("email")

	chk:=isAllowed(EMAIL)
	if chk !=true{
		w.WriteHeader(http.StatusFailedDependency) // 424 Statuscode
		return
	}

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest) // Here status code is 400 something went wrong at your side!
		return
	}
	// Check if a document with the given username already exists
	var result bson.M
	err = mndb.CollectionHandler.FindOne(context.TODO(), bson.M{"email": EMAIL}).Decode(&result)	//Cannot create account if email exists
	if err == nil {
		w.WriteHeader( http.StatusNotAcceptable)  //406 statuscode
		return
	} else{			// Here if error is not nil, means document is not found, so free to create new document for the user

		err = mndb.CollectionHandler.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result) //Check for username conflicts
		if err == nil {
			w.WriteHeader( http.StatusConflict)  //409 statuscode
			return
		} else{	

			chk,OTP:=twofa.TwoFA_Send(username,EMAIL)
			if chk!=true{
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			//Handle JWT signing and header creation 
			token,err:=jwt.SignHandler(username)
			if err!=nil{ 	
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			tokenString := "Bearer "+token
			w.Header().Set("Authorization",tokenString)
			
			//setting cookie based session
			session.CreateRegTempSession(w,username,password,EMAIL,OTP)
			return
 
 		}
}
}


func VerifyRegisterUser(w http.ResponseWriter, r *http.Request){

			//CSRF handling
			check:=csrf.HandleSubmit(w,r)
			if check!=true{
				return
			}
			check,username,password,EMAIL:=auth.Temp_Reg_auth(w,r)
			if check!=true{
				return
			}
		
		OTP := r.FormValue("otp")
	
		if OTP == "" {
			w.WriteHeader(http.StatusBadRequest) // 400 
		} else {
	
	
	
			chk:=twofa.TwoFA_Verify(username,OTP)					//Deletes the <username> temporary key when true
			if chk!=true{
				w.WriteHeader(http.StatusUnauthorized)			
				return
			}

	volumeOpts:=volume.CreateOptions{
		Name:username,
		Driver:"local",
	}

	_,err:=containers.Cli.VolumeCreate(context.TODO(),volumeOpts)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
// Insert the new user into the database
	hashedPassword:=mndb.ComputeHash(password)
_, err = mndb.CollectionHandler.InsertOne(context.TODO(), bson.M{"username": username, "password": hashedPassword,"email":EMAIL})
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
	check:=csrf.HandleSubmit(w,r)
	if check!=true{
		return
	}
	
	// Parse the POST request body and retrieve the form values
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}


	// Validate the form values
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	// Check if a document with the given username already exists
	type resultStruct struct{
		Username string
		Password string
	}

	result:=resultStruct{}
	err = mndb.CollectionHandler.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err == nil {

		//Check if the password matches
		if password == result.Password{
		// Remove the user document from the user_details

		err:=containers.Cli.VolumeRemove(context.TODO(),username,true)
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_,err = mndb.CollectionHandler.DeleteOne(context.TODO(),bson.M{"username": username})
		if err!=nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)} else{
			w.WriteHeader(http.StatusUnauthorized)   // 401
			return
		}	

	} else{

		w.WriteHeader(http.StatusNotFound)  //   404
}
}
