package sendMail

import (
	"net/smtp"
	 "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/twofa"
)
import "net/http"

//route handler
func SendMail(username string,eMail string) int{

    if eMail == ""{   
   return http.StatusBadRequest
}
    toEmailAddress := eMail //Should be dynamic generated format
    to := []string{toEmailAddress}

    host := "smtp.gmail.com"
    port := "587"
    address := host + ":" + port 

	msg := []byte("To: "+toEmailAddress+"\r\n" +
		"Subject: Your OTP for Two-Factor Authentication\n\n\r\n" +
		"\r\n" +
		"Dear User,\n\nYour OTP for two-factor authentication is: " + "\n\nPlease enter this otp in your app to complete the authentication process.\n\nBest regards,\nDYPLUG\r\n")

    //subject := "Your OTP for Two-Factor Authentication\n\n"
    //body := "Dear User,\n\nYour OTP for two-factor authentication is: " + OTP + "\n\nPlease enter this otp in your app to complete the authentication process.\n\nBest regards,\nDYPLUG"
    //message := []byte(subject + body)                           //Don't use colon(:)

    auth := smtp.PlainAuth("", twofa.Official_Email, twofa.Official_Email_Password, host)

    err := smtp.SendMail(address, auth, twofa.Official_Email, to, msg)
    if err != nil {
        return http.StatusInternalServerError
    }
return http.StatusOK
}