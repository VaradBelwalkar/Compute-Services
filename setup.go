package main

import (
	"context"
	"log"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	mndb "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"
	"github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/redis"
	jwt "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/jwt"
    "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/twofa"
	containers "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/query_handling/containers"
	"github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/container_apis"
	"github.com/VaradBelwalkar/Private-Cloud-MongoDB/routes"
	"github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/recovery"
	"github.com/docker/docker/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)


//const url = "mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=myRS"

var mongoURL string
type dbConfig struct {
	Name string     	 `json:"name"`
	Collections []string `json:"collections"`
}

type config struct {
	Emails []string `json:"emails"`
	Database dbConfig `json:"database"`
	Port string `json:"server_port"`
	Images []string `json:"images"`
}


func Setup_Env(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURL = os.Getenv("MONGODB_URI")
	redis.Redis_URL = os.Getenv("REDIS_URL")
    redis.Redis_Password = os.Getenv("REDIS_PASSWORD")
	mndb.PassHashKey = os.Getenv("PASSWORD_HASH_SECRET")
	jwt.JWTSigningKey = os.Getenv("JWT_SECRET")
	twofa.Official_Email = os.Getenv("OFFICIAL_EMAIL")
	twofa.Official_Email_Password =  os.Getenv("OFFICIAL_EMAIL_APP_PASSWORD")
}


// The main function manages all the query handling and manages the database as well
func Setup() (*mux.Router,string){
    
	Setup_Env()
    //Initiate Mongo client
	mongo_client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL)) 
	if err != nil {
		log.Fatal("Please start MongoDB")
	}
	ctx := context.Background()
	err = mongo_client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

    // Initiate Docker client
    containers.Cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
       panic("Failed to create Docker Client, ensure that Docker Daemon is running\n")
    }

	//Get handler for the "user_details" collection (creates collection if not exists)
    mndb.CollectionHandler,mndb.Sys_CollectionHandler=mndb.InitiateMongoDB(mongo_client);
	redis.Initiate_Redis()
    recovery.UpdateContainerStatus()
    //login to be handled separatly
	router:=routes.NewRouter()

	//Get server port from configuration
		// read the config file
		data, err := ioutil.ReadFile("config.json")
		if err != nil {
			fmt.Println("Config file not found, Please provide Configuration file!")
			return nil,""
		}
	
		// unmarshal the JSON data into a Config struct
		var config config
		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Println("Configuration file format invalid!")
			return nil,""
		}

		container_apis.ImageArray = config.Images


	return router,config.Port
}