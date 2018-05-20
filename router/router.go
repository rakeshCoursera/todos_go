package router

import (
	"net/http"
    
    "../todo"
    "../logger"
    "github.com/gorilla/mux"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"TodoCreate",
		"POST",
		"/todos",
		todo.TodoCreate,
    },
    Route{
        "TodoShow",
        "GET",
        "/todos/{todoId}",
        todo.TodoShow,
	},
	Route{
		"TodoList",
		"GET",
		"/todos",
		todo.TodoList,
    },
    Route{
		"TodoUpdate",
		"PUT",
		"/todos/update/{todoId}",
		todo.TodoUpdate,
    },
    Route{
		"TodoDelete",
		"DELETE",
		"/todos/delete/{todoId}",
		todo.TodoDelete,
    },
}

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler
        handler = route.HandlerFunc
        handler = logger.Logger(handler, route.Name)

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)

    }
    return router
}