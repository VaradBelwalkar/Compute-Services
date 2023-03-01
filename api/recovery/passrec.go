package recovery

import(
	"net/http"
	csrf "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/csrf"
	jwt "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/jwt"
	twofa "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/twofa"
	session "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/sessions"
)

func RecoverPass(w http.ResponseWriter, r *http.Request) {

	//CSRF handling
	check:=csrf.HandleSubmit(w,r)
	if check!=true{
		return
	}

username := r.FormValue("username")
email := r.FormValue("email")

if username == "" || email == ""{
	w.WriteHeader(http.StatusBadRequest) // 400 
} else {
	actEmail:=twofa.RetrieveEmail(username)
	if email!=actEmail{
		w.WriteHeader(http.StatusUnauthorized)
		return 
	}


	chk,OTP:=twofa.TwoFA_Send(username,actEmail)
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
	session.CreatePassResetSession(w,username,OTP)
	return
	//redirectTarget = "/internal"
}
//http.Redirect(response, request, redirectTarget, 302)
}



func RecoverPassCheck(w http.ResponseWriter, r *http.Request) {

}