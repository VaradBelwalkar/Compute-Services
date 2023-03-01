package container_apis

import (
	"context"
	"github.com/docker/docker/client"
)


//Gives information about the containers that user holds
func OwnedContainerInfo(ctx context.Context,username string)(map[string]string,int){

	documentData,err := get_document(context.TODO(),username)
	if err!=200{
		return nil,err
	}
	containerMap:=make(map[string]string)
	for k, v := range documentData.ContainerInfo {
		containerMap[k]=v.(map[string]interface{})["status"].(string)
	}
	return containerMap,200

}



//Gives images available on the server (the information is available on "system_details" collection)
func ImageInfo(ctx context.Context,cli *client.Client) []string{
	
	var imageArray=[]string{"ubuntu","development_server"}

	return imageArray

}











