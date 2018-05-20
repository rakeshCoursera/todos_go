package todo

import (
    "fmt"
    "time"
    "net/http"
    "io/ioutil"
    "encoding/json"

    "../config"
    "github.com/rs/xid"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func TodoCreate(w http.ResponseWriter, r *http.Request) {
    // load environment variables
	configs:= config.LoadConfiguration()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(configs.Region),
		Credentials: credentials.NewStaticCredentials(configs.Access_Key, configs.Secret_Key, configs.Access_Token),
	})

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
        Id: xid.New().String(),
        Title: todo.Title,
        Completed: todo.Completed,
        Due: todo.Due,
        CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
        UpdatedAt: "",
    }
		
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
        TableName: aws.String(configs.Table_Name),
    }

    _, err = svc.PutItem(input)

    if err != nil {
        fmt.Println(err)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(500) // Internal Server Error
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    } else {
        fmt.Printf("Successfully added a item in the table: %s \n", configs.Table_Name)
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated) // 200
	if err := json.NewEncoder(w).Encode(item); err != nil {
		panic(err)
	}
}
