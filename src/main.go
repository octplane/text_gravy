package main

import "code.google.com/p/flickgo"
import "net/http"
import "fmt"
import "flag"
import "os"
import "log"

var (
	apiKey string
	secret string
)

// Define and parse flags
func init() {
	flag.StringVar(&apiKey, "api_key", "", "Api Key")
	flag.StringVar(&secret, "secret", "", "Secret Key")
}

func main() {
	var err error
	var sresponse *flickgo.SearchResponse

	flag.Parse()

	if apiKey == "" || secret == "" {
		fmt.Println("Missing one or more command line options.")
		flag.PrintDefaults()
		os.Exit(2)
	}

	client := flickgo.New(apiKey, secret, http.DefaultClient)

	sparams := make(map[string]string)
	sparams["text"] = "meditation"
	sparams["sort"] = "interestingness-desc"
	sparams["licence"] = "4,7"

	sresponse, err = client.Search(sparams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", sresponse.Photos[0].URL("n"))

	log.Println("Starting Server")
	http.Handle("/", http.FileServer(http.Dir("./public/")))

	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
