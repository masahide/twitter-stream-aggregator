package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/oauth1"
)

const (
	streamingEndpoint = `https://stream.twitter.com/1.1/statuses/filter.json`
)

type TwitterConfig struct {
	Config *oauth1.Config
	Token  *oauth1.Token
}

func main() {

	conf := TwitterConfig{}

	jsonString, err := ioutil.ReadFile("apikeys.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(jsonString, &conf)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	httpClient := conf.Config.Client(ctx, conf.Token)

	queryString := `?track=pig`
	req, err := http.NewRequest("POST", streamingEndpoint+queryString, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)

}
