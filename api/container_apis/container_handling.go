package cotainer_apis

import (
	"context"
	"io"
	"os"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gocql/gocql"
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

func getTar(publicKey string) bytes.Buffer{

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err = tw.WriteHeader(&tar.Header{
		Name: "authorized_keys",            // filename
		Mode: 600,                // permissions
		Size: int64(len(publicKey)), // filesize
	})
	if err != nil {
		return nil, fmt.Errorf("docker copy: %v", err)
	}
	tw.Write([]byte(publicKey))
	tw.Close()
	return buf;
}


//Function to return err if document not found
func get_document(ctx context.Context,username string)(map[string]interface{}, err){

	var documentData map[string]interface{} 
	//Check user-document exists in the collection 
	//document_handler of type *SingleResult, see github code for more details
	err := CollectionHandler.Findone(ctx,bson.M{"username":userName}).Deocde(&documentData)
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
func ContainerCreate(ctx context.Context,cli *Client,imageName string){

	var documentData map[string]interface{} 
	documentData, err = get_document()
	if err!= nil{
		return err
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
        panic(err)
    }

	privateKey,publicKey,err:= MakeSSHKeyPair()
	if err!=nil {
		panic(err)
	}

	//First make a tar archive for the public key generated above 

	err:=cli.CopyToContainer(context.Background(), resp.ID, "/home/user/", getTar(publicKey),types.CopyToContainerOptions{})
	if err!=nil{
		panic(err)
	}
		// Handle db call to store the resp.ID into the appropriate row for the user
	

	CollectionHandler.UpdateOne(ctx,bson.M{"username":userName},bson.M{ownedContainers})

}



//Stop the container 
func ContainerStop(ctx context.Context,cli *Client,containerName string){
	//** Check if containerName is valid or ContainerStart requires id to start the container
	//Handle db call to retrieve the 'id' for the container required to start the container
	var documentData map[string]interface{} 
	//Make db call to retrieve user info about the containers it holds
	documentData,err = get_document(ctx,username)
	if err == mongo.ErrNoDocuments{
		//Give appropriate response
		fmt.Fprintf(w, "You haven't registered yet!\nRegister first")
	}else{
		return err //Here the system failure has occured
	}

    if err := cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
        panic(err)
    }

    statusCh, errCh := cli.ContainerWait(ctx, id, container.WaitConditionNotRunning)
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

	err:=cli.CopyToContainer(context.Background(), id, "/home/user/", getTar(publicKey),types.CopyToContainerOptions{})
	if err!=nil{
		panic(err)
	}

}



//Start the container if already created
func ContainerStart(ctx context.Context,cli *Client,containerName string){
	//** Check if containerName is valid or ContainerStart requires id to start the container
	//Handle db call to retrieve the 'id' for the container required to start the container
	var documentData map[string]interface{} 
	//Make db call to retrieve user info about the containers it holds
	documentData,err = get_document(ctx,username)
	if err == mongo.ErrNoDocuments{
		//Give appropriate response

	}else{
		return err //Here the system failure has occured
	}

    if err := cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
        panic(err)
    }

    statusCh, errCh := cli.ContainerWait(ctx, id, container.WaitConditionNotRunning)
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

	err:=cli.CopyToContainer(context.Background(), id, "/home/user/", getTar(publicKey),types.CopyToContainerOptions{})
	if err!=nil{
		panic(err)
	}

}

//Gives information about the containers that user holds
func OwnedContainerInfo(ctx context.Context,cli *Client){
	var documentData map[string]interface{} 
	//Make db call to retrieve user info about the containers it holds
	documentData,err = get_document(ctx,username)
	if err == mongo.ErrNoDocuments{
		//Give appropriate response

	}else{
		return err //here system failure has happened
	}

}



//Gives images available on the server (the information is available on "system_details" collection)
func ImageInfo(ctx context.Context,cli *Client){
	
	//Make db call to retrieve the available ssh-able images 


}



func getVM(){

//Make the db call to ensure one user only gets only one time access to the VM



}











