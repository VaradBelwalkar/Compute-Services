package recovery

import(
	"net/http"
	"context"
	csrf "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/csrf"
	jwt "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/jwt"
	mndb "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"
	twofa "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/twofa"
	"github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/auth"
	"go.mongodb.org/mongo-driver/bson"
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


//POST -> CHECK OTP
func RecoverPassCheck(w http.ResponseWriter, r *http.Request) {
		//CSRF handling
		check:=csrf.HandleSubmit(w,r)
		if check!=true{
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		chk,username,OTP:=auth.Temp_Pass_auth(w,r)
		if chk!=true{
			return
		}

		justCheck:=twofa.TwoFA_Verify(username,OTP)
		if justCheck!=true{
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		session.CreatePassAuthSession(w,username)
		

}



func ChangePass(w http.ResponseWriter, r *http.Request){
			//CSRF handling
			check:=csrf.HandleSubmit(w,r)
			if check!=true{
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			chk,username:=auth.Verify_Pass_auth(w,r)
			if chk!=true{
				w.WriteHeader(http.StatusUnauthorized)
			}


				// Parse the POST request body and retrieve the form values
	err := r.ParseForm()
	if err != nil {		
		w.WriteHeader(http.StatusBadRequest)  // Here status code is 400 something went wrong at your side!
		return
	}

	password := r.Form.Get("password")

	hashedPassword:=mndb.ComputeHash(password)
	filter:=bson.M{
		"username":username,
	}
	update:=bson.M{ "$unset":bson.M{
		"password":"",
	},
	"$set":bson.M{
		"password":hashedPassword,
	},
	}
	updateResult,ch:=mndb.CollectionHandler.UpdateOne(context.TODO(),filter,update)
	if ch!=nil || updateResult.MatchedCount!=1{
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}
	w.WriteHeader(http.StatusOK)

}