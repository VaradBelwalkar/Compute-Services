package twofa

import (
	"net/smtp"
    "encoding/json"
    rds "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/redis"
)

func PassReset_Send(username string,eMail string)(bool,string) {
    OTP:=Generate_OTP()
    if eMail == ""{
    eMail=RetrieveEmail(username)
    if eMail == ""{
        return false,""
    }
}
    from := "officialdyplug@gmail.com"
    password := "oskpnzbzxwzkvpmu" //Change to get from the env
    toEmailAddress := eMail //Should be dynamic generated format
    to := []string{toEmailAddress}

    host := "smtp.gmail.com"
    port := "587"
    address := host + ":" + port 

	msg := []byte("To: "+toEmailAddress+"\r\n" +
		"Subject: Your OTP for Password Reset\n\n\r\n" +
		"\r\n" +
		"Dear User,\n\nYour OTP for two-factor authentication is: " + OTP + "\n\nPlease enter this otp in your app to complete the password reset process.\n\nRegards,\nDYPLUG\r\n")

    //subject := "Your OTP for Two-Factor Authentication\n\n"
    //body := "Dear User,\n\nYour OTP for two-factor authentication is: " + OTP + "\n\nPlease enter this otp in your app to complete the authentication process.\n\nBest regards,\nDYPLUG"
    //message := []byte(subject + body)                           //Don't use colon(:)

    auth := smtp.PlainAuth("", from, password, host)

    err := smtp.SendMail(address, auth, from, to, msg)
    if err != nil {
        return false,""
    }
return true,OTP
}


func PassReset_Verify(username string, SentOTP string) bool{
    UserInstance:=make(map[string]string)
    jsonString:=rds.Redis_Get_Value(username)
    err := json.Unmarshal([]byte(jsonString), &UserInstance)
    if err != nil {
        return false
    }
    StoredOTP:=UserInstance["OTP"]
    if StoredOTP == "" || SentOTP == ""{
        return false
    }
    if StoredOTP == SentOTP{
        _=rds.Redis_Delete_key(username)
        return true
    }else {
        return false
    }
}

