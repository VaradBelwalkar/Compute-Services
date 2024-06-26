package auth_service

import (
	"net/http"
	"encoding/json"
	"errors"
	rds "github.com/VaradBelwalkar/Compute-Services/api/database_handling/redis"
)


func savePassResetSession(sessionID string,username string,OTP string)error{
	err:=rds.Redis_Set_Value_With_Timeout(sessionID,username,5)
	if err!=true{
		return errors.New("errorHolder")
	}
	val:=make(map[string]string)
	val["Authentication"]="passreset"
	val["JWT"]="issued"
	val["OTP"]=OTP
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

func CreatePassResetSession(w http.ResponseWriter,username string,OTP string){
		// Create a new session
		sessionID := generateSessionID(10)
		// Save the session

		err:=savePassResetSession(sessionID,username,OTP)
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
func RetrievePassResetSession(r *http.Request) (string,string,string,error){
	// Get the session ID from the request
	sessionID, err := r.Cookie("session")
	if err != nil {
		return "","","",err // Here the error means cookie doesn't exist
	}
	// Get the session by ID

	username:=rds.Redis_Get_Value(sessionID.Value)
	if username == ""{
		return "","","",errors.New("errorHolder")
	}

	UserInstance:=make(map[string]string)
    jsonString:=rds.Redis_Get_Value(username)
    err = json.Unmarshal([]byte(jsonString), &UserInstance)
    if err != nil {
        return "","","",nil
    }
	if UserInstance["Authentication"] == "passreset"{
		return sessionID.Value,username,UserInstance["OTP"],nil
	}else{
		return "","","",errors.New("errorHolder")}
}

