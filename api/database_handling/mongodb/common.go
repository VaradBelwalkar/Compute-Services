package mongodb

import (
	"crypto/hmac"
    "crypto/sha256"
	"crypto/subtle"
    "encoding/hex"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var PassHashKey string

//Global Object for user details
type resultStruct struct{
	Username string `bson:"username"`
	Password string `bson:"password"`					// In the hash format
    Email string `bson:"email"`
	ContainerInfo map[string]interface{} `bson:"containerInfo"`
	TotalOwnedContainers int `bson:"totalOwnedContainers,omitempty"`
}

var CollectionHandler *mongo.Collection 
var Sys_CollectionHandler *mongo.Collection 
var DatabaseHandler *mongo.Database

//Register system information here (e.g docker images available)
var sys_info = bson.M {"server":"private_cloud","docker_images":[]string{"ubuntu","nginx"}}

type DBConfig struct {
	Name string     	 `json:"name"`
	Collections []string `json:"collections"`
}

type Config struct {
	Emails []string `json:"emails"`
	Database DBConfig `json:"database"`
}

func ComputeHash(password string) string {
    key := []byte(PassHashKey)
    h := hmac.New(sha256.New, key)
    h.Write([]byte(password))
    hash := hex.EncodeToString(h.Sum(nil))
    return hash
}

func compareHashAndPassword(hash, password string) bool {
    expectedHash := ComputeHash(password)
    return subtle.ConstantTimeCompare([]byte(hash), []byte(expectedHash)) == 1
}



