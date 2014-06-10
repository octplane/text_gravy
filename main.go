package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "github.com/gorilla/mux"
  "github.com/octplane/flickgo"
  "log"
  "math"
  "net/http"
  "os"
  "strconv"
  "strings"
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

type Photo struct {
  Id    string `json:"id"`
  Title string `json:"title"`
  Thumb string `json:"thumb"` // n
  Large string `json:"large"` // z
}

type PhotoInfo struct {
  ID          string        `json:"id"`
  Rotation    string        `json:"rotation,attr"`
  License     string        `json:"license,attr"`
  Title       string        `json:"title"`
  Description string        `json:"description"`
  Dates       flickgo.Dates `xml:"dates"`
  Urls        []flickgo.Url `xml:"urls>url"`
  Thumb       string        `json:"thumb"`
  Large       string        `json:"large"`
  TagsRef     []string      `json:"tags"`
}

type Tag struct {
  ID        string `json:"id"`
  Text      string `json:"text"`
  Photoinfo string `json:"photoinfo_id"`
}

func TagsToTags(tags []flickgo.Tag, photoinfo string) []Tag {
  ret := make([]Tag, 0, 20)

  for _, tag := range tags {
    ids := strings.Split(tag.ID, "-")
    ret = append(ret, Tag{
      ID:        ids[len(ids)-1],
      Text:      tag.Text,
      Photoinfo: photoinfo})
  }
  return ret
}

func TagsToId(tags []Tag) []string {
  ret := make([]string, 0, 20)

  for _, tag := range tags {
    ret = append(ret, tag.ID)
  }
  return ret
}

func photoHandler(w http.ResponseWriter, r *http.Request) {
  client := flickgo.New(apiKey, secret, http.DefaultClient)
  photo_id := mux.Vars(r)["photo_id"]

  finfo, err := client.GetInfo(photo_id)
  if err != nil {
    panic(err)
  }

  tags := TagsToTags(finfo.Tags, finfo.ID)

  info := PhotoInfo{
    ID:          finfo.ID,
    Rotation:    finfo.Rotation,
    License:     finfo.License,
    Title:       finfo.Title,
    Description: finfo.Description,
    Dates:       finfo.Dates,
    Urls:        finfo.Urls,
    Thumb:       finfo.Photo.URL("n"),
    Large:       finfo.Photo.URL("z"),
    TagsRef:     TagsToId(tags),
  }

  // Serialize the modified kitten to JSON
  j, err := json.MarshalIndent(map[string]interface{}{
    "tags":      tags,
    "photoinfo": &info,
  }, "", " ")
  if err != nil {
    panic(err)
  }

  w.Write(j)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
  photo_id := mux.Vars(r)["photo_id"]
  query := r.FormValue("query")
  fmt.Printf("Searching for %s\n", query)

  client = flickgo.New(apiKey, secret, http.DefaultClient)

  max := 50
  photos := make([]Photo, 0, max)

  if photo_id != "" {
    finfo, err := client.GetInfo(photo_id)
    if err != nil {
      panic(err)
    }

    info := Photo{
      Id:    finfo.ID,
      Title: finfo.Title,
      Thumb: finfo.Photo.URL("n"),
      Large: finfo.Photo.URL("z"),
    }
    photos = append(photos, info)

  } else {
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

    for i := 0; i < int(math.Min(float64(resCount), float64(max))); i++ {

      tempPhoto = Photo{Id: sresponse.Photos[i].ID,
        Thumb: sresponse.Photos[i].URL("n"),
        Large: sresponse.Photos[i].URL("z"),
        Title: sresponse.Photos[i].Title,
      }
      photos = append(photos, tempPhoto)
    }
  }

  // Serialize the modified kitten to JSON
  j, err := json.MarshalIndent(map[string][]Photo{"photos": photos}, "", " ")
  if err != nil {
    panic(err)
  }
  w.Header().Set("Content-Type", "application/json")
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
  r.HandleFunc("/api/v1/photos/{photo_id:[a-z0-9]+}", searchHandler).Methods("GET")
  r.HandleFunc("/api/v1/photoinfos/{photo_id:[a-z0-9]+}", photoHandler).Methods("GET")
  r.HandleFunc("/api/v1/photos", searchHandler).Methods("GET")

  http.Handle("/api/", r)
  http.Handle("/", http.FileServer(http.Dir("./public/")))

  port := os.Getenv("PORT")

  if port == "" {
    port = "8080"
  }
  log.Printf("Listening on %s\n", port)
  http.ListenAndServe(":"+port, nil)
}
