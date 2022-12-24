package main

import (
	"context"
	"io"
	"os"
    "src/database_handling"
    "src/query_handling"
    "net/http"
    "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// The main function manages all the query handling and manages the database as well

const url = "mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=myRS"

func main() {
	mongo_client, err := mongo.NewClient(options.Client().ApplyURI(url)) // Give appropriate port here
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



Cli_native, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
if err != nil {
    panic(err)
}
defer cli.Close()


InitiateMongoDB(mongo_client);

// Get userinfo-collection handler 
collection := mongo_client.Database("private-cloud").Collection("userinfo")



  // create a mux to hold url and handlers
  mux := http.NewServeMux()
  log.Println("created mux")

  // register handlers using HandleFunc
  mux.HandleFunc("/container/run/*", Container_Run(ResponseWriter,*Request,collection))
  mux.HandleFunc("/container/resume/*", Container_Resume(ResponseWriter,*Request,collection))
  mux.HandleFunc("/container/<regex>/*", Container_Stop_or_Remove(ResponseWriter,*Request,collection))
  mux.HandleFunc("/container/list/*", Container_List(ResponseWriter,*Request,collection))
  //Handle POST Request
  mux.HandleFunc("/upload_file/", upload_file(ResponseWriter,*Request,collection))
  mux.HandleFunc("/upload_folder/", upload_folder(ResponseWriter,*Request,collection))
  //Here we have registered the handlers

  // register mux with server and listen for requests
  log.Println("starting server")
  log.Fatal(http.ListenAndServe(":8080", mux))
}
