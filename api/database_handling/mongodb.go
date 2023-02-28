package database_handling

import (
	"context"
	"log"
	"fmt"
	"crypto/hmac"
    "crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"crypto/subtle"
    "encoding/hex"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var PassHashKey string

//Global Objects 
type resultStruct struct{
	Username string `bson:"username"`
	Password string `bson:"password"`					// In the hash format
    Email string `bson:"email"`
	ContainerInfo map[string]interface{} `bson:"containerInfo"`
	TotalOwnedContainers int `bson:"totalOwnedContainers,omitempty"`
}

var CollectionHandler *mongo.Collection 
var Sys_CollectionHandler *mongo.Collection 
var DatabaseHandler *mongo.Database

//Register system information here (e.g docker images available)
var sys_info = bson.M {"server":"private_cloud","docker_images":[]string{"ubuntu","nginx"}}

type DBConfig struct {
	Name string     	 `json:"name"`
	Collections []string `json:"collections"`
}

type Config struct {
	Emails []string `json:"emails"`
	Database DBConfig `json:"database"`
}

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
		
		err:= DatabaseHandler.CreateCollection(context.TODO(), "user_details", &opts)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !sysFound {
		// Create the collection
		err:= DatabaseHandler.CreateCollection(context.TODO(), "system_details", &opts)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	//Handle adding system details here by creating new document int the system_details
	
return DatabaseHandler.Collection("user_details"),DatabaseHandler.Collection("system_details")
	
}

func ComputeHash(password string) string {
    key := []byte(PassHashKey)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(password))
    hash := hex.EncodeToString(h.Sum(nil))
    return hash
}

func compareHashAndPassword(hash, password string) bool {
    expectedHash := ComputeHash(password)
    return subtle.ConstantTimeCompare([]byte(hash), []byte(expectedHash)) == 1
}

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


