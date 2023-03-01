package container_apis

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"
)


//Stop the container 
func ContainerStop(ctx context.Context,cli *client.Client,containerName string,username string) int{
	//** Check if containerName is valid or ContainerStart requires id to start the container
	//Handle db call to retrieve the 'id' for the container required to start the container
	//var documentData map[string]interface{} 
	//Make db call to retrieve user info about the containers it holds
	//var err error
	documentData, status := get_document(context.TODO(),username)
	if status!= 200{
		return status
	}
	var containerID string

	nesting1:=documentData.ContainerInfo//[containerName].(map[string]interface{})["containerID"].(string)
	if nesting2,ok:= nesting1[containerName]; ok{
		containerID=nesting2.(map[string]interface{})["containerID"].(string)
	} else{
		return 404 // Send StatusNotFound
	}

	stopOpts:=container.StopOptions{}
    if err := cli.ContainerStop(context.TODO(),containerID,stopOpts); err != nil {
        return 500
    }

	filter:=bson.M{
		"username":username,
	}
	update:=bson.M{ "$set":bson.M{
		"containerInfo."+containerName+"."+"status":"stopped",
	}}

	updateResult,err:=db.CollectionHandler.UpdateOne(ctx,filter,update)
	if err!=nil || updateResult.MatchedCount!=1{
		return 500
	}	
	
	return 200

}
















