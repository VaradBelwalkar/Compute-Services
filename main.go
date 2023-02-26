package main

import (
	"context"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"github.com/joho/godotenv"
    db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
    qh "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/query_handling"
	as "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service"
	"github.com/docker/docker/client"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gorilla/handlers"

)


//const url = "mongodb://host1:27017,host2:27017,host3:27017/?replicaSet=myRS"

var mongoURL string

func Setup_Env(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURL = os.Getenv("MONGODB_URI")
	db.Redis_URL = os.Getenv("REDIS_URL")
    db.Redis_Password = os.Getenv("REDIS_PASSWORD")
	db.PassHashKey = os.Getenv("PASSWORD_HASH_SECRET")
	as.JWTSigningKey = os.Getenv("JWT_SECRET")
}


// The main function manages all the query handling and manages the database as well
func main() {
    
	Setup_Env()
    // server main method

    var router = mux.NewRouter()

    //Initiate Mongo client
	mongo_client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL)) 
	if err != nil {
		log.Fatal("Please start MongoDB")
	}
	ctx := context.Background()
	err = mongo_client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer mongo_client.Disconnect(ctx)



    // Initiate Docker client
    qh.Cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
       panic("Failed to create Docker Client, ensure that Docker Daemon is running\n")
    }
    defer qh.Cli.Close()

	//Get handler for the "user_details" collection (creates collection if not exists)
    db.CollectionHandler,db.Sys_CollectionHandler=db.InitiateMongoDB(mongo_client);
	db.Initiate_Redis()
    
    //login to be handled separatly

	router.HandleFunc("/", as.RenderForm)								//DONE			
	router.HandleFunc("/login", as.LoginHandler).Methods("POST")		//DONE
	router.HandleFunc("/login", as.RenderForm).Methods("GET")			//DONE
	router.HandleFunc("/logout", as.LogoutHandler).Methods("POST")		//DONE
	router.HandleFunc("/otphandler", as.RenderForm).Methods("GET")
	router.HandleFunc("/otphandler", as.OTPHandler).Methods("POST")
	router.HandleFunc("/container/run/{image}", qh.Container_Run).Methods("GET")
	router.HandleFunc("/container/run/{image}", qh.Container_Run).Methods("POST")
	router.HandleFunc("/container/resume/{container}", qh.Container_Resume).Methods("GET")
	router.HandleFunc("/container/stop/{container}", qh.Container_Stop).Methods("GET")
	router.HandleFunc("/container/remove/{container}", qh.Container_Remove).Methods("GET")
	router.HandleFunc("/container/list/containers", qh.Container_List).Methods("GET")
    router.HandleFunc("/register",as.RenderForm).Methods("GET")	   	//DONE
	router.HandleFunc("/register",qh.RegisterUser).Methods("POST")     //DONE
	router.HandleFunc("/regotphandler",qh.VerifyRegisterUser).Methods("POST") 
	router.HandleFunc("/remove_account",as.RenderForm).Methods("GET")	//DONE
	router.HandleFunc("/remove_account",qh.RemoveAccount).Methods("POST")  //DONE

	http.Handle("/", router)
	http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}