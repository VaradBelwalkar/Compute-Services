package auth_service
import(
	"net/http"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
	"github.com/gorilla/securecookie"
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





// login handlersession_username

func LoginHandler(w http.ResponseWriter, r *http.Request) {

		//CSRF handling
		check:=HandleSubmit(w,r)
		if check!=true{
			return
		}

	username := r.FormValue("username")
	pass := r.FormValue("password")

	if username == "" || pass == ""{
		w.WriteHeader(http.StatusBadRequest) // 400 
	} else {
		// .. check credentials against db entry
        check:=db.Authenticate_user(username,pass)
		if check !=200{
			if check == 500{
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else if check == 404{
				w.WriteHeader(http.StatusNotFound)
				return
			} else if check == 401{
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			return
		}

		chk,OTP:=TwoFA_Send(username,"")
		if chk!=true{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//Handle JWT signing and header creation 
		token,err:=SignHandler(username)
		if err!=nil{ 	
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tokenString := "Bearer "+token
		w.Header().Set("Authorization",tokenString)
		
        //setting cookie based session
		CreateTempSession(w,username,OTP)
		return
		//redirectTarget = "/internal"
	}
	//http.Redirect(response, request, redirectTarget, 302)
}


//OTP Handler 

// Don't need to check auth during GET for otphandler as the form request doesn't require it
func OTPHandler(w http.ResponseWriter, r *http.Request){

		//CSRF handling
		check:=HandleSubmit(w,r)
		if check!=true{
			return
		}
		check,username:=Temp_auth(w,r)
		if check!=true{
			return
		}
	
	OTP := r.FormValue("otp")

	if OTP == "" {
		w.WriteHeader(http.StatusBadRequest) // 400 
	} else {



		chk:=TwoFA_Verify(username,OTP)					//Deletes the <username> temporary key when true
		if chk!=true{
			w.WriteHeader(http.StatusUnauthorized)			
			return
		}
		//Handle JWT signing and header creation 
		token,err:=SignHandler(username)		
		if err!=nil{ 	
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tokenString := "Bearer "+token
		w.Header().Set("Authorization",tokenString)
		
        //setting cookie based session
		CreateSession(w,username)					//Creates new session valid for 24hrs
		//redirectTarget = "/internal"
	}
	//http.Redirect(response, request, redirectTarget, 302)
	
}




// logout handler

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID,session_username,chk,err:= RetrieveAuthorizedSession(r)
	if err!=nil || session_username == "" || chk!=true{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	username,status:=VerifyHandler(r)
	if status!=200 || username == ""{
		if status == 401{
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if status == 500{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if session_username!=username{
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	DeleteSession(sessionID)

	w.WriteHeader(http.StatusOK)
	//http.Redirect(response, request, "/", 302)
}



func Temp_auth(w http.ResponseWriter, r *http.Request) (bool,string) {
	_,session_username,err:= RetrieveTempSession(r)
	if err!=nil || session_username == ""{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,""
	}
	username,status:=VerifyHandler(r)
	if status!=200 || username == ""{	
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,""
	}
	if session_username!=username{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,""
	}
	return true,username

}


func Temp_Reg_auth(w http.ResponseWriter, r *http.Request) (bool,string,string,string) {
	_,session_username,password,email,err:= RetrieveRegTempSession(r)
	if err!=nil || session_username == ""{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,"","",""
	}
	username,status:=VerifyHandler(r)
	if status!=200 || username == ""{	
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,"","",""
	}
	if session_username!=username{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,"","",""
	}
	return true,username,password,email

}



func Verify_Auth(w http.ResponseWriter, r *http.Request)(bool,string){
	_,session_username,chk,err:= RetrieveAuthorizedSession(r)
	if err!=nil || session_username == "" || chk!=true{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,""
	}
	username,status:=VerifyHandler(r)
	if status!=200 || username == ""{	
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,""
	}
	if session_username!=username{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,""
	}
	return true,username
}