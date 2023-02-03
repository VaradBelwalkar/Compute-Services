package main

import (
	"fmt"
	"context"
	"io"
	"os"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
    "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
    "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/query_handling"
    "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	
	"github.com/docker/docker/pkg/stdcopy"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)


const url = "mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=myRS"



// The main function manages all the query handling and manages the database as well
func main() {

    // server main method

    var router = mux.NewRouter()

    //Initiate Mongo client
	mongo_client, err := mongo.NewClient(options.Client().ApplyURI(url)) 
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	err = mongo_client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mongo_client.Disconnect(ctx)



    // Initiate Docker client
    query_handling.Cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
       panic(err)
    }
    defer Cli.Close()

	//Get handler for the "user_details" collection (creates collection if not exists)
    database_handling.CollectionHandler,database_handling.Sys_CollectionHandler=InitiateMongoDB(mongo_client);
    
    //login to be handled separatly

	router.HandleFunc("/", RenderForm)
	router.HandleFunc("/internal", internalPageHandler)
	router.HandleFunc("/login", HandleSubmit).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/container/run/*", Container_Run)
	router.HandleFunc("/container/resume/*", Container_Resume)
	router.HandleFunc("/container/<regex>/*", Container_Stop_or_Remove)
	router.HandleFunc("/container/list/*", Container_List)
	router.HandleFunc("/upload_file/", upload_file)
	router.HandleFunc("/upload_folder/", upload_folder)
    router.HandleFunc("/register",RegisterUser)
	router.HandleFunc("/remove_account",RemoveAccount)

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}