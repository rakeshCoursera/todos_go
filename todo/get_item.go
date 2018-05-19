package main

import (
    "fmt"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// structs to hold info about an item
type Todo struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Due       string		`json:"due"`
}

func TodoList(w http.ResponseWriter, r *http.Request) {
    // Initialize a session in us-west-2 that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials.
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-south-1")},
    )

    // Create DynamoDB client
    svc := dynamodb.New(sess)

    result, err := svc.GetItem(&dynamodb.GetItemInput{
        TableName: aws.String("todo_app_table"),
        Key: map[string]*dynamodb.AttributeValue{
            "id": {
                N: aws.String("2"),
            },
            "due": {
                S: aws.String("The Big New Movie"),
            },
        },
    })

    if err != nil {
        fmt.Println(err.Error())
        return
    }

    item := Item{}

    err = dynamodbattribute.UnmarshalMap(result.Item, &item)

    if err != nil {
        panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
    }

    if item.Title == "" {
        fmt.Println("Could not find 'The Big New Movie' (2015)")
        return
    }

    fmt.Println("Found item:")
    fmt.Println("Year:  ", item.Year)
    fmt.Println("Title: ", item.Title)
    fmt.Println("Plot:  ", item.Info.Plot)
    fmt.Println("Rating:", item.Info.Rating)
}