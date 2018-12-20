package main

import (
	 "context"
	"log"
	// "encoding/json"
	"fmt"
	"net/http"
	"math"
	"reflect"
	"strconv"
	"time"
	elastic "gopkg.in/olivere/elastic.v5"
	"github.com/speps/go-hashids"
	 "github.com/gin-gonic/gin"
	 "github.com/gin-contrib/cors"

)


var prefix string = "http://localhost:8080/"
// var ctx context
var id_count int = 1


type ShortUrl struct {
	Hash     string                `json:"hash"`
	Original  string                `json:"original_url"`
	Shortened  string                `json:"shortened_url"`
}
const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"entry":{
			"properties":{
				"hash":{
					"type":"keyword"
				},
				"original_url":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"short_url":{
					"type":"text",
					"store": true,
					"fielddata": true
				}
			}
		}
	}
}`
func collison_check(c *gin.Context, hash string) int{
	client, err := elastic.NewClient()
	if err != nil {
	// 	// Handle error
		panic(err)
	}
	termQuery := elastic.NewTermQuery("hash", hash)
	searchResult, err := client.Search().
				Index("shorturl").
				Type("entry"). // search in type
				Query(termQuery).
				From(0). // Starting from this result
				Size(5).  // Limit of responds
				Do(c.Request.Context())
	return int(searchResult.Hits.TotalHits)
}
func createHandle(c *gin.Context) {

client, err := elastic.NewClient()
if err != nil {
// 	// Handle error
	panic(err)
}
term:= c.PostForm("url")
fmt.Print(term)
termQuery := elastic.NewTermQuery("original_url", term)
searchResult, err := client.Search().
			Index("shorturl").
			Type("entry"). // search in type
			Query(termQuery).
			From(0). // Starting from this result
			Size(5).  // Limit of responds
			Do(c.Request.Context())
			if err != nil {
				// Handle error
				panic(err)
			}
fmt.Print(searchResult.Hits.TotalHits)
fmt.Print(c.PostForm("url"))
if (searchResult.Hits.TotalHits == 0)  {

	hd := hashids.NewData()
	// hd.Salt = "this is my salt"
	h,_ := hashids.NewWithData(hd)
	now := time.Now()
	fmt.Print(now)
	bod := ShortUrl{Hash: "Kunjan", Original:c.PostForm("url") , Shortened: "TF"}

	bod.Hash,_ = h.Encode([]int{int(math.Mod(float64(now.Unix()),10000))})
	index := 10000
	var loop_cond  int = collison_check(c,bod.Hash)
	for loop_cond>0{
		bod.Hash,_ = h.Encode([]int{int(math.Mod(float64(now.Unix()),float64(index)))})
		index = index * 10
		loop_cond  = collison_check(c,bod.Hash)
	}
	bod.Shortened = prefix + bod.Hash
	fmt.Print(bod)
	st_id := strconv.Itoa(id_count)
	id_count += 1

	put1, err := client.Index().
		Index("shorturl").
		Type("entry").
		Id(st_id).
		BodyJson(bod).
		Do(c.Request.Context())
		if err != nil {
			// Handle error
			panic(err)
		}
		fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
		c.Status(http.StatusOK)
		return
}
c.Status(208)
}
func redirectHandle(c *gin.Context){
	term := c.Param("hash")
	ci,err := elastic.NewClient()
	if err != nil {
	// 	// Handle error
		panic(err)
	}
	fmt.Print(term)
	termQuery := elastic.NewTermQuery("hash", term)
	searchResult, err := ci.Search().
				Index("shorturl").
				Type("entry"). // search in type
        Query(termQuery).
        From(0). // Starting from this result
        Size(5).  // Limit of responds
        Do(c.Request.Context())         // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	var ttyp ShortUrl
	var msg string
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(ShortUrl); ok {
			fmt.Printf("Tweet by %s: %s\n", t.Original, t.Shortened)
			msg = t.Original
		}
	}
	c.JSON(http.StatusOK,gin.H{"url":msg})


}
func prettyHandle(c *gin.Context){
	term := c.Param("orig")
	termQuery := elastic.NewTermQuery("original_url", term)
	ci,err := elastic.NewClient()
	if err != nil {
	// 	// Handle error
		panic(err)
	}
	searchResult, err := ci.Search().
				Index("shorturl").
				Type("entry"). // search in type
        Query(termQuery).
        From(0). // Starting from this result
        Size(5).  // Limit of responds
        Do(c.Request.Context())         // execute
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	var ttyp ShortUrl
	var msg string
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(ShortUrl); ok {
			fmt.Printf("Tweet by %s: %s\n", t.Original, t.Shortened)
			msg = t.Shortened
		}
	}
	c.JSON(http.StatusOK,gin.H{"url":msg})
}
func main() {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()
	r := gin.Default()
	// // Obtain a client and connect to the default Elasticsearch installation
	// // on 127.0.0.1:9200. Of course you can configure your client to connect
	// // to other hosts and configure it in various other ways.
	client, err := elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("shorturl").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("shorturl").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	query := elastic.NewMatchAllQuery()
	searchResult, err := client.Search().
				Index("shorturl").
				Type("entry"). // search in type
        Query(query).
				Do(ctx)
				if err != nil {
					// Handle error
					panic(err)
				}
	id_count = int(searchResult.Hits.TotalHits) + 1
	fmt.Print(id_count)
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))
	r.POST("/create", createHandle)
	r.GET("/redirect/:hash",redirectHandle)
	r.GET("/pretty/:orig",prettyHandle)
	log.Fatal(r.Run(":8000"))
	// log.Fatal(http.ListenAndServe(":8080", r))
}
