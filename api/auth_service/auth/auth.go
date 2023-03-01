package auth
import(
	"net/http"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"
	twofa "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/twofa"
	"github.com/gorilla/securecookie"
)




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