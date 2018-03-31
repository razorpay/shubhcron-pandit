package main

import (
  "github.com/kelvins/sunrisesunset"
  "io/ioutil"
  "log"
  "math"
  "os"
  "strconv"
  "time"
)

type Chowgadhiya int
type Phase int

const (
  Day Phase = iota
  Night
)

const (
  Chal Chowgadhiya = iota
  Amrit
  Kaal
  Labh
  Rog
  Shubh
  Udveg
)

// Ultimate default coordinates
const DEFAULT_LATITUDE  string = "26.7880"
const DEFAULT_LONGITUDE string = "82.1986"

// https://hinduism.stackexchange.com/questions/26242/how-is-the-first-choghadiya-decided
// Golang does not allow constant maps, but a literal map is close enough
var CHOWGADHIYA_LIST = map[Phase]map[time.Weekday][]Chowgadhiya{
  Day: map[time.Weekday][]Chowgadhiya{
    time.Sunday    : []Chowgadhiya{Udveg , Chal  , Labh  , Amrit , Kaal  , Shubh , Rog   , Udveg} ,
    time.Monday    : []Chowgadhiya{Amrit , Kaal  , Shubh , Rog   , Udveg , Chal  , Labh  , Amrit} ,
    time.Tuesday   : []Chowgadhiya{Rog   , Udveg , Chal  , Labh  , Amrit , Kaal  , Shubh , Rog}   ,
    time.Wednesday : []Chowgadhiya{Labh  , Amrit , Kaal  , Shubh , Rog   , Udveg , Chal  , Labh}  ,
    time.Thursday  : []Chowgadhiya{Shubh , Rog   , Udveg , Chal  , Labh  , Amrit , Kaal  , Shubh} ,
    time.Friday    : []Chowgadhiya{Chal  , Labh  , Amrit , Kaal  , Shubh , Rog   , Udveg , Chal}  ,
    time.Saturday  : []Chowgadhiya{Kaal  , Shubh , Rog   , Udveg , Chal  , Labh  , Amrit , Kaal}  ,
  },
  Night: map[time.Weekday][]Chowgadhiya{
    time.Sunday    : []Chowgadhiya{Shubh , Amrit , Chal  , Rog   , Kaal  , Labh  , Udveg , Shubh} ,
    time.Monday    : []Chowgadhiya{Chal  , Rog   , Kaal  , Labh  , Udveg , Shubh , Amrit , Chal}  ,
    time.Tuesday   : []Chowgadhiya{Kaal  , Labh  , Udveg , Shubh , Amrit , Chal  , Rog   , Kaal}  ,
    time.Wednesday : []Chowgadhiya{Udveg , Shubh , Amrit , Chal  , Rog   , Kaal  , Labh  , Udveg} ,
    time.Thursday  : []Chowgadhiya{Amrit , Chal  , Rog   , Kaal  , Labh  , Udveg , Shubh , Amrit} ,
    time.Friday    : []Chowgadhiya{Rog   , Kaal  , Labh  , Udveg , Shubh , Amrit , Chal  , Rog}   ,
    time.Saturday  : []Chowgadhiya{Labh  , Udveg , Shubh , Amrit , Chal  , Rog   , Kaal  , Labh}  ,
  },
}

var chowgadhiyaToStringMap = map[Chowgadhiya]string{
  Chal  : "chal",
  Amrit : "amrit",
  Kaal  : "kaal",
  Labh  : "labh",
  Rog   : "rog",
  Shubh : "shubh",
  Udveg : "udveg",
}

/**
 * Returns the list of Chowgadhiyas in Order for Daytime
 */
func getChowgadhiyaListFromWeekday(day time.Weekday, phase Phase) []Chowgadhiya {
  return CHOWGADHIYA_LIST[phase][day]
}

/**
 * Takes time and returns the correct Chowgadhiya
 */
func getChowgadhiya(t time.Time) Chowgadhiya {
  sunrise, sunset, nextSunrise := getVedicDay(t)

  debug("Next sunrise:", nextSunrise)
  debug("Current time:", t)

  if t.Before(sunrise) || t.After(nextSunrise) {
    panic("current time does not fall between Sunrise and Sunset")
  }

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

  timePassedInCurrentPhase := t.Sub(baseTime).Seconds()
  debug("timePassedInCurrentPhase:", timePassedInCurrentPhase)
  debug("offsetInSeconds:", offsetInSeconds)
  numberOfChowgadhiyaPassed := timePassedInCurrentPhase / offsetInSeconds
  debug("numberOfChowgadhiyaPassed:", numberOfChowgadhiyaPassed)
  chowgadhiyaIndex := int(math.Floor(numberOfChowgadhiyaPassed))
  debug("chowgadhiyaIndex:", chowgadhiyaIndex)
  list := getChowgadhiyaListFromWeekday(sunrise.Weekday(), phase)
  debug("phase:", phase)
  debug("list:", list)
  return list[chowgadhiyaIndex]
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

func getEnv(key, fallback string) string {
  if value, ok := os.LookupEnv(key); ok {
    return value
  }
  return fallback
}

func getSunriseSunset(t time.Time) (time.Time, time.Time) {
  reference_time := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

  _, offset := t.Zone()

  fractional_offset := (float64(offset)/60/60);

  if fractional_offset > 12 {
    fractional_offset = 12 - fractional_offset
  }

  latitude, _  := strconv.ParseFloat(getEnv("LATITUDE", DEFAULT_LATITUDE), 64)
  longitude, _ := strconv.ParseFloat(getEnv("LONGITUDE", DEFAULT_LONGITUDE), 64)

  p := sunrisesunset.Parameters{
    Latitude:  latitude,
    Longitude: longitude,
    UtcOffset: fractional_offset,
    Date:      reference_time,
  }

  sunrise, sunset, err := p.GetSunriseSunset()

  if err == nil {
    return sunrise, sunset
  }
  panic("sunrise/sunset calculations failed")
}

func debug(strings ...interface{}) {
  Debug := log.New(os.Stdout,
    "DEBUG:",
    log.Ldate|log.Ltime|log.Lshortfile)

  _, debug := os.LookupEnv("DEBUG")

  if debug != true {
    Debug.SetOutput(ioutil.Discard)
  }

  Debug.Println(strings...)
}

func getVedicDay(now time.Time) (time.Time, time.Time, time.Time) {

  var sunrise, sunset, nextSunrise time.Time

  sunrise, sunset = getSunriseSunset(now)

  yesterday := now.AddDate(0, 0, -1)
  tomorrow := now.AddDate(0, 0, 1)

  loc := now.Location()
  tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
  yesterday = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())

  sunrise = time.Date(now.Year(), now.Month(), now.Day(), sunrise.Hour(), sunrise.Minute(), sunrise.Second(), sunrise.Nanosecond(), loc)
  sunset = time.Date(now.Year(), now.Month(), now.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), sunset.Nanosecond(), loc)

  // Sun has not risen yet
  // So check the sunrise for yesterday
  if now.Before(sunrise) {
    debug("Sun is not yet up, go back to bed")
    nextSunrise, sunset = getSunriseSunset(yesterday)

    sunset = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), sunset.Hour(), sunset.Minute(), sunset.Second(), sunset.Nanosecond(), loc)
    nextSunrise = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), nextSunrise.Hour(), nextSunrise.Minute(), nextSunrise.Second(), nextSunrise.Nanosecond(), loc)
  } else {
    debug("Sun is up, rise and shine")
    // Calculate the sunrise time for tomorrow
    nextSunrise, _ = getSunriseSunset(tomorrow)
    nextSunrise = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), nextSunrise.Hour(), nextSunrise.Minute(), nextSunrise.Second(), nextSunrise.Nanosecond(), loc)
  }

  // Now we have a definite sunrise time for the "vedic day"

  debug("Sunrise:", sunrise)
  debug("Sunset:", sunset)
  debug("Next sunrise:", nextSunrise)

  return sunrise, sunset, nextSunrise
}
