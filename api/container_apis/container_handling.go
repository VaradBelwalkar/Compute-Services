package container_apis

import (
	"context"
	"strings"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"golang.org/x/crypto/ssh"
	"encoding/pem"
	"github.com/VaradBelwalkar/Private-Cloud/api/database_handling"
)
//Creating private-public key pairs to be used by the client to ssh into registered container
func MakeSSHKeyPair() (string, string, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
    if err != nil {
        return "", "", err
    }

    // generate and write private key as PEM
    var privKeyBuf strings.Builder

    privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
    if err := pem.Encode(&privKeyBuf, privateKeyPEM); err != nil {
        return "", "", err
    }

    // generate and write public key
    pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
    if err != nil {
        return "", "", err
    }

    var pubKeyBuf strings.Builder
    pubKeyBuf.Write(ssh.MarshalAuthorizedKey(pub))

    return pubKeyBuf.String(), privKeyBuf.String(), nil
}


//Function to return err if document not found
func get_document(ctx context.Context,username string)(map[string]interface{}, error){

	var documentData map[string]interface{} 
	//Check user-document exists in the collection 
	//document_handler of type *SingleResult, see github code for more details
	err := database_handling.CollectionHandler.FindOne(ctx,bson.M{"username":username}).Decode(&documentData)
	//If not then use following	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil,err
		} else {
			return nil,err
		}
	}	

return documentData,nil

}


//Create a new container
func ContainerCreate(ctx context.Context,cli *client.Client,imageName string) (string,error){

	var documentData map[string]interface{} 
	var err error
	documentData, err = get_document(context.TODO(),"CHANGE_USER_NAME_HERE")
	if err!= nil{
		return "",err
	}
	//Here we get the document to work with

	
	totalOwnedContainers := documentData["totalOwnedContainers"].(int)
	
		if totalOwnedContainers>5 {
			//Handle response
		}
	
	//Else do allocate the container	

    resp, err := cli.ContainerCreate(ctx, &container.Config{ 
        Image: imageName,
    }, nil, nil, nil, "")
    if err != nil {
        return "",err
    }

	privateKey,publicKey,err:= MakeSSHKeyPair()
	if err!=nil {
		return "",err
	}

	//First make a tar archive for the public key generated above 
	buf := strings.NewReader(publicKey)
	if err == nil{
	err =cli.CopyToContainer(context.Background(), resp.ID, "/home/user/",buf ,types.CopyToContainerOptions{})
	if err!=nil{
		return "",err
	}}
	// Handle db call to store the resp.ID into the appropriate row for the user
	

	database_handling.CollectionHandler.UpdateOne(ctx,bson.M{"username":"CHANGE_USER_NAME_HERE"},bson.M{"CHANGE":"ownedContainers"})

return privateKey,nil
}



//Stop the container 
func ContainerStop(ctx context.Context,cli *client.Client,containerName string) error{
	//** Check if containerName is valid or ContainerStart requires id to start the container
	//Handle db call to retrieve the 'id' for the container required to start the container
	//var documentData map[string]interface{} 
	//Make db call to retrieve user info about the containers it holds
	//var err error
	_,err := get_document(ctx,"CHANGE_USER_NAME_HERE")

	if err !=nil{
	if err == mongo.ErrNoDocuments{
		//Give appropriate response
		panic("You haven't registered yet!\nRegister first")
	}else{
		return err //Here the system failure has occured
	}}
	var CHANGE time.Duration = 1029
    if err := cli.ContainerStop(ctx, "id",&CHANGE); err != nil {
        panic(err)
    }

    statusCh, errCh := cli.ContainerWait(ctx, "id", container.WaitConditionNotRunning)
    select {
    case err := <-errCh:
        if err != nil {
            panic(err)
        }
    case <-statusCh:
    }

	return nil

}


//Stop the container 
func ContainerRemove(ctx context.Context,cli *client.Client,containerName string) error{
	//** Check if containerName is valid or ContainerStart requires id to start the container
	//Handle db call to retrieve the 'id' for the container required to start the container
	//var documentData map[string]interface{} 
	//Make db call to retrieve user info about the containers it holds
	//var err error
	_,err := get_document(ctx,"CHANGE_USER_NAME_HERE")

	if err !=nil{
	if err == mongo.ErrNoDocuments{
		//Give appropriate response
		panic("You haven't registered yet!\nRegister first")
	}else{
		return err //Here the system failure has occured
	}}
	var options types.ContainerRemoveOptions
    if err := cli.ContainerRemove(ctx, "id",options); err != nil {
        panic(err)
    }

    statusCh, errCh := cli.ContainerWait(ctx, "id", container.WaitConditionNotRunning)
    select {
    case err := <-errCh:
        if err != nil {
            panic(err)
        }
    case <-statusCh:
    }

	return nil

}







//Start the container if already created
func ContainerStart(ctx context.Context,cli *client.Client,containerName string) (string,error){
	//** Check if containerName is valid or ContainerStart requires id to start the container
	//Handle db call to retrieve the 'id' for the container required to start the container
	var _ map[string]interface{} 
	//Make db call to retrieve user info about the containers it holds
	var err error
	_,err = get_document(ctx,"GET THE USER NAME SOMEHOW")
	if err == mongo.ErrNoDocuments{
		//Give appropriate response

	}else{
		return "",err //Here the system failure has occured
	}

    if err := cli.ContainerStart(ctx, "id", types.ContainerStartOptions{}); err != nil {
        panic(err)
    }

    statusCh, errCh := cli.ContainerWait(ctx, "id", container.WaitConditionNotRunning)
    select {
    case err := <-errCh:
        if err != nil {
            panic(err)
        }
    case <-statusCh:
    }

	privateKey,publicKey,err:= MakeSSHKeyPair()
	if err!=nil {
		panic(err)
	}

	//First make a tar archive for the public key generated above 
	buf := strings.NewReader(publicKey)
	err =cli.CopyToContainer(context.Background(), "id", "/home/user/", buf,types.CopyToContainerOptions{})
	if err!=nil{
		panic(err)
	}

	return privateKey,nil

}

//Gives information about the containers that user holds
func OwnedContainerInfo(ctx context.Context,cli *client.Client)(string,error){
	var _ map[string]interface{} 
	//Make db call to retrieve user info about the containers it holds
	var err error
	_,err = get_document(ctx,"username")
	if err == mongo.ErrNoDocuments{
		//Give appropriate response
		return "",err

	}else{
		return "",err //here system failure has happened
	}

}



//Gives images available on the server (the information is available on "system_details" collection)
func ImageInfo(ctx context.Context,cli *client.Client){
	
	//Make db call to retrieve the available ssh-able images 


}



func getVM(){

//Make the db call to ensure one user only gets only one time access to the VM



}











