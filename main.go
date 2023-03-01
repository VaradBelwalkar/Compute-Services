package main

import (
	"net/http"
	"os"
	"github.com/gorilla/handlers"

)



// The main function manages all the query handling and manages the database as well
func main() {
    
	Setup_Env()
    // server main method
    var router = Setup()

	http.Handle("/", router)
	http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}