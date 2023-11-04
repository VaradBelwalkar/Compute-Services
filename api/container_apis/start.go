package container_apis

import (
	"context"
	"io/ioutil"
	"strings"
	"bytes"
	"archive/tar"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	db "github.com/VaradBelwalkar/Compute-Services/api/database_handling/mongodb"
)


//Start the container if already created
func ContainerStart(ctx context.Context,cli *client.Client,containerName string,username string) (string,string,int){
	//** Check if containerName is valid or ContainerStart requires id to start the container
	//Handle db call to retrieve the 'id' for the container required to start the container

	documentData,status := get_document(ctx,username)
	if status!=200{
		return "","",status //Here the system failure has occured
	}


	var containerID string

	nesting1:=documentData.ContainerInfo//[containerName].(map[string]interface{})["containerID"].(string)
	if nesting2,ok:= nesting1[containerName]; ok{
		containerID=nesting2.(map[string]interface{})["containerID"].(string)
	} else{
		return "","",404 // Send StatusNotFound
	}
	// If returns error if container is already running, first do inspect the container and then only run



	containerINFO, errCase := cli.ContainerInspect(context.Background(), containerID)
	if errCase != nil {
		return "","",500
	}
	var old_container_ip string
	if containerINFO.State.Running == false{
    if err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
        return "","",500
    }
} else {
	old_container_ip=strings.Replace(containerINFO.NetworkSettings.Networks["docker-dhcp"].IPAddress, "_", ".", -1)
}

	privateKey,publicKey,err:= MakeSSHKeyPair()
	if err!=200 {
		return "","",err
	}

	str := []byte(publicKey)
	b := new(bytes.Buffer)

	// Create a new tar archive
	tw := tar.NewWriter(b)

	// Add the string to the archive
	hdr := &tar.Header{
		Name: "authorized_keys",
		Mode: 0600,
		Size: int64(len(str)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return "","",500
	}
	if _, err := tw.Write(str); err != nil {
		return "","",500
	}
	if err := tw.Close(); err != nil {
		return "","",500
	}
	
	r := bytes.NewReader(b.Bytes())
	t, errInfo := ioutil.ReadAll(r)
	if errInfo != nil {
		return "","",500
	}
	tempString:=string(t)
	readBuf := strings.NewReader(tempString)
	check :=cli.CopyToContainer(context.Background(), containerID, "/root/.ssh/", readBuf,types.CopyToContainerOptions{})
	if check!=nil{
		return "","",500
	}
	if containerINFO.State.Running == false{
	containerJSON,check:=cli.ContainerInspect(context.TODO(),containerID)
	if check!=nil{
	return "","",500
	}
	container_ip:=containerJSON.NetworkSettings.Networks["docker-dhcp"].IPAddress
	newContainerName:=username+"_"+strings.Replace(container_ip, ".", "_", -1)
	

	filter:=bson.M{
		"username":username,
	}
	update:=bson.M{ "$unset":bson.M{
		"containerInfo."+containerName:"",
	},
	"$set":bson.M{
		"containerInfo."+newContainerName:bson.M{"containerID":containerID,"container_ip":strings.Replace(container_ip, ".", "_", -1),"status":"running"},
	},
	}
	updateResult,check:=db.CollectionHandler.UpdateOne(ctx,filter,update)
	if check!=nil || updateResult.MatchedCount!=1{
		return "","",500
	}
	
	return privateKey,container_ip,200
	
	}else{
		return privateKey,old_container_ip,200
	}

}












