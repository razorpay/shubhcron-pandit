package main

import (
    "fmt"
    "net/http"
    "log"
    "time"
    "encoding/json"
)

var chowgadhiyaToStringMap = map[Chowgadhiya]string{
  Chal  : "chal",
  Amrit : "amrit",
  Kaal  : "kaal",
  Labh  : "labh",
  Rog   : "rog",
  Shubh : "shubh",
  Udveg : "udveg",
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    now := time.Now()
    chowgadhiya := getChowgadhiya(now)

    response := make(map[string]string)
    response["current"] = chowgadhiyaToStringMap[chowgadhiya]

    jresponse, _ := json.Marshal(response)

    fmt.Fprintf(w, string(jresponse))
}


func main() {
    http.HandleFunc("/chowgadhiya", sayhelloName) // set router
    err := http.ListenAndServe(":9090", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}