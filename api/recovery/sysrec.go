package recovery

import (
	"context"
	"log"
    "runtime"
    "github.com/docker/docker/client"
	"go.mongodb.org/mongo-driver/bson"
    mndb "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"

)

type resultStruct struct{
	Username string `bson:"username"`
	Password string `bson:"password"`					// In the hash format
    Email string `bson:"email"`
	ContainerInfo map[string]interface{} `bson:"containerInfo"`
	TotalOwnedContainers int `bson:"totalOwnedContainers,omitempty"`
}


func UpdateContainerStatus() error {
    // Get a cursor for all the documents in the collection
    cursor, err := mndb.DatabaseHandler.Collection("user_details").Find(context.Background(), bson.M{})
    if err != nil {
        return err
    }
    defer cursor.Close(context.Background())

    // Create a channel for the user documents
    userChan := make(chan resultStruct)

    // Launch multiple goroutines to update the container status for each user document
    for i := 0; i < runtime.NumCPU(); i++ {
        go func() {
            ctx := context.Background()
            cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
            if err != nil {
                log.Fatal(err)
            }
            user:= <-userChan
            for containerName,ContainerOBJ := range user.ContainerInfo {
                // Check if the user has a container and get its ID
                containerID := ContainerOBJ.(map[string]interface{})["containerID"].(string)

                    containerStatus, err := cli.ContainerInspect(ctx, containerID)
                    if err != nil {
                        if client.IsErrNotFound(err) {
                            filter:=bson.M{
                                "username":user.Username,	
                            }
                            update:=bson.M{ "$unset":bson.M{
                                "containerInfo."+containerName:"",
                            },
                            "$set":bson.M{"totalOwnedContainers":user.TotalOwnedContainers-1},
                        
                        }
                            updateResult,check:=mndb.CollectionHandler.UpdateOne(ctx,filter,update)
                            if check!=nil || updateResult.MatchedCount!=1{
                                return
                            }
                        } else {
                            log.Fatal(err)
                        }
                    } else {
                        // Update the container status in the user document
                        status := "running"
                        if !containerStatus.State.Running {
                            status = "stopped"
                        }
                        filter:=bson.M{
                            "username":user.Username,
                        }
                        update:=bson.M{ "$set":bson.M{
                            "containerInfo."+containerName+"."+"status":status,
                        }}
                    
                        updateResult,err:=mndb.CollectionHandler.UpdateOne(ctx,filter,update)
                        if err!=nil || updateResult.MatchedCount!=1{
                            return 
                        }
                    }
                
            }
        }()
    }

    // Loop through all the documents and send them to the channel
    for cursor.Next(context.Background()) {
        var user resultStruct
        if err := cursor.Decode(&user); err != nil {
            return err
        }
        userChan <- user
    }

    // Close the channel to signal the goroutines to exit
    close(userChan)

    return cursor.Err()
}
