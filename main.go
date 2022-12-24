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
	"github.com/gocql/gocql"
)

// The main function manages all the query handling and manages the database as well


func main() {

 //Initiating Cassandra 
// connect to the cluster
cluster := gocql.NewCluster("PublicIP", "PublicIP", "PublicIP") //replace PublicIP with the IP addresses used by your cluster.
cluster.Consistency = gocql.Quorum
cluster.ProtoVersion = 4
cluster.ConnectTimeout = time.Second * 10
cluster.Authenticator = gocql.PasswordAuthenticator{Username: "Username", Password: "Password"} //replace the username and password fields with their real settings.
session, err := cluster.CreateSession()
    if err != nil {
    log.Println(err)
    return
}
defer session.Close() 

// Initiate Docker client



Cli_native, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
if err != nil {
    panic(err)
}
defer cli.Close()


InitiateCassandra(session);

  // create a mux to hold url and handlers
  mux := http.NewServeMux()
  log.Println("created mux")

  // register handlers using HandleFunc
  mux.HandleFunc("/container/run/*", Container_Run)
  mux.HandleFunc("/container/resume/*", Container_Resume)
  mux.HandleFunc("/container/<regex>/*", Container_Stop_or_Remove)
  mux.HandleFunc("/container/list/*", Container_List)
  //Handle POST Request
  mux.HandleFunc("/upload_file/", upload_file)
  mux.HandleFunc("/upload_folder/", upload_folder)
  //Here we have registered the handlers

  // register mux with server and listen for requests
  log.Println("starting server")
  log.Fatal(http.ListenAndServe(":8080", mux))
}

