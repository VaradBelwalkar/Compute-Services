package query_handling

import (

	"context"
	"net/http"
	"github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/container_apis"

	"path/filepath"
	"github.com/docker/docker/client"
)
// Here the user will be authenticated first and then request will be fulfilled

var Cli *client.Client

  // HandlerFunc to be registered
func Container_Run(w http.ResponseWriter, req *http.Request) {
 	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := req.URL.Path // Suppose "/foo/bar/something"
	imageName:=filepath.Base(URLPath) // will output "something"
	//Get the requested Image from from the request-URL and pass it to the Container handler
	container_apis.ContainerCreate(ctx,Cli,imageName)


}
  
func Container_Resume(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := req.URL.Path // Suppose "/foo/bar/something"
	containerName:=filepath.Base(URLPath) // will output "something"	
	container_apis.ContainerStart(ctx,Cli,containerName)
}
  
func Container_List(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	//Retrieve username appropriatly

	container_apis.OwnedContainerInfo(ctx,"username")
  
}

func Container_Stop_or_Remove(w http.ResponseWriter, _ *http.Request) {
	//_ := context.Background()
	//Extracting required string from the request Structure
    //URLPath := _.URL.Path // Suppose "/foo/bar/something"
	//imageName:=filepath.Base(URLPath) // will output "something" 
  
}

