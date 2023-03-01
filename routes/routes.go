package routes

import(
    db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
	rec "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/recovery"
    qh "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/query_handling"
	as "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service"
)

func addSignHandler(router *mux.Router) {
	router.HandleFunc("/", as.RenderForm)
	router.HandleFunc("/login", as.LoginHandler).Methods("POST")	
	router.HandleFunc("/login", as.RenderForm).Methods("GET")	
	router.HandleFunc("/logout", as.LogoutHandler).Methods("POST")	
	router.HandleFunc("/otphandler", as.RenderForm).Methods("GET")
	router.HandleFunc("/otphandler", as.OTPHandler).Methods("POST")
	router.HandleFunc("/recoverpass", as.RenderForm).Methods("GET")
	router.HandleFunc("/recoverpass", as.RecoverPass).Methods("POST")
	router.HandleFunc("/recoverpasscheck", as.RenderForm).Methods("GET")
	router.HandleFunc("/recoverpasscheck", as.RecoverPassCheck).Methods("POST")
	router.HandleFunc("/recoverpass", as.RenderForm).Methods("POST")
	router.HandleFunc("/container/run/{image}", qh.Container_Run).Methods("GET")
	router.HandleFunc("/container/run/{image}", qh.Container_Run).Methods("POST")
	router.HandleFunc("/container/resume/{container}", qh.Container_Resume).Methods("GET")
	router.HandleFunc("/container/stop/{container}", qh.Container_Stop).Methods("GET")
	router.HandleFunc("/container/remove/{container}", qh.Container_Remove).Methods("GET")
	router.HandleFunc("/container/list/containers", qh.Container_List).Methods("GET")
    router.HandleFunc("/register",as.RenderForm).Methods("GET")
	router.HandleFunc("/register",qh.RegisterUser).Methods("POST") 
	router.HandleFunc("/regotphandler",qh.VerifyRegisterUser).Methods("POST") 
	router.HandleFunc("/remove_account",as.RenderForm).Methods("GET")	
	router.HandleFunc("/remove_account",qh.RemoveAccount).Methods("POST")  
}