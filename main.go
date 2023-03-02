package main

import (
	"net/http"
	"os"
	"github.com/gorilla/handlers"

)

func main() {
    // server main method
    var router,server_port = Setup()

	http.Handle("/", router)
	http.ListenAndServe(":"+server_port, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}	