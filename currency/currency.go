package currency

import (
  "net/http"
  "io/ioutil"
  "strings"
  "encoding/json"
  "strconv"
)

const endpoint string = "live"
const accessKey string = "" //enter access key here
const format string = "1"

type Conversion struct {
  ConversionRate float64
  ConvertedValue float64
  ValuetoConvert float64
  Target string
  Errors []string
  ErrorString string
}

func (c Conversion) ErrorPrettyPrint() string {
  var result string = ""
  for i, err := range c.Errors {
    result += err
    if i != len(c.Errors) - 1 {
      result += ", "
    }
  }
  return result
}

func Convert(target string) Conversion {
  var result Conversion
  result.ConversionRate = 0
  result.ConvertedValue = 0

  var url string = "http://apilayer.net/api/" + endpoint + "?access_key=" + accessKey + "&currencies=" + target + "&format=" + format

  req, _ := http.NewRequest("GET", url, nil)

  res, _ := http.DefaultClient.Do(req)

  defer res.Body.Close()
  body, _ := ioutil.ReadAll(res.Body)

  var status bool = strings.Contains(res.Status, "200")

  var raw map[string]interface{}
  json.Unmarshal(body, &raw)
  var successfulBytes, _ = json.Marshal(raw["success"])
  var successfulString = string(successfulBytes)
  var successfulBool, _ = strconv.ParseBool(successfulString)

  if !status || !successfulBool {
    result.Errors = append(result.Errors, "unsuccessful connection to conversion api")
  } else {
    stringQuotes, _ := json.Marshal(raw["quotes"])
    var rawQuotes map[string]interface{}
    json.Unmarshal(stringQuotes, &rawQuotes)
    var currencyString string = "USD" + target
    var conversionRateBytes, _ = json.Marshal(rawQuotes[currencyString])
    var converstionRateString = string(conversionRateBytes)
    var conversionRateFloat, _ = strconv.ParseFloat(converstionRateString, 64)
    result.ConversionRate = conversionRateFloat
  }

  return result
}
