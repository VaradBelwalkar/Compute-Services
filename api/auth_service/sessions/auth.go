package auth_service

import (
	"net/http"
	"math/rand"
	"encoding/json"
	"time"
	"errors"
	rds "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/redis"
)



// Save a session (only occurs when new session is going to be created)
func saveSession(sessionID string,username string) error{
	err:=rds.Redis_Set_Value_With_Timeout(sessionID,username,1440)
	if err!=true{
		return errors.New("errorHolder")
	}
	val:=make(map[string]string)
	val["Authentication"]="success"
	val["JWT"]="issued"
	jsonFormat, chk := json.Marshal(val)
    if chk != nil {
        return errors.New("errorHolder")
    }
	err=rds.Redis_Set_Value_With_Timeout(username,string(jsonFormat),1440)
	if err!=true{
		return errors.New("errorHolder")
	}
	return nil	
}


// A handler function that creates a new session and saves it
func CreateSession(w http.ResponseWriter,username string) {
	// Create a new session
	rand.Seed(time.Now().UnixNano())
	sessionID := generateSessionID(10)
	// Save the session
	err:=saveSession(sessionID,username)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Set the session ID as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: sessionID,
	})
	w.WriteHeader(http.StatusOK)
}


func RetrieveAuthorizedSession(r *http.Request) (string,string,bool,error){
	// Get the session ID from the request
	sessionID, err := r.Cookie("session")
	if err != nil {
		return "","",false,err // Here the error means cookie doesn't exist
	}
	// Get the session by ID

	username:=rds.Redis_Get_Value(sessionID.Value)
	if username == ""{
		return "","",false,errors.New("errorHolder")
	}

    UserInstance:=make(map[string]string)
    jsonString:=rds.Redis_Get_Value(username)
    err = json.Unmarshal([]byte(jsonString), &UserInstance)
    if err != nil {
        return "","",false,nil
    }
	if UserInstance["Authentication"] == "success"{
		return sessionID.Value,username,true,nil
	}else{
	return sessionID.Value,username,false,nil}

}

