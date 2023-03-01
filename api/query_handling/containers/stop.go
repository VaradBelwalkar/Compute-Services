package containers

import (

	"context"
	"net/http"
	"github.com/gorilla/mux"
	ca "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/container_apis"
	auth "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service/auth"
)

func Container_Stop(w http.ResponseWriter, r *http.Request) {
	check,username:=auth.Verify_Auth(w,r)
	if check!=true{
		return
	}
 	ctx := context.Background()
	//Extracting required string from the request Structure
	vars := mux.Vars(r)
	//Get the requested Image from from the request-URL and pass it to the Container handler
	err:=ca.ContainerStop(ctx,Cli,vars["container"],username)	
	if err!=200{
		if err ==500{
		w.WriteHeader(http.StatusInternalServerError)
		return
		}
		if err== 404{
			w.WriteHeader(http.StatusNotFound)
			return
		}
	return
	}

	w.WriteHeader(http.StatusOK)

	//Send reponse in the body


}