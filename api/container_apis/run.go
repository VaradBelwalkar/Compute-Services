package container_apis

import (
	"context"
	"io/ioutil"
	"strings"
	"bytes"
	"archive/tar"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"
)

//Create a new container
func ContainerCreate(ctx context.Context,cli *client.Client,imageName string,username string) (string,string,int){
	documentData, status := get_document(context.TODO(),username)
	if status!= 200{
		return "","",status
	}
	//Here we get the document to work with

	if !checkImage(imageName) {
		return "","",404
	}
	
	totalOwnedContainers := documentData.TotalOwnedContainers
	
		if totalOwnedContainers>=5 {
			return "","",403		// http.StatusForbidden  //This is because system cannot allocate more than 5 containers per user
		}
	
	//Else do allocate the container	
	containerCfg := &container.Config {
		Image: imageName,
		AttachStdin:false,
		AttachStdout:false,
		AttachStderr:false,
		OpenStdin:false,
		Cmd: []string{"service","ssh","start", "-D", "daemon on;"},
		ExposedPorts: nat.PortSet{
			//nat.Port("443/tcp"): {},
			nat.Port("22/tcp"): {},
		},
	}	
	volumeBinding:=username+":/mnt:rw"
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			//nat.Port("443/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "443"}},
			nat.Port("22/tcp"): []nat.PortBinding{{HostIP: "0.0.0.0"}}, //Here excluding HostPort to assign a random port 
		},
		Binds: []string{volumeBinding},
	}
	
    resp, err := cli.ContainerCreate(ctx,containerCfg,hostConfig,nil,nil,"")
    if err != nil {
        return "","",500
    }

    if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
        return "","",500
    }
	
	privateKey,publicKey,check:= MakeSSHKeyPair()
	if check!=200 {
		
		return "","",check
	}
	
	//First make a tar archive for the public key generated above 
	//buf := strings.NewReader(publicKey)

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

	t, err := ioutil.ReadAll(r)
	if err != nil {
		return "","",500
	}
	tempString:=string(t)
	readBuf := strings.NewReader(tempString)
	err =cli.CopyToContainer(context.Background(), resp.ID, "/root/.ssh/",readBuf ,types.CopyToContainerOptions{})
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
	//opts := options.Update().SetUpsert(true)
	//fmt.Println(opts)
	filter:=bson.M{
		"username":username,
	}

	update:=bson.M{ "$set":bson.M{
		"totalOwnedContainers":totalOwnedContainers+1,
		"containerInfo."+containerName:bson.M{"containerID":resp.ID,"port":port,"status":"running"},
	}}

	updateResult,err:=db.CollectionHandler.UpdateOne(ctx,filter,update)
	if err!=nil || updateResult.MatchedCount!=1{
		return "","",500
	}
return privateKey,port,200
}

