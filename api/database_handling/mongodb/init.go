package mongodb

import (
	"context"
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func setupDB() (string,string,string){
	// read the config file
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Config file not found, Please provide Configuration file!")
		return "","",""
	}

	// unmarshal the JSON data into a Config struct
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println("Configuration file format invalid!")
		return "","",""
	}

	return config.Database.Name,config.Database.Collections[0],config.Database.Collections[1]

}


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

	dbName,userDetails,sysDetails:=setupDB()
	if dbName=="" || userDetails == "" || sysDetails == ""{
		return nil,nil
	}
	DatabaseHandler = m.Database(dbName)

	// Check if the collection already exists (for users)
	names, err := DatabaseHandler.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Something went wrong!")
	}
	var sysFound, userFound bool = false, false
	for _, name := range names {
		if !userFound || !sysFound {
			if name == userDetails{
				userFound = true
			}
			if name == sysDetails{
				sysFound = true
			}
		continue
		}
		break
	}

	if !userFound {

		// Create the collection
		userDetails="user_details"
		err:= DatabaseHandler.CreateCollection(context.TODO(), "user_details", &opts)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !sysFound {
		// Create the collection
		sysDetails="system_details"
		err:= DatabaseHandler.CreateCollection(context.TODO(), "system_details", &opts)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	//Handle adding system details here by creating new document int the system_details
	
return DatabaseHandler.Collection(userDetails),DatabaseHandler.Collection(sysDetails)
	
}



