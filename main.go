package main

import (
    "log"
    "net/http"
    "strconv"

    "./config"
    "./router"
)

func main() {
    configs:= config.LoadConfiguration()
 
    router := router.NewRouter()

    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(configs.Port), router))
}