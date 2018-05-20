package todo

import (
    "net/http"
    "encoding/json"

    "../config"
    "github.com/gorilla/mux"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/expression"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// structs to hold info about an item
type Todo struct {
	Id					string		`json:"id"`
	Title				string		`json:"title"`
	Completed		bool			`json:"completed"`
	Due					string		`json:"due"`
	CreatedAt		string		`json:"createdAt"`
	UpdatedAt		string		`json:"updatedAt"`
}


func TodoShow(w http.ResponseWriter, r *http.Request) {
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

	// Create the Expression to fill the input struct with.
	// Get all item for the requested todoId
	filt := expression.Name("id").Equal(expression.Value(todoId))

	// Get back the id, name, completed, and due
	proj := expression.NamesList(expression.Name("id"), expression.Name("title"), expression.Name("completed"), expression.Name("due"), expression.Name("createdAt"), expression.Name("updatedAt"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

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
        FilterExpression: expr.Filter(),
        ProjectionExpression: expr.Projection(),
		TableName: aws.String(configs.Table_Name),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

    if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(500) // Internal Server Error
		if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
		}
    }

    Todos:= []Todo{}
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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated) // 200
	if err := json.NewEncoder(w).Encode(Todos); err != nil {
			panic(err)
	}
}