package database_handling

import (
	"context"
	"io"
	"os"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"    
)
//Currently working on just retrieving data

func InitiateMongoDB(m session) {
//Create an empty Database first within MongoDB
// Create appropriate collection which will contain information about user
// The preferred keys in document ----> username, containerObj {containerName, port}

 
}


// To be called when authorized user requests storage of data
//Depends on the structure of the documents
func upload_file(){
//Handle POST request to get the bytes for the file 

}

func upload_folder(){
//Handle POST request to get teh bytes for the file


}












