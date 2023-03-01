package auth
import(
	"net/http"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"
	twofa "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/twofa"
	csrf "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/csrf"
	ssn "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/sessions"
	jwt "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/jwt"
)



// login handlersession_username

func LoginHandler(w http.ResponseWriter, r *http.Request) {

		//CSRF handling
		check:=csrf.HandleSubmit(w,r)
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

		chk,OTP:=twofa.TwoFA_Send(username,"")
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
		ssn.CreateTempSession(w,username,OTP)
		return
		//redirectTarget = "/internal"
	}
	//http.Redirect(response, request, redirectTarget, 302)
}

