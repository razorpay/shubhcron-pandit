package main

import (
    "fmt"
    "net/http"
    "log"
    "time"
    "encoding/json"
)

func getChowgadhiyaResponse(w http.ResponseWriter, r *http.Request) {
    now := time.Now()
    chowgadhiya := getChowgadhiya(now)

    response := make(map[string]string)
    response["current"] = chowgadhiyaToStringMap[chowgadhiya]
    response["list"] = getChowgadhiyaList(now)

    jresponse, _ := json.Marshal(response)

    fmt.Fprintf(w, string(jresponse))
}

func main() {
    http.HandleFunc("/chowgadhiya", getChowgadhiyaResponse) // set router
    err := http.ListenAndServe(":9090", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}