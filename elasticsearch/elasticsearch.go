package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	elastic "gopkg.in/olivere/elastic.v5"
	"log"
	"os"
	"github.com/hectorgool/api-rest-elasticsearch-gin/common"
	//"github.com/satori/go.uuid"
)

type ID struct {
	_Id string `json:"_id"` 
}

type Document struct {
	Id         string `json:"id"`
	Ciudad     string `json:"ciudad"`
	Colonia    string `json:"colonia"`
	Cp         string  `json:"cp"`
	Delegacion string `json:"delegacion"`
	Location   `json:"location"`
}

type Location struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type User struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	Nickname string `json:"nickname"`
}

var client *elastic.Client

func init() {

	var err error

	client, err = elastic.NewClient(
		elastic.SetURL(os.Getenv("ELASTICSEARCH_ENTRYPOINT")),
		elastic.SetBasicAuth(os.Getenv("ELASTICSEARCH_USERNAME"), os.Getenv("ELASTICSEARCH_PASSWORD")),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	common.CheckError(err)

}

func Ping() (string, error) {

	ctx := context.Background()
	info, code, err := client.Ping(os.Getenv("ELASTICSEARCH_ENTRYPOINT")).Do(ctx)
	common.CheckError(err)

	msg := fmt.Sprintf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
	return msg, nil

}

func TermToJson(term string) (string, error) {

	if len(term) == 0 {
		return "", errors.New("No string supplied")
	}
	searchJson := fmt.Sprintf(
	`{
	   "query": {
	     	"match": {
	        	"_all": {
	            	"operator": "and",
	            	"query": "%v"
	         	}
	      	}
	   },
	   "size": 10,
	   "sort": [
	      	{
	        	"colonia": {
	            	"order": "asc"
	         	}
	      	}
	   ]
	}`, term )

	return searchJson, nil

}

func SearchTerm(term string) (*elastic.SearchResult, error) {

	ctx := context.Background()
	if len(term) == 0 {
		return nil, errors.New("No string supplied")
	}

	//Convert string to json query for elasticsearch
	searchJson, err := TermToJson(term)
	common.CheckError(err)

	// Search with a term source
	searchResult, err := client.Search().
		Index(os.Getenv("ELASTICSEARCH_INDEX")).
		Type(os.Getenv("ELASTICSEARCH_TYPE")).
		Source(searchJson).
		Do(ctx)
	common.CheckError(err)

	return searchResult, nil

}

func DisplayResults( searchResult *elastic.SearchResult ) ([]*Document, error) {

    var Documents []*Document

    for _, hit := range searchResult.Hits.Hits {
        d := &Document{}
        //parses *hit.Source into the instance of the Document struct
       	err := json.Unmarshal(*hit.Source, &d)
        common.CheckError(err)
        //Puts d into a map for later access
        Documents = append(Documents, d)
    }
    return Documents, nil

}

func Search(term string) ([]*Document, error) {

	searchResult, err := SearchTerm(term)
	common.CheckError(err)

	result, err := DisplayResults(searchResult)
	common.CheckError(err)

	return result, nil

}

func DeleteDocument(id string) {

	ctx := context.Background()
	res, err := client.Delete().
		Index(os.Getenv("ELASTICSEARCH_INDEX")).
		Type(os.Getenv("ELASTICSEARCH_TYPE")).
	    Id(id).
	    Do(ctx)
	common.CheckError(err)

	if res.Found {
	    fmt.Print("Document deleted from from index\n")
	}

}

func ReadDocument(id string) {

	ctx := context.Background()
	get, err := client.Get().
		Index(os.Getenv("ELASTICSEARCH_INDEX")).
		Type(os.Getenv("ELASTICSEARCH_TYPE")).
	    Id(id).
	    Do(ctx)
	common.CheckError(err)

	if get.Found {
	    fmt.Printf("Got document %s in version %d from index %s, type %s\n", get.Id, get.Version, get.Index, get.Type)
	}

}

/*
func CreateDocument(d Document) {

	id := uuid.Must(uuid.NewV4())
	doc := Document{
		"Ciudad" : d.ciudad,
		"Colonia" : d.colonia,
		"Cp" : d.cp,
		"Delegacion" : d.delegacion,
		Location{
			"Lat" : d.lat,
			"Lon" : d.lon,
		},
	}
	ctx := context.Background()
	get, err := client.Index().
		Index(os.Getenv("ELASTICSEARCH_INDEX")).
		Type(os.Getenv("ELASTICSEARCH_TYPE")).
	    Id(id).
	    BodyString(doc).
	    Do(ctx)
	common.CheckError(err)
	if err != nil {
		// Handle error
		panic(err)
	}
	return doc

}
*/

