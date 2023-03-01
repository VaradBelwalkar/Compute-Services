package containers

import (

	"context"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	ca "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/container_apis"
	as "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service"
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
