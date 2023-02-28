package auth_service

import (
	"context"
	"net/smtp"
    "time"
    "encoding/json"
	"strconv"
    "math/rand"
    "go.mongodb.org/mongo-driver/bson"
    db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
)
type resultStruct struct{
	Username string `bson:"username"`
	Password string `bson:"password"`
    Email string `bson:"email"`
	ContainerInfo map[string]interface{} `bson:"containerInfo"`
	TotalOwnedContainers int `bson:"totalOwnedContainers,omitempty"`
}

func Generate_OTP() string{
    rand.Seed(time.Now().UnixNano())
    min := 100000
    max := 999999
    otp := rand.Intn(max-min+1) + min
    return strconv.Itoa(otp)
}


func RetrieveEmail(username string) string{
	documentData:=resultStruct{}
    ctx:=context.Background()
	//Check user-document exists in the collection 
	//document_handler of type *SingleResult, see github code for more details
	err := db.CollectionHandler.FindOne(ctx,bson.M{"username":username}).Decode(&documentData)
	//If not then use following	
	if err != nil {		
			return ""	// Internal server error
	}

    return documentData.Email

}




func TwoFA_Send(username string,eMail string)(bool,string) {
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
		"Subject: Your OTP for Two-Factor Authentication\n\n\r\n" +
		"\r\n" +
		"Dear User,\n\nYour OTP for two-factor authentication is: " + OTP + "\n\nPlease enter this otp in your app to complete the authentication process.\n\nBest regards,\nDYPLUG\r\n")

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


func TwoFA_Verify(username string, SentOTP string) bool{
    UserInstance:=make(map[string]string)
    jsonString:=db.Redis_Get_Value(username)
    err := json.Unmarshal([]byte(jsonString), &UserInstance)
    if err != nil {
        return false
    }
    StoredOTP:=UserInstance["OTP"]
    if StoredOTP == "" || SentOTP == ""{
        return false
    }
    if StoredOTP == SentOTP{
        _=db.Redis_Delete_key(username)
        return true
    }else {
        return false
    }
}

