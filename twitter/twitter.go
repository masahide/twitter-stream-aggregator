package twitter

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/oauth1"
)

const (
	streamingEndpoint = `https://stream.twitter.com/1.1/statuses/filter.json`
	apiKeyFileName    = `apikeys.json`
)

// TwitterConfig twitter config file struct
type twitterConfig struct {
	Config *oauth1.Config
	Token  *oauth1.Token
}

func loadConfig() (conf twitterConfig, err error) {
	jsonString, err := ioutil.ReadFile(apiKeyFileName)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonString, &conf)
	return
}

// Twitter twitter stream api filter
type Twitter struct {
	client *http.Client
}

// NewTwitter twitter stream
func NewTwitter(ctx context.Context) (*Twitter, error) {
	res := Twitter{}
	conf, err := loadConfig()
	if err != nil {
		return nil, err
	}
	res.client = conf.Config.Client(ctx, conf.Token)
	return &res, nil
}

// StreamFilter  twitter stream filter
func (ts *Twitter) StreamFilter(queryString string) (*http.Response, error) {
	req, err := http.NewRequest("POST", streamingEndpoint+queryString, nil)
	if err != nil {
		return nil, err
	}
	return ts.client.Do(req)
}
