package todo

import (
    "fmt"
    "net/http"
    "encoding/json"
    "io/ioutil"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func TodoCreate(w http.ResponseWriter, r *http.Request) {
    // Initialize a session in us-west-2 that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials.
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-south-1")},
    )

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

    item := Todo {
        Id: todo.Id,
        Name: todo.Name,
        Completed: todo.Completed,
        Due: todo.Due,
    }
		
	fmt.Println("Item: ", item)

    av, err := dynamodbattribute.MarshalMap(item)

    if err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(400) // Bad Request
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    // Create item in table Movies
    input := &dynamodb.PutItemInput{
        Item: av,
        TableName: aws.String("todo_app_table"),
    }

    _, err = svc.PutItem(input)

    if err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(500) // Internal Server Error
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    fmt.Println("Successfully added a item to todo_app table")

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated) // 200
	if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
	}
}
