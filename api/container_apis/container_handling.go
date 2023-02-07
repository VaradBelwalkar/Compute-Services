package container_apis

import (
	"context"
	"strings"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"golang.org/x/crypto/ssh"
	"encoding/pem"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
)
//Creating private-public key pairs to be used by the client to ssh into registered container
func MakeSSHKeyPair() (string, string, int) {
    privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
    if err != nil {
        return "", "", 500
    }

    // generate and write private key as PEM
    var privKeyBuf strings.Builder

    privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
    if err := pem.Encode(&privKeyBuf, privateKeyPEM); err != nil {
        return "", "", 500
    }

    // generate and write public key
    pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
    if err != nil {
        return "", "", 500
    }

    var pubKeyBuf strings.Builder
    pubKeyBuf.Write(ssh.MarshalAuthorizedKey(pub))

    return pubKeyBuf.String(), privKeyBuf.String(), 200
}


//Function to return err if document not found
func get_document(ctx context.Context,username string)(map[string]interface{}, int){

	var documentData map[string]interface{} 
	//Check user-document exists in the collection 
	//document_handler of type *SingleResult, see github code for more details
	err := db.CollectionHandler.FindOne(ctx,bson.M{"username":username}).Decode(&documentData)
	//If not then use following	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil,404
		} else {
			return nil,500
		}
	}	

return documentData,200

}


//Create a new container
func ContainerCreate(ctx context.Context,cli *client.Client,imageName string,username string) (string,string,int){

	documentData, status := get_document(context.TODO(),username)
	if status!= 200{
		return "","",status
	}
	//Here we get the document to work with

	
	totalOwnedContainers := documentData["totalOwnedContainers"].(int)
	
		if totalOwnedContainers>=5 {
			return "","",406		// http.StatusNotAcceptable
		}
	
	//Else do allocate the container	

	containerConfig:=container.Config{Image:imageName}
	portBinding:=nat.PortBinding{HostIP:"0.0.0.0"}	
	portBindings:=make(nat.PortMap)
	portBindings["22/tcp"][0]=portBinding		// A PortMap
	hostConfig:=container.HostConfig{PortBindings:portBindings}

    resp, err := cli.ContainerCreate(ctx,&containerConfig,&hostConfig,nil,nil,"")
    if err != nil {
        return "","",500
    }

	privateKey,publicKey,check:= MakeSSHKeyPair()
	if check!=200 {
		return "","",check
	}

	//First make a tar archive for the public key generated above 
	buf := strings.NewReader(publicKey)
	err =cli.CopyToContainer(context.Background(), resp.ID, "/home/user/.ssh/",buf ,types.CopyToContainerOptions{})
	if err!=nil{
		return "","",500
	}
	// Handle db call to store the resp.ID into the appropriate row for the user
	containerJSON,err:=cli.ContainerInspect(ctx,resp.ID)
	if err!=nil{
	return "","",500
	}
	port:=containerJSON.NetworkSettings.NetworkSettingsBase.Ports["22/tcp"][0].HostPort
	containerName:=username+"_"+port
	// Here count is updated but not container information, hence do update that
	filter:=bson.M{
		"username":username,
	}
	update:=bson.M{
		"totalOwnedContainers":totalOwnedContainers+1,
		"containerInfo":bson.M{containerName:bson.M{"containerID":resp.ID,"port":port,"status":"running"}},
	}

	updateResult,err:=db.CollectionHandler.UpdateOne(ctx,filter,update)
	if err!=nil || updateResult.MatchedCount!=1{
		return "","",500
	}

return privateKey,port,200
}



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

	nesting1:=documentData["containerInfo"].(map[string]interface{})//[containerName].(map[string]interface{})["containerID"].(string)
	if nesting2,ok:= nesting1[containerName]; ok{
		containerID=nesting2.(map[string]interface{})["containerID"].(string)
	} else{
		return 404 // Send StatusNotFound
	}

	var CHANGE time.Duration = 1029
    if err := cli.ContainerStop(ctx,containerID,&CHANGE); err != nil {
        return 500
    }

    statusCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
    select {
    case err := <-errCh:
        if err != nil {
			return 500
        }
    case <-statusCh:
    }

	//DO APPROPRIATE CHANGES HERE

	filter:=bson.M{
		"username":username,
	}
	update:=bson.M{
		"containerInfo":bson.M{containerName:bson.M{"status":"stopped"}},
	}

	updateResult,err:=db.CollectionHandler.UpdateOne(ctx,filter,update)
	if err!=nil || updateResult.MatchedCount!=1{
		return 500
	}	
	
	return 200

}


//Stop the container 
func ContainerRemove(ctx context.Context,cli *client.Client,containerName string,username string) int{
	//** Check if containerName is valid or ContainerStart requires id to start the container
	//Handle db call to retrieve the 'id' for the container required to start the container
	//var documentData map[string]interface{} 
	//Make db call to retrieve user info about the containers it holds
	//var err error
	documentData,err := get_document(ctx,username)
	if err!=200{
		return err
	}

	var containerID string

	nesting1:=documentData["containerInfo"].(map[string]interface{})//[containerName].(map[string]interface{})["containerID"].(string)
	if nesting2,ok:= nesting1[containerName]; ok{
		containerID=nesting2.(map[string]interface{})["containerID"].(string)
	} else{
		return 404 // Send StatusNotFound
	}

	var options types.ContainerRemoveOptions

    if err := cli.ContainerRemove(ctx, containerID,options); err != nil {
        return 500
    }

    statusCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
    select {
    case err := <-errCh:
        if err != nil {
            return 500
        }
    case <-statusCh:
    }

	filter:=bson.M{
		"username":username,
	}
	update:=bson.M{
		"containerInfo":bson.M{"$unset":bson.M{containerName:"",}},
	}
	updateResult,check:=db.CollectionHandler.UpdateOne(ctx,filter,update)
	if check!=nil || updateResult.MatchedCount!=1{
		return 500
	}	
	
	return 200

}







//Start the container if already created
func ContainerStart(ctx context.Context,cli *client.Client,containerName string,username string) (string,string,int){
	//** Check if containerName is valid or ContainerStart requires id to start the container
	//Handle db call to retrieve the 'id' for the container required to start the container

	documentData,status := get_document(ctx,username)
	if status!=200{
		return "","",status //Here the system failure has occured
	}


	var containerID string

	nesting1:=documentData["containerInfo"].(map[string]interface{})//[containerName].(map[string]interface{})["containerID"].(string)
	if nesting2,ok:= nesting1[containerName]; ok{
		containerID=nesting2.(map[string]interface{})["containerID"].(string)
	} else{
		return "","",404 // Send StatusNotFound
	}
	// If returns error if container is already running, first do inspect the container and then only run
    if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
        return "","",500
    }

    statusCh, errCh := cli.ContainerWait(ctx, "id", container.WaitConditionNotRunning)
    select {
    case err := <-errCh:
        if err != nil {
          return "","",500
        }
    case <-statusCh:
    }

	privateKey,publicKey,err:= MakeSSHKeyPair()
	if err!=200 {
		return "","",err
	}

	//First make a tar archive for the public key generated above 
	buf := strings.NewReader(publicKey)
	check :=cli.CopyToContainer(context.Background(), containerID, "/home/user/.ssh/", buf,types.CopyToContainerOptions{})
	if check!=nil{
		return "","",500
	}

	containerJSON,check:=cli.ContainerInspect(ctx,containerID)
	if check!=nil{
	return "","",500
	}
	port:=containerJSON.NetworkSettings.NetworkSettingsBase.Ports["22/tcp"][0].HostPort
	newContainerName:=username+"_"+port
	

	filter:=bson.M{
		"username":username,
	}
	update:=bson.M{
		"containerInfo":bson.M{"$unset":bson.M{containerName:"",},"$set":bson.M{newContainerName:bson.M{"containerID":containerID,"port":port,"status":"running"},}},
		
	}
	updateResult,check:=db.CollectionHandler.UpdateOne(ctx,filter,update)
	if check!=nil || updateResult.MatchedCount!=1{
		return "","",500
	}

	return privateKey,port,200

}

//Gives information about the containers that user holds
func OwnedContainerInfo(ctx context.Context,username string)([]string,int){

	documentData,err := get_document(ctx,username)
	if err!=200{
		return nil,err
	}
	var containerArray []string
	for k, _ := range documentData["containerInfo"].(map[string]interface{}) { 
		containerArray= append(containerArray,k)
		
	}
	
	return containerArray,200

}



//Gives images available on the server (the information is available on "system_details" collection)
func ImageInfo(ctx context.Context,cli *client.Client) []string{
	
	var imageArray=[]string{"ubuntu","development_server"}

	return imageArray

}











