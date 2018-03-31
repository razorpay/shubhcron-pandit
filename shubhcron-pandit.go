package main

import (
    "fmt"
    "net/http"
    "log"
    "time"
    "encoding/json"
)

type Response struct {
  current string
  list ChowgadhiyaTimeList
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

    current := chowgadhiyaToStringMap[chowgadhiya]
    list := getChowgadhiyaList(now)

    response := Response{current, list}

    jresponse, err := json.MarshalIndent(response, "", "  ")

    if err != nil {
        panic(err)
    }
    fmt.Fprintf(w, string(jresponse))
}

func main() {
    http.HandleFunc("/chowgadhiya", getChowgadhiyaResponse) // set router
    err := http.ListenAndServe(":9090", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}