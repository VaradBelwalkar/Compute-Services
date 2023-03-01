package auth
import(
	"net/http"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"
	twofa "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/twofa"
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

