package query_handling

import (

	"fmt"
	"net/http"
	"src/container_apis"
	"main"

)
// Here the user will be authenticated first and then request will be fulfilled


// HandlerFunc to be registered
func Container_Run(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	//Get the requested Image from from the request-URL and pass it to the Container handler
	ContainerCreate(ctx,Cli,imageName)


  }
  
  func Container_Resume(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	ContainerStart(ctx,Cli)
  }
  
  func Container_List(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	 
  
  }

  func Container_Stop_or_Remove(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	 
  
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

