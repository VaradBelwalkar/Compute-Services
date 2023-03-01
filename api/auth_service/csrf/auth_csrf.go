package csrf
import (
	"net/http"
	jwt "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/jwt"
	rds "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/redis"

)

func chkCsrf(givenCsrfToken string,givenUsername string)bool{

	actUname:=rds.Redis_Get_Value(givenCsrfToken)
	if givenUsername == actUname{
		return true
	}else{
		return false
	}
}


// This handles the CSRF submission
func HandleAuthSubmit(w http.ResponseWriter, r *http.Request) bool {
	// Check the request method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)   // 405
		return false
	}
	username,chk:=jwt.VerifyHandler(r)
	if chk!=200{
		w.WriteHeader(http.StatusPreconditionFailed) 
		return false
	}
	// Get the CSRF token from the request
	csrfToken := r.FormValue("csrf")

	// Get the CSRF token from the user's session
	sessionCSRFToken, err := r.Cookie("csrftoken")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)   // 400
		return false
	}

	// Compare the CSRF tokens
	if csrfToken != sessionCSRFToken.Value {
		w.WriteHeader(http.StatusPreconditionFailed)   // 412 precondition i.e csrf is failed
		return false
	}
	
	check:=chkCsrf(csrfToken,username)
	if check!=true{
		w.WriteHeader(http.StatusUnauthorized) 
		return false
	}

	//Returns boolean value whether user is credible or not
	return true

}