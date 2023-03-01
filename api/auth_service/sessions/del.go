package auth_service

import (
	rds "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/redis"
)


//To be called in Logout handler
func DeleteSession(sessionID string) bool{
	check:=rds.Redis_Delete_key(sessionID)
	if check!=true{
		return false
	}
	return true
}

