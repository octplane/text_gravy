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

func infoURL(c *Client, args map[string]string) string {
	argsCopy := clone(args)
	argsCopy["extras"] += ",url_t"
	return makeURL(c, "flickr.photos.getinfo", argsCopy, true)
}

// Searches for photos.  args contains search parameters as described in
// http://www.flickr.com/services/api/flickr.photos.search.html.
func (c *Client) Info(photoId string) (*SearchResponse, error) {
	r := struct {
		Stat   string         `xml:"stat,attr"`
		Err    flickrError    `xml:"err"`
		Photos SearchResponse `xml:"photos"`
	}{}
	if err := flickrGet(c, searchURL(c, args), &r); err != nil {
		return nil, err
	}
	if r.Stat != "ok" {
		return nil, r.Err.Err()
	}

	for i, ph := range r.Photos.Photos {
		h, hErr := strconv.ParseFloat(ph.Height_T, 64)
		w, wErr := strconv.ParseFloat(ph.Width_T, 64)
		if hErr == nil && wErr == nil {
			// ph is apparently just a copy of r.Photos.Photos[i], so we are
			// updating the original.
			r.Photos.Photos[i].Ratio = w / h
		}
	}
	return &r.Photos, nil
}

// Define and parse flags
func init() {
	// https://www.flickr.com/services/apps/by/oct
	flag.StringVar(&apiKey, "api_key", "", "Api Key")
	flag.StringVar(&secret, "secret", "", "Secret Key")
}

type Photo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Thumb string `json:"thumb"`
	Large string `json:"large"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	client = flickgo.New(apiKey, secret, http.DefaultClient)

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")

	w.Header().Set("Content-Type", "application/json")

	client = flickgo.New(apiKey, secret, http.DefaultClient)

	fmt.Printf("Searching for %s\n", query)

	sparams := make(map[string]string)
	sparams["text"] = query
	sparams["sort"] = "relevance" //interestingness-desc"
	sparams["licence"] = "4,7"
	//	sparams["others"] = "url_n,tags"

	var err error
	var sresponse *flickgo.SearchResponse

	sresponse, err = client.Search(sparams)
	if err != nil {
		panic(err)
	}

	var resCount int
	var tempPhoto Photo

	resCount, err = strconv.Atoi(sresponse.Total)
	if err != nil {
		panic(err)
	}

	max := 50

	photos := make([]Photo, 0, max)

	for i := 0; i < int(math.Min(float64(resCount), float64(max))); i++ {

		tempPhoto = Photo{}
		tempPhoto.Id = sresponse.Photos[i].ID + "@flickr"
		tempPhoto.Thumb = sresponse.Photos[i].URL("n")
		tempPhoto.Large = sresponse.Photos[i].URL("z")
		tempPhoto.Title = sresponse.Photos[i].Title
		photos = append(photos, tempPhoto)
	}

	// Serialize the modified kitten to JSON
	j, err := json.Marshal(map[string][]Photo{"photos": photos})
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
	r.HandleFunc("/api/v1/photos", searchHandler).Methods("GET")
	r.HandleFunc("/api/v1/photos/:id", photoHandler).Methods("GET")
	// /api/v1/photos/2251667479@flickr
	http.Handle("/api/", r)

	http.Handle("/", http.FileServer(http.Dir("./public/")))

	log.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
