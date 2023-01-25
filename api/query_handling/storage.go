package query_handling

import (

	"fmt"
	"net/http"
	"github.com/VaradBelwalkar/Private-Cloud/api/container_apis"
	"path/filepath"
	"github.com/docker/docker/client"
)
// Here the user will be authenticated first and then request will be fulfilled

func upload_file(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()


}

func upload_folder(w http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()

}
