package query_handling

import (

	"context"
	"net/http"
	"io"
	"encoding/json"
	as "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/auth_service"
	db "github.com/VaradBelwalkar/Private-Cloud-MongoDB/api/database_handling"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)
// Here the user will be authenticated first and then request will be fulfilled

type storageInfoStruct struct{
	storageInfo map[string]string
}


//Function to retrieve file-data from request and upload it to the MongoDB 
func Upload_file(w http.ResponseWriter, r *http.Request) {
	check,username:=as.Handle_auth(w,r)
	if check!=true{
		return
	}
	// Read the file from the request
	filename:= r.Form.Get("filename")
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	defer file.Close()

	// Convert the file to a byte slice
	buf := make([]byte, 1024)
	var fileBytes []byte
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if n == 0 {
			break
		}
		fileBytes = append(fileBytes, buf[:n]...)
	}

	// Check if a document with the specified username already exists
	var result bson.M
	err = db.CollectionHandler.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} 
		


	filter:=bson.M{
		"username":username,
	}
	update:=bson.M{
		"storage":bson.M{"$set":bson.M{filename:bson.M{"type":"file","bytes":fileBytes}}},
	}

	updateResult,err:=db.CollectionHandler.UpdateOne(context.TODO(),filter,update)
	if err!=nil || updateResult.MatchedCount!=1{
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusOK)


}




//Function to retrieve file-data from request and upload it to the MongoDB 
func Upload_folder(w http.ResponseWriter, r *http.Request) {
	check,username:=as.Handle_auth(w,r)
	if check!=true{
		return
	}
	// Read the file from the request
	foldername:=r.Form.Get("foldername")
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer file.Close()

	// Convert the file to a byte slice
	buf := make([]byte, 1024)
	var fileBytes []byte
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if n == 0 {
			break
		}
		fileBytes = append(fileBytes, buf[:n]...)
	}

	// Check if a document with the specified username already exists
	var result bson.M
	err = db.CollectionHandler.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} 
		


	filter:=bson.M{
		"username":username,
	}
	update:=bson.M{
		"storage":bson.M{"$set":bson.M{foldername:bson.M{"type":"folder","bytes":fileBytes}}},
	}

	updateResult,err:=db.CollectionHandler.UpdateOne(context.TODO(),filter,update)
	if err!=nil || updateResult.MatchedCount!=1{
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusOK)


}


func fetch_document(ctx context.Context,username string)(map[string]interface{}, int){

	var documentData map[string]interface{} 
	//Check user-document exists in the collection 
	//document_handler of type *SingleResult, see github code for more details
	err := db.CollectionHandler.FindOne(ctx,bson.M{"username":username}).Decode(&documentData)
	//If not then use following	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil,404
		} else {
			return nil,500
		}
	}	

return documentData,200

}


func storageInfo(w http.ResponseWriter, r *http.Request){
	check,username:=as.Handle_auth(w,r)
	if check!=true{
		return
	}

	documentData,err := fetch_document(context.TODO(),username)
	if err!=200{
		w.WriteHeader(http.StatusInternalServerError)
	}

	var containerArray map[string]string
	for k, _ := range documentData["storage"].(map[string]interface{}) { 

		containerArray[k]=documentData["storage"].(map[string]interface{})[k].(map[string]string)["type"]
		
	}
	
	resp:=storageInfoStruct{storageInfo:containerArray}
	json.NewEncoder(w).Encode(resp)
	w.Header().Set("Content-Type", "application/json")

}