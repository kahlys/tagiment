package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/cdipaolo/sentiment"
	twitter "github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func main() {
	token := flag.String("token", "", "twitter API token")
	tag := flag.String("tag", "", "twitter hashtag you want to analyze")
	flag.Parse()

	// init ml model
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		log.Fatal(err)
	}

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	tweets, err := client.TweetRecentSearch(
		context.Background(),
		fmt.Sprintf("#%v lang:en", strings.TrimSpace(*tag)),
		twitter.TweetRecentSearchOpts{
			MaxResults: 100,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// calculate score
	set := map[uint8]int{0: 0, 1: 0}
	for _, t := range tweets.Raw.TweetDictionaries() {
		score := sentimentModel.SentimentAnalysis(t.Tweet.Text, sentiment.English)
		set[score.Score] += 1
	}

	fmt.Println(math.Round(float64(set[1]*100.0) / float64((set[0] + set[1]))))
}
