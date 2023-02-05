package database_handling

import (
	"context"
	"log"
	"golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
//Currently working on just retrieving data
//Global Object 

var CollectionHandler *mongo.Collection 
var Sys_CollectionHandler *mongo.Collection 


//Register system information here (e.g docker images available)
var sys_info = bson.M {"server":"private_cloud","docker_images":[]string{"ubuntu","nginx"}}


func InitiateMongoDB(m *mongo.Client) (*mongo.Collection,*mongo.Collection) {
//Create an empty Database first within MongoDB
// Create appropriate collection which will contain information about user
// The preferred keys in document ----> username, containerObj {containerName, port}

	db := m.Database("private_cloud")
	
	var coll *mongo.Collection
	var sys_coll *mongo.Collection

	// Check if the collection already exists (for users)
	names, err := db.ListCollectionNames(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	var sysFound, userFound bool = false, false
	for _, name := range names {

		if !userFound || !sysFound {
			if name == "user_details"{
				userFound = true
			}
			if name == "sys_details"{
				sysFound = false
			}
		continue
		}
		break
	}

	if !userFound {
		// Create the collection
		opts := options.CreateCollection().SetMaxDocuments(1000).SetCapped(true)
		err:= db.CreateCollection(context.TODO(), "user_details", opts)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !sysFound {
		// Create the collection
		opts := options.CreateCollection().SetMaxDocuments(1000).SetCapped(true)
		err:= db.CreateCollection(context.TODO(), "system_details", opts)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	//Handle adding system details here by creating new document int the system_details
	


	return coll,sys_coll
}


//Authenticate user against DB entry
//Returns appropriate statusCodes
func Authenticate_user(username string,password string)(int){
	//CHANGE THIS LATER
	var result struct{
		username string
		password string
	}

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
		if string(hashedPassword[:]) == result.password{
			return 200
		}

	}
return 404

}


