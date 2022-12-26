package query_handling

import (

	"fmt"
	"net/http"
	"github.com/VaradBelwalkar/Private-Cloud/src/container_apis"
	"github.com/VaradBelwalkar/Private-Cloud"
	"path/filepath"

)
// Here the user will be authenticated first and then request will be fulfilled


// HandlerFunc to be registered
func Container_Run(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := _.URL.Path // Suppose "/foo/bar/something"
	imageName:=filepath.Base(URLPath) // will output "something"
	//Get the requested Image from from the request-URL and pass it to the Container handler
	ContainerCreate(ctx,Cli,imageName)


  }
  
  func Container_Resume(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := _.URL.Path // Suppose "/foo/bar/something"
	containerName:=filepath.Base(URLPath) // will output "something"	
	ContainerStart(ctx,Cli,containerName)
  }
  
  func Container_List(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	OwnedContainerInfo(ctx,)
  
  }

  func Container_Stop_or_Remove(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := _.URL.Path // Suppose "/foo/bar/something"
	imageName:=filepath.Base(URLPath) // will output "something" 
  
  }

  func Container_List(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
  
  }

  func upload_file(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()


  }

  func upload_folder(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()

  }

func register_user(){


}


func remove_user(){



}