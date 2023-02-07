package database_handling

import (
	"context"
	"log"
	"fmt"
	"golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Global Objects 

var CollectionHandler *mongo.Collection 
var Sys_CollectionHandler *mongo.Collection 


//Register system information here (e.g docker images available)
var sys_info = bson.M {"server":"private_cloud","docker_images":[]string{"ubuntu","nginx"}}


func InitiateMongoDB(m *mongo.Client) (*mongo.Collection,*mongo.Collection) {
//Create an empty Database first within MongoDB
// Create appropriate collection which will contain information about user
// The preferred keys in document ----> username, containerObj {containerName, port}
var temp bool =true
var byteSize int64=120000
var maxDoc int64=120
opts:=options.CreateCollectionOptions{
	Capped:&temp,
	SizeInBytes:&byteSize,
	MaxDocuments:&maxDoc,
}

	db := m.Database("private_cloud")

	// Check if the collection already exists (for users)
	names, err := db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Something went wrong!")
	}
	var sysFound, userFound bool = false, false
	for _, name := range names {
		if !userFound || !sysFound {
			if name == "user_details"{
				userFound = true
			}
			if name == "system_details"{
				sysFound = true
			}
		continue
		}
		break
	}

	if !userFound {

		// Create the collection
		
		err:= db.CreateCollection(context.TODO(), "user_details", &opts)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !sysFound {
		// Create the collection
		err:= db.CreateCollection(context.TODO(), "system_details", &opts)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	//Handle adding system details here by creating new document int the system_details
	
return db.Collection("user_details"),db.Collection("system_details")
	
}


//Authenticate user against DB entry
//Returns appropriate statusCodes
func Authenticate_user(username string,password string)(int){
	//CHANGE THIS LATER
	type resultStruct struct{
		username string
		password []byte
	}

	result:=resultStruct{}

	err := CollectionHandler.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err == mongo.ErrNoDocuments {
	
		return 404
	} else if err != nil {
		return 500
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return 500
		}
		// If a document with the specified username already exists, update it
		if string(hashedPassword[:]) == string(result.password[:]){
			return 200
		}
		fmt.Println("again problem with hashing")
	}
return 404

}


