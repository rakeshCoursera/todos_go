package todo

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"../config"
	"github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TodoUpdate(w http.ResponseWriter, r *http.Request) {
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
			panic(err)
	}

	var todo Todo
	err = json.Unmarshal(body, &todo)
	if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
					panic(err)
			}
	}

	// Create item in table Movies
	input := &dynamodb.UpdateItemInput{
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":c": {
				BOOL: aws.Bool(todo.Completed),
			},
			":t": {
				S: aws.String(todo.Title),
			},
			":d": {
				S: aws.String(todo.Due),
			},
			":u": {
				S: aws.String(time.Now().Format("2006-01-02 15:04:05")),
			},
		},
		TableName: aws.String(configs.Table_Name),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(todoId),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set completed = :c, title = :t, due = :d, updatedAt = :u"),
  }

  _, err = svc.UpdateItem(input)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // Internal Server Error
		if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
		}
	} else {
		fmt.Println("Successfully updated the table", configs.Table_Name)
	}
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated) // 200
	if err := json.NewEncoder(w).Encode("Successfully Updated"); err != nil {
		panic(err)
	}
}