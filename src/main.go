package main

import (
	"code.google.com/p/flickgo"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

var (
	apiKey string
	secret string
	client *flickgo.Client
)

// Define and parse flags
func init() {
	// https://www.flickr.com/services/apps/by/oct
	flag.StringVar(&apiKey, "api_key", "", "Api Key")
	flag.StringVar(&secret, "secret", "", "Secret Key")
}

type Kitten struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type KittenJSON struct {
	Kitten Kitten
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")

	w.Header().Set("Content-Type", "application/json")

	client = flickgo.New(apiKey, secret, http.DefaultClient)

	fmt.Printf("Searching for %s\n", query)

	sparams := make(map[string]string)
	sparams["text"] = query
	sparams["sort"] = "interestingness-desc"
	sparams["licence"] = "4,7"

	var err error
	var sresponse *flickgo.SearchResponse

	sresponse, err = client.Search(sparams)
	if err != nil {
		panic(err)
	}

	var resCount int
	var tempKitty Kitten
	resCount, err = strconv.Atoi(sresponse.Total)
	if err != nil {
		panic(err)
	}

	kittens := make([]Kitten, 0, 10.0)

	for i := 0; i < int(math.Min(float64(resCount), float64(10))); i++ {
		tempKitty = Kitten{}
		tempKitty.Id = sresponse.Photos[i].ID
		tempKitty.Picture = sresponse.Photos[i].URL("n")
		kittens = append(kittens, tempKitty)
	}

	// Serialize the modified kitten to JSON
	j, err := json.Marshal(map[string][]Kitten{"kittens": kittens})
	if err != nil {
		panic(err)
	}

	w.Write(j)
}

func main() {

	flag.Parse()

	if apiKey == "" || secret == "" {
		fmt.Println("Missing one or more command line options.")
		flag.PrintDefaults()
		os.Exit(2)
	}

	log.Println("Starting Server")

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/kittens", searchHandler).Methods("GET")
	http.Handle("/api/", r)

	http.Handle("/", http.FileServer(http.Dir("./public/")))

	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
