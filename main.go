package main

import (
	"context"
	"log"
	"github.com/gorilla/mux"
	"net/http"
    db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
    qh "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/query_handling"
	as "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service"
	"github.com/docker/docker/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
    Cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
       panic(err)
    }
    defer Cli.Close()

	//Get handler for the "user_details" collection (creates collection if not exists)
    db.CollectionHandler,db.Sys_CollectionHandler=db.InitiateMongoDB(mongo_client);
    
    //login to be handled separatly

	router.HandleFunc("/", as.RenderForm)								//DONE			
	router.HandleFunc("/login", as.LoginHandler).Methods("POST")		//DONE
	router.HandleFunc("/login", as.RenderForm).Methods("GET")			//DONE
	router.HandleFunc("/logout", as.LogoutHandler).Methods("POST")		//DONE
	router.HandleFunc("/container/run/*", qh.Container_Run)
	router.HandleFunc("/container/resume/*", qh.Container_Resume)
	router.HandleFunc("/container/stop/*", qh.Container_Stop)
	router.HandleFunc("/container/remove/*", qh.Container_Remove)
	router.HandleFunc("/container/list/*", qh.Container_List)
	router.HandleFunc("/upload_file/", qh.Upload_file)			//yet to be determined
	router.HandleFunc("/upload_folder/", qh.Upload_folder)		//yet to be determined
    router.HandleFunc("/register",as.RenderForm).Methods("GET")	   	//DONE
	router.HandleFunc("/register",qh.RegisterUser).Methods("POST")     //DONE
	router.HandleFunc("/remove_account",as.RenderForm).Methods("GET")	//DONE
	router.HandleFunc("/remove_account",qh.RemoveAccount).Methods("POST")  //DONE

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}