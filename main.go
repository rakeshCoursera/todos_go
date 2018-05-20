package main

import (
    "log"
    "net/http"
    "strconv"

    "./config"
)

func main() {
    configs:= config.LoadConfiguration()
 
    router := NewRouter()

    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(configs.Port), router))
}