package todo

import (
	"fmt"
	"net/http"
	"encoding/json"

	"../config"
	"github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TodoDelete(w http.ResponseWriter, r *http.Request) {
  // load environment variables
	configs:= config.LoadConfiguration()

	vars := mux.Vars(r)
	todoId := vars["todoId"]

	sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(configs.Region),
			Credentials: credentials.NewStaticCredentials(configs.Access_Key, configs.Secret_Key, configs.Access_Token),
	})
	
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(401) // unauthorised request
		if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
		}
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(todoId),
			},
		},
		TableName: aws.String(configs.Table_Name),
	}

	_, err = svc.DeleteItem(input)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // Internal Server Error
		if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
		}
	} else {
		fmt.Println("Successfully deleted a item with Id: ", todoId)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated) // 200
	if err := json.NewEncoder(w).Encode("Item successfully deleted"); err != nil {
		panic(err)
	}
}