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
  NextShubh int64
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

func otherPhase(p Phase) Phase {
  if p == Day {
    return Night
  }
  return Day
}

func getSoonestShubhTime(list map[string]int64) int64 {
  // Arbitrary minimum, one more digit than are in timestamps at the moment
  // This means this script will stop working in the year 2286
  // Hope I learn how to write Golang before that
  min := int64(10000000000)
  for _, value := range list {
    if value < min {
      min = value
    }
  }
  return min
}

func getChowgadhiyaList(t time.Time) map[string]int64 {
  sunrise, sunset, nextSunrise := getVedicDay(t)

  var baseTime time.Time
  var nextBase time.Time
  var phase Phase
  var offsetInSeconds float64

  if t.Before(sunset) {
    // Daytime
    phase = Day
    baseTime = sunrise
    nextBase = sunset
    offsetInSeconds = (sunset.Sub(sunrise) / 8).Seconds()
  } else {
    // Nighttime
    phase = Night
    baseTime = sunset
    nextBase = nextSunrise
    offsetInSeconds = (nextSunrise.Sub(sunset) / 8).Seconds()
    debug("time difference:", nextSunrise.Sub(t).Hours())
  }

  todayList := getChowgadhiyaListFromWeekday(t.Weekday(), phase)
  // fmt.Println(todayList)
  nextDayList := getChowgadhiyaListFromWeekday(t.Weekday()+1, otherPhase(phase))
  // fmt.Println(nextDayList)


  cList := make(map[string]int64)

  for index, element := range todayList {
    delta := (float64(index) * offsetInSeconds)
    startTime := (int64(delta)+baseTime.Unix())
    if (startTime > t.Unix() && isChowgadhiyaConsideredShubh(element)) {
      cList[chowgadhiyaToStringMap[element]] = startTime
    }
  }

  if len(cList) == 0 {
    for index, element := range nextDayList {
      delta := (float64(index) * offsetInSeconds)
      startTime := (int64(delta)+nextBase.Unix())
      if (startTime > t.Unix() && isChowgadhiyaConsideredShubh(element)) {
        cList[chowgadhiyaToStringMap[element]] = startTime
      }
    }
  }

  return cList
}

func getChowgadhiyaResponse(w http.ResponseWriter, r *http.Request) {
  now := time.Now()
  chowgadhiya := getChowgadhiya(now)

  isShubh := isShubh(now)
  current := chowgadhiyaToStringMap[chowgadhiya]
  list := getChowgadhiyaList(now)
  nextShubh := getSoonestShubhTime(list)

  response := Response{isShubh, nextShubh, current, list}

  fmt.Println(response)
  jResponse, _ := json.Marshal(response)
  fmt.Fprintf(w, string(jResponse))
}

func determineListenAddress() (string, error) {
  port := os.Getenv("PORT")
  if port == "" {
    return "", fmt.Errorf("$PORT not set")
  }
  return ":" + port, nil
}

func main() {
  http.HandleFunc("/chowgadhiya", getChowgadhiyaResponse) // set router
  addr, err := determineListenAddress()
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("Listening on %s...\n", addr)
  err = http.ListenAndServe(addr, nil) // set listen port
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
