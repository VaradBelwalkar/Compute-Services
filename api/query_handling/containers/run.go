package containers

import (

	"context"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	ca "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/container_apis"
	auth "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/auth"
)
  // HandlerFunc to be registered
  func Container_Run(w http.ResponseWriter, r *http.Request) {
	check,username:=auth.Verify_Auth(w,r)
	if check!=true{
		return
	}
	//Extracting required string from the request Structure
	vars := mux.Vars(r)
	//Get the requested Image from from the request-URL and pass it to the Container handler
	privateKey,Port,err:=ca.ContainerCreate(context.TODO(),Cli,vars["image"],username)	
	if err!=200{
		if err ==500{
		w.WriteHeader(http.StatusInternalServerError)
		return
		}
		if err==403{
			w.WriteHeader(http.StatusForbidden) // Means no more than 5 containers are allowed
			return
		} 
		if err == 404{
			w.WriteHeader(http.StatusNotFound)
			return
		}
		return
	}

	resp:=map[string]string{"privatekey":privateKey,"port":Port}
	b, _ := json.Marshal(resp)
	w.Write(b)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//Send reponse in the body


}