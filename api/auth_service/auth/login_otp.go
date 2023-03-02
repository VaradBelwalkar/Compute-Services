package auth
import(
	"net/http"
	twofa "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/twofa"
	csrf "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/csrf"
	jwt "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/jwt"
	ssn "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/sessions"
)


//OTP Handler 

// Don't need to check auth during GET for otphandler as the form request doesn't require it
func OTPHandler(w http.ResponseWriter, r *http.Request){

		//CSRF handling
		check:=csrf.HandleSubmit(w,r)
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



		chk:=twofa.TwoFA_Verify(username,OTP)					//Deletes the <username> temporary key when true
		if chk!=true{
			w.WriteHeader(http.StatusUnauthorized)			
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
		ssn.CreateSession(w,username)					//Creates new session valid for 24hrs
		//redirectTarget = "/internal"
	}
	//http.Redirect(response, request, redirectTarget, 302)
	
}


