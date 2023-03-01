package auth_service

import (
	"net/http"
	"encoding/json"
	"errors"
	rds "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/redis"
)


func saveRegTempSession(sessionID string,username string,password string,email string,OTP string)error{
	err:=rds.Redis_Set_Value_With_Timeout(sessionID,username,5)
	if err!=true{
		return errors.New("errorHolder")
	}
	val:=make(map[string]string)
	val["Authentication"]="register"
	val["JWT"]="issued"
	val["OTP"]=OTP
	val["Password"]=password
	val["Email"]=email
	jsonFormat, chk := json.Marshal(val)
    if chk != nil {
        return errors.New("errorHolder")
    }

	err=rds.Redis_Set_Value_With_Timeout(username,string(jsonFormat),5)
	if err!=true{
		return errors.New("errorHolder")
	}
	return nil	
}

func CreateRegTempSession(w http.ResponseWriter,username string,password string,email string,OTP string){
		// Create a new session
		sessionID := generateSessionID(10)
		// Save the session

		err:=saveRegTempSession(sessionID,username,password,email,OTP)
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
		return 
}

// A handler function that retrieves a session by ID
func RetrieveRegTempSession(r *http.Request) (string,string,string,string,error){
	// Get the session ID from the request
	sessionID, err := r.Cookie("session")
	if err != nil {
		return "","","","",err // Here the error means cookie doesn't exist
	}
	// Get the session by ID

	username:=rds.Redis_Get_Value(sessionID.Value)
	if username == ""{
		return "","","","",errors.New("errorHolder")
	}

	UserInstance:=make(map[string]string)
    jsonString:=rds.Redis_Get_Value(username)
    err = json.Unmarshal([]byte(jsonString), &UserInstance)
    if err != nil {
        return "","","","",nil
    }
	if UserInstance["Authentication"] == "register"{
		return sessionID.Value,username,UserInstance["Password"],UserInstance["Email"],nil
	}else{
		return "","","","",errors.New("errorHolder")}
}
