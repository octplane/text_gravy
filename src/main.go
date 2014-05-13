package main

import "code.google.com/p/flickgo"
import "net/http"
import "fmt"
import "flag"

func main() {
  var err error
  var sresponse * flickgo.SearchResponse

  // Get the API keys from https://www.flickr.com/services/apps/by/oct
  var api_key string
  var secret  string


  flag.StringVar(&api_key, "api_key", "1234", "Api Key")
  flag.StringVar(&secret, "secret", "12", "Secret Key")
  flag.Parse()

  fmt.Printf("%s %s", api_key, secret)


  client := flickgo.New( api_key, secret, http.DefaultClient )


  sparams := make(map[string]string)
  sparams["text"] = "meditation"
  sparams["sort"] = "interestingness-desc"
  sparams["licence"] = "4,7"


  sresponse, err = client.Search(sparams)
  if err != nil {
    panic(err)
  }
  fmt.Println("%s", sresponse.Photos[0].URL("n"))
}
