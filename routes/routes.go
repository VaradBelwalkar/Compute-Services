package routes

import(
	"github.com/gorilla/mux"
	auth "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/auth"
	csrf "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/csrf"
	containers "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/query_handling/containers"
	"github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/query_handling/accounts"
	"github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/recovery"
)

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/", csrf.RenderForm)
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")	
	router.HandleFunc("/login", csrf.RenderForm).Methods("GET")	
	router.HandleFunc("/logout", auth.LogoutHandler).Methods("POST")	
	router.HandleFunc("/otphandler", csrf.RenderForm).Methods("GET")
	router.HandleFunc("/otphandler", auth.OTPHandler).Methods("POST")
	router.HandleFunc("/recoverpass", csrf.RenderForm).Methods("GET")
	router.HandleFunc("/recoverpass", recovery.RecoverPassCheck).Methods("POST")
	router.HandleFunc("/recoverpasscheck", recovery.RecoverPassCheck).Methods("POST")
	router.HandleFunc("/recoverpasscheck", recovery.ChangePass).Methods("POST")
	router.HandleFunc("/recoverpass", csrf.RenderForm).Methods("POST")
	router.HandleFunc("/container/run/{image}", containers.Container_Run).Methods("GET")
	router.HandleFunc("/container/run/{image}", containers.Container_Run).Methods("POST")
	router.HandleFunc("/container/resume/{container}", containers.Container_Resume).Methods("GET")
	router.HandleFunc("/container/stop/{container}", containers.Container_Stop).Methods("GET")
	router.HandleFunc("/container/remove/{container}", containers.Container_Remove).Methods("GET")
	router.HandleFunc("/container/list/containers", containers.Container_List).Methods("GET")
    router.HandleFunc("/register",csrf.RenderForm).Methods("GET")
	router.HandleFunc("/register",accounts.RegisterUser).Methods("POST") 
	router.HandleFunc("/regotphandler",accounts.VerifyRegisterUser).Methods("POST") 
	router.HandleFunc("/remove_account",csrf.RenderForm).Methods("GET")	
	router.HandleFunc("/remove_account",accounts.RemoveAccount).Methods("POST")  
}