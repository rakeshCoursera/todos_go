# Golang Todo REST API Example
A RESTful API example for simple todo application with Golang and DynamoDB.

## Prerequisite
* Should have an AWS `access key` and `secret key`. This can be get by creating a user on AWS with programatic access.
* Must create a DynamoDB table with primary key as `id` with type `string`.


## Installation & Run

```bash
# clone the project
> git clone github.com/rakeshCoursera/todos_go

# go into the project folder
> cd todos_go 
```
In the project folder look for file `config.json` in `config` folder and fill the values:
```
{
  "port": 3000,
  "region" : "ap-south-1",
  "access_key": "enter your aws access key",
  "secret_key": "enter your aws secret key",
  "access_token": "this is optional, leave empty",
  "table_name": "DynamoDB table name created in 2 step"
}
```

```
# download the packages used in the project
> go get

# build the project
> go build

# run the project
> go run main.go

# API Endpoint : http://localhost:3000
```

## APIs

#### /todos
* `GET` : Get all todos tasks
* `POST` : Create a new task

#### /todos/{todoId}
* `GET` : Get a todo task

#### /todos/update/{todos}
* `PUT` : Update a task

#### /todos/delete/{todos}
* `DELETE`: delete a task


