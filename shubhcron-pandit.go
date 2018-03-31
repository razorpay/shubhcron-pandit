package main

import (
  "os"
  "fmt"
  "net/http"
  "log"
  "time"
  "encoding/json"
)

type Response struct {
  IsShubh bool
  Current string `json:"current"`
  List ChowgadhiyaTimeList `json:"list"`
}

type ChowgadhiyaTimeList map[string]int64

var chowgadhiyaToStringMap = map[Chowgadhiya]string{
  Chal  : "chal",
  Amrit : "amrit",
  Kaal  : "kaal",
  Labh  : "labh",
  Rog   : "rog",
  Shubh : "shubh",
  Udveg : "udveg",
}

func getChowgadhiyaList(t time.Time) map[string]int64 {
  sunrise, sunset, nextSunrise := getVedicDay(t)

  var baseTime time.Time
  var phase Phase
  var offsetInSeconds float64

  if t.Before(sunset) {
    // Daytime
    phase = Day
    baseTime = sunrise
    offsetInSeconds = (sunset.Sub(sunrise) / 8).Seconds()
  } else {
    // Nighttime
    phase = Night
    baseTime = sunset
    offsetInSeconds = (nextSunrise.Sub(sunset) / 8).Seconds()
    debug("time difference:", nextSunrise.Sub(t).Hours())
  }

  list := getChowgadhiyaListFromWeekday(t.Weekday(), phase)

  cList := make(map[string]int64)

  for index, element := range list {
    delta := (float64(index) * offsetInSeconds)
    cList[chowgadhiyaToStringMap[element]] = (int64(delta)+baseTime.Unix())
  }

  return cList
}

func getChowgadhiyaResponse(w http.ResponseWriter, r *http.Request) {
  now := time.Now()
  chowgadhiya := getChowgadhiya(now)

  isShubh := isShubh(now)
  current := chowgadhiyaToStringMap[chowgadhiya]
  list := getChowgadhiyaList(now)

  response := Response{isShubh, current, list}

  fmt.Println(response)
  jResponse, _ := json.Marshal(response)
  fmt.Fprintf(w, string(jResponse))
}

func determineListenAddress() (string, error) {
  port := os.Getenv("PORT")
  if port == "" {
    return "", fmt.Errorf("$PORT not set")
  }
  fmt.Println("PORT")
  return ":" + port, nil
}

func main() {
  http.HandleFunc("/chowgadhiya", getChowgadhiyaResponse) // set router
  addr, err := determineListenAddress()
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("Listening on %s...\n", addr)
  fmt.Println(addr)
  err = http.ListenAndServe(addr, nil) // set listen port
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
