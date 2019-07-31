package main

import (
    "log"
    "net/http"
    "./currency"
    "strconv"
    "html/template"
    "github.com/bmizerany/pat"
)

func main() {
  mux := pat.New()
  mux.Get("/", http.HandlerFunc(index))
  mux.Post("/", http.HandlerFunc(process))

  log.Println("Listening...")
  http.ListenAndServe(":8080", mux)
}

func index(w http.ResponseWriter, r *http.Request) {
  render(w, "./static/index.html", nil)
}

func process(w http.ResponseWriter, r *http.Request) {
  var valueToConvert, err = strconv.ParseFloat(r.FormValue("number"), 64)
  var target string = r.FormValue("target")

  var conversion currency.Conversion = currency.Convert(target)

  conversion.ValuetoConvert = valueToConvert
  conversion.Target = target

  if target == "" {
    conversion.Errors = append(conversion.Errors, "target currency not selected")
  }
  if err != nil {
    conversion.Errors = append(conversion.Errors, "number entered is not valid")
  }
  conversion.ConvertedValue = conversion.ValuetoConvert * conversion.ConversionRate

  conversion.ErrorString = conversion.ErrorPrettyPrint()
  render(w, "./static/index.html", conversion)
}

func render(w http.ResponseWriter, filename string, data interface{}) {
  tmpl, err := template.ParseFiles(filename)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  if err := tmpl.Execute(w, data); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
