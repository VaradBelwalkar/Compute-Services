package mongodb

import (
	"context"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Authenticate user against DB entry
//Returns appropriate statusCodes
func Authenticate_user(username string,password string)(int){
	//CHANGE THIS LATER
	result:=resultStruct{}

	err := CollectionHandler.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err == mongo.ErrNoDocuments {	
		return 404
	} else if err != nil {
		return 500
	} else {
		chk:=compareHashAndPassword(result.Password,password)
		if chk==true{
			return 200
		}else{
			return 401
		}
		// If a document with the specified username already exists, update it
	}


}


