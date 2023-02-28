package recovery

import (
	"context"
	"log"
    "runtime"
    "github.com/docker/docker/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

)

type resultStruct struct{
	Username string `bson:"username"`
	Password string `bson:"password"`					// In the hash format
    Email string `bson:"email"`
	ContainerInfo map[string]interface{} `bson:"containerInfo"`
	TotalOwnedContainers int `bson:"totalOwnedContainers,omitempty"`
}


func (s resultStruct) updateStatus(){

}

func updateContainerStatus(db *mongo.Database) error {
    // Get a cursor for all the documents in the collection
    cursor, err := db.Collection("user_details").Find(context.Background(), bson.M{})
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
                containerID := ContainerOBJ.(map[string]string)["containerID"]

                    containerStatus, err := cli.ContainerInspect(ctx, containerID)
                    if err != nil {
                        if client.IsErrNotFound(err) {
                            update:=bson.M{ "$unset":bson.M{
                                "containerInfo."+containerName:"",
                            },
                            "$set":bson.M{"totalOwnedContainers":user.TotalOwnedContainers-1},
                        
                        }
                            if _, err := db.Collection("users").UpdateByID(ctx, user.Username, update); err != nil {
                                log.Fatal(err)
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
                        update:=bson.M{ "$set":bson.M{
                            "containerInfo."+containerName+"."+"status":status,
                        }}
                        if _, err := db.Collection("users").UpdateByID(ctx, user.Username, update); err != nil {
                            log.Fatal(err)
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
