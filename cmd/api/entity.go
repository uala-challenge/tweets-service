package main

import (
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/get_item"
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/save_item"
	"github.com/uala-challenge/tweets-service/cmd/api/get_tweet"
	"github.com/uala-challenge/tweets-service/cmd/api/post_tweet"
	"github.com/uala-challenge/tweets-service/internal/platform/publish_tweet_event_sns"
	"github.com/uala-challenge/tweets-service/internal/retrieve_tweet"
	"github.com/uala-challenge/tweets-service/internal/store_tweet"
)

type repositories struct {
	PublishTweet publish_tweet_event_sns.Service
	SaveTweet    save_item.Service
	GetTweet     get_item.Service
}

type useCases struct {
	StoreTweet    store_tweet.Service
	RetrieveTweet retrieve_tweet.Service
}

type handlers struct {
	PostTweet post_tweet.Service
	GetTweet  get_tweet.Service
}
