package container_apis

import (
	"context"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"golang.org/x/crypto/ssh"
	"encoding/pem"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling/mongodb"
)

type resultStruct struct{
	Username string `bson:"username"`
	Password string `bson:"password"`
    Email string `bson:"email"`
	ContainerInfo map[string]interface{} `bson:"containerInfo"`
	TotalOwnedContainers int `bson:"totalOwnedContainers,omitempty"`
}

var ImageArray []string

func checkImage(userImage string) bool{
	for _,availImage := range ImageArray{
		if availImage == userImage{
			return true
		}
	}
	return false
}


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

    return  privKeyBuf.String(),pubKeyBuf.String(), 200
}


//Function to return err if document not found
func get_document(ctx context.Context,username string)(resultStruct, int){

	documentData:=resultStruct{}
	//Check user-document exists in the collection 
	//document_handler of type *SingleResult, see github code for more details
	err := db.CollectionHandler.FindOne(ctx,bson.M{"username":username}).Decode(&documentData)
	//If not then use following	
	if err != nil {		
			return resultStruct{},500	// Internal server error
	}	
return documentData,200

}













