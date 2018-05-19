package todo

import (
    "fmt"
    "net/http"
    "encoding/json"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
		"github.com/aws/aws-sdk-go/service/dynamodb"
		"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
		"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// structs to hold info about an item
type Todo struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Due       string	`json:"due"`
}

var Todos []Todo

func TodoList(w http.ResponseWriter, r *http.Request) {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
			Region: aws.String("ap-south-1")},
	)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
		}
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// Create the Expression to fill the input struct with.
	// Get back the id, name, completed, and due
	proj := expression.NamesList(expression.Name("id"), expression.Name("name"), expression.Name("completed"), expression.Name("due"))

	expr, err := expression.NewBuilder().WithProjection(proj).Build()

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
		}
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ProjectionExpression: expr.Projection(),
		TableName: aws.String("todo_app_table"),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	fmt.Println("error: ", err)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // Internal Server Error
		if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
		}
	}

	

	for _, i := range result.Items {
		item := Todo{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(500) // Internal Server Error
			if err := json.NewEncoder(w).Encode(err); err != nil {
					panic(err)
			}
		}
		Todos = append(Todos, item)
	}

	fmt.Println("Items: ", Todos)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated) // 200
	if err := json.NewEncoder(w).Encode(Todos); err != nil {
			panic(err)
	}
}