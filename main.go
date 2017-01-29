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
	apiKeyFileName    = `apikeys.json`
)

type TwitterConfig struct {
	Config *oauth1.Config
	Token  *oauth1.Token
}

func checkFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig() (conf TwitterConfig, err error) {
	jsonString, err := ioutil.ReadFile(apiKeyFileName)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonString, &conf)
	return
}

type TwitterStream struct {
	client *http.Client
}

func NewTwitterStream(ctx context.Context) (*TwitterStream, error) {
	res := TwitterStream{}
	conf, err := loadConfig()
	if err != nil {
		return nil, err
	}
	res.client = conf.Config.Client(ctx, conf.Token)
	return &res, nil
}

func (ts *TwitterStream) filter(queryString string) (*http.Response, error) {
	req, err := http.NewRequest("POST", streamingEndpoint+queryString, nil)
	if err != nil {
		return nil, err
	}
	return ts.client.Do(req)
}

func main() {

	ctx := context.Background()
	ts, err := NewTwitterStream(ctx)
	checkFatal(err)
	resp, err := ts.filter(`?track=pig`)
	checkFatal(err)
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

}
