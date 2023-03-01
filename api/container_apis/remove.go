package container_apis

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"
)


//Stop the container 
func ContainerRemove(ctx context.Context,cli *client.Client,containerName string,username string) int{
	//** Check if containerName is valid or ContainerStart requires id to start the container
	//Handle db call to retrieve the 'id' for the container required to start the container
	//var documentData map[string]interface{} 
	//Make db call to retrieve user info about the containers it holds
	//var err error
	documentData,err := get_document(context.TODO(),username)
	if err!=200{
		return err
	}
	totalOwnedContainers:=documentData.TotalOwnedContainers
	var containerID string

	nesting1:=documentData.ContainerInfo//[containerName].(map[string]interface{})["containerID"].(string)
	if nesting2,ok:= nesting1[containerName]; ok{
		containerID=nesting2.(map[string]interface{})["containerID"].(string)
	} else{
		return 404 // Send StatusNotFound
	}

	options:=types.ContainerRemoveOptions{Force:true}

    if err := cli.ContainerRemove(context.TODO(), containerID,options); err != nil {
        return 500
    }

	filter:=bson.M{
		"username":username,	
	}
	update:=bson.M{ "$unset":bson.M{
		"containerInfo."+containerName:"",
	},
	"$set":bson.M{"totalOwnedContainers":totalOwnedContainers-1},

}
	updateResult,check:=db.CollectionHandler.UpdateOne(ctx,filter,update)
	if check!=nil || updateResult.MatchedCount!=1{
		return 500
	}		
	return 200

}
