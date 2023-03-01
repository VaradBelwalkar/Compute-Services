package containers

import (
	"github.com/docker/docker/client"
)
// Here the user will be authenticated first and then request will be fulfilled

var Cli *client.Client
type responseStruct struct{
	privatekey string
	port string
}

type infoStruct struct{
	containerInfo []string
}
