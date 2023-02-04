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





// login handler

func LoginHandler(w http.ResponseWriter, r *http.Request) {

		//CSRF handling
		check:=HandleSubmit(w,r)
		if check!=true{
			return
		}

	username := r.FormValue("username")
	pass := r.FormValue("password")

	if username != "" || pass != ""{
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// .. check credentials against db entry
        check:=db.Authenticate_user(username,pass)
		if check !=200{
			if check == 500{
				w.WriteHeader(http.StatusInternalServerError)
			} else if check == 404{
				w.WriteHeader(http.StatusNotFound)
			}
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
		CreateSession(w,username)
		//redirectTarget = "/internal"
	}
	//http.Redirect(response, request, redirectTarget, 302)
}

// logout handler

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID,session_username,err:= RetrieveSession(r)
	if err!=nil || session_username == ""{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	username,status:=VerifyHandler(r)
	if status!=200 || username == ""{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if session_username!=username{
		w.WriteHeader(http.StatusForbidden)
		return
	}

	DeleteSession(sessionID)
	//http.Redirect(response, request, "/", 302)
}



func Handle_auth(w http.ResponseWriter, r *http.Request) int {
	_,session_username,err:= RetrieveSession(r)
	if err!=nil || session_username == ""{
		w.WriteHeader(http.StatusNotFound)
		return 304
	}
	username,status:=VerifyHandler(r)
	if status!=200 || username == ""{	
		return 404
	}
	if session_username!=username{
		w.WriteHeader(http.StatusForbidden)
		return 404
	}
	
	return 200

}