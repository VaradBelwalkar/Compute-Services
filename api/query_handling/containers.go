package query_handling

import (

	"context"
	"net/http"
	"encoding/json"
	ca "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/container_apis"
	as "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service"
	"path/filepath"
	"github.com/docker/docker/client"
)
// Here the user will be authenticated first and then request will be fulfilled

var Cli *client.Client
type responseStruct struct{
	privatekey string
	port string
}

type infoStruct struct{
	containerInfo []string
}


  // HandlerFunc to be registered
func Container_Run(w http.ResponseWriter, r *http.Request) {
	check,username:=as.Handle_auth(w,r)
	if check!=true{
		return
	}

 	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := r.URL.Path // Suppose "/foo/bar/something"
	imageName:=filepath.Base(URLPath) // will output "something"
	//Get the requested Image from from the request-URL and pass it to the Container handler
	privateKey,Port,err:=ca.ContainerCreate(ctx,Cli,imageName,username)	
	if err!=200{
		if err ==500{
		w.WriteHeader(http.StatusInternalServerError)
		return
		}
	}

	resp:=responseStruct{privatekey:privateKey,port:Port}
	json.NewEncoder(w).Encode(resp)
	w.Header().Set("Content-Type", "application/json")

	//Send reponse in the body


}
  
func Container_Resume(w http.ResponseWriter, r *http.Request) {
	check,username:=as.Handle_auth(w,r)
	if check!=true{
		return
	}

 	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := r.URL.Path // Suppose "/foo/bar/something"
	containerName:=filepath.Base(URLPath) // will output "something"
	//Get the requested Image from from the request-URL and pass it to the Container handler
	privateKey,Port,err:=ca.ContainerStart(ctx,Cli,containerName,username)	
	if err!=200{
		if err ==500{
		w.WriteHeader(http.StatusInternalServerError)
		return
		}
	}

	resp:=responseStruct{privatekey:privateKey,port:Port}
	json.NewEncoder(w).Encode(resp)
	w.Header().Set("Content-Type", "application/json")

	//Send reponse in the body


}
  
func Container_List(w http.ResponseWriter, r *http.Request) {
	check,username:=as.Handle_auth(w,r)
	if check!=true{
		return
	}

 	ctx := context.Background()
	//Extracting required string from the request Structure
	//Get the requested Image from from the request-URL and pass it to the Container handler
	containerArray,err:=ca.OwnedContainerInfo(ctx,username)	
	if err!=200{
		if err ==500{
		w.WriteHeader(http.StatusInternalServerError)
		return
		}
	}

	resp:=infoStruct{containerInfo:containerArray}
	json.NewEncoder(w).Encode(resp)
	w.Header().Set("Content-Type", "application/json")

	//Send reponse in the body


}

func Container_Stop(w http.ResponseWriter, r *http.Request) {
	check,username:=as.Handle_auth(w,r)
	if check!=true{
		return
	}

 	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := r.URL.Path // Suppose "/foo/bar/something"
	containerName:=filepath.Base(URLPath) // will output "something"
	//Get the requested Image from from the request-URL and pass it to the Container handler
	err:=ca.ContainerStop(ctx,Cli,containerName,username)	
	if err!=200{
		if err ==500{
		w.WriteHeader(http.StatusInternalServerError)
		return
		}
		if err== 404{
			w.WriteHeader(http.StatusNotFound)
		}
	}

	w.WriteHeader(http.StatusOK)

	//Send reponse in the body


}


func Container_Remove(w http.ResponseWriter, r *http.Request) {
	check,username:=as.Handle_auth(w,r)
	if check!=true{
		return
	}

 	ctx := context.Background()
	//Extracting required string from the request Structure
    URLPath := r.URL.Path // Suppose "/foo/bar/something"
	containerName:=filepath.Base(URLPath) // will output "something"
	//Get the requested Image from from the request-URL and pass it to the Container handler
	err:=ca.ContainerRemove(ctx,Cli,containerName,username)	
	if err!=200{
		if err ==500{
		w.WriteHeader(http.StatusInternalServerError)
		return
		}
		if err == 404{
			w.WriteHeader(http.StatusNotFound)
		}
	}
	
	w.WriteHeader(http.StatusOK)

	//Send reponse in the body


}

