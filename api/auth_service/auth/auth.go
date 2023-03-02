package auth
import(
	"net/http"
	jwt "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/jwt"
	ssn "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/sessions"
)



// For authenticating user against OTP after successful password confirmation
func Temp_auth(w http.ResponseWriter, r *http.Request) (bool,string) {
	_,session_username,err:= ssn.RetrieveTempSession(r)
	if err!=nil || session_username == ""{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,""
	}
	username,status:=jwt.VerifyHandler(r)
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

// For authenticating user against OTP using temporary session created after user checked against already available  
func Temp_Reg_auth(w http.ResponseWriter, r *http.Request) (bool,string,string,string) {
	_,session_username,password,email,err:= ssn.RetrieveRegTempSession(r)
	if err!=nil || session_username == ""{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,"","",""
	}
	username,status:=jwt.VerifyHandler(r)
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

// For authenticating user against OTP using temporary session created after user checked against already available  
func Temp_Pass_auth(w http.ResponseWriter, r *http.Request) (bool,string,string) {
	_,session_username,OTP,err:= ssn.RetrievePassResetSession(r)
	if err!=nil || session_username == ""{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,"",""
	}
	username,status:=jwt.VerifyHandler(r)
	if status!=200 || username == ""{	
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,"",""
	}
	if session_username!=username{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,"",""
	}
	return true,username,OTP

}


func Verify_Pass_auth(w http.ResponseWriter, r *http.Request) (bool,string) {
	_,session_username,err:= ssn.RetrievePassAuthSession(r)
	if err!=nil || session_username == ""{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,""
	}
	username,status:=jwt.VerifyHandler(r)
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


func Verify_Auth(w http.ResponseWriter, r *http.Request)(bool,string){
	_,session_username,chk,err:= ssn.RetrieveAuthorizedSession(r)
	if err!=nil || session_username == "" || chk!=true{
		w.WriteHeader(http.StatusUnauthorized) // 401 meaning user should login again
		return false,""
	}
	username,status:=jwt.VerifyHandler(r)
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