package auth_service

import (
	"math/rand"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
)


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func generateSessionID(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
	ok := db.Redis_Check_Key(string(b))
	for ok==true{
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		ok = db.Redis_Check_Key(string(b))
	}
    return string(b)
}