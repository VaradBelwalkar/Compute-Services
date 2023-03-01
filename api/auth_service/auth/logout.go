package auth
import(
	"net/http"
	jwt "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/jwt"
	ssn "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/sessions"
)


// logout handler

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID,session_username,chk,err:= ssn.RetrieveAuthorizedSession(r)
	if err!=nil || session_username == "" || chk!=true{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	username,status:=jwt.VerifyHandler(r)
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

	ssn.DeleteSession(sessionID)

	w.WriteHeader(http.StatusOK)
	//http.Redirect(response, request, "/", 302)
}

