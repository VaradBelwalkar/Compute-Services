package redis

import (
	"time"
	"context"
    //"os"
	"github.com/go-redis/redis/v8"
)
var(
    Redis_Client *redis.Client
    Redis_URL string
    Redis_Password string 
)
func Initiate_Redis() {
   // redisPassword := os.Getenv("redis_password")
	Redis_Client = redis.NewClient(&redis.Options{
		Addr: Redis_URL,					
		Password: Redis_Password,							
		DB: 0,
	})

}   


//returns string
func Redis_Get_Value(key string) string {
    // create context
    ctx := context.Background()

    // get value for key
    value, err := Redis_Client.Get(ctx, key).Result()
    if err != nil {
        return ""
    }
	return value
}


func Redis_Set_Value(key string,value string) bool{
    // create context
    ctx := context.Background()

    // set key with value and timeout
    err := Redis_Client.Set(ctx, key, value, 0).Err()
    if err != nil {
        panic(err)
    }

    // output set status
    if set, _ := Redis_Client.Exists(ctx, key).Result(); set == 1 {
        return true
    } else {
        return false
    }

}

// Use 5 minutes for OTP login 
func Redis_Set_Value_With_Timeout(key string,value string,timeout int) bool{
    // create context
    ctx := context.Background()

    // set key with value and timeout
    err := Redis_Client.Set(ctx, key, value, time.Duration(timeout)*time.Minute).Err()
    if err != nil {
        return false
    }

    // output set status
    if set, _ := Redis_Client.Exists(ctx, key).Result(); set == 1 {
        return true
    } else {
        return false
    }

}


func Redis_Delete_key(key string) bool{
	    // create context
		ctx := context.Background()

		// delete key
		err := Redis_Client.Del(ctx,key).Err()
		if err != nil {
			return false
		}
	
		// output deletion status
		if deleted, _ := Redis_Client.Exists(ctx,key).Result(); deleted == 0 {
			return true
		} else {
			return false
		}
	
}


func Redis_Check_Key(key string) bool{
	ctx := context.Background()
	if check, _ := Redis_Client.Exists(ctx,key).Result(); check == 0 {
		return false
	} else {
		return true
	}
}