package auth_service

import (
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
)


//To be called in Logout handler
func DeleteSession(sessionID string) bool{
	check:=db.Redis_Delete_key(sessionID)
	if check!=true{
		return false
	}
	return true
}

