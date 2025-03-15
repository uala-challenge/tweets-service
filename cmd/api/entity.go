package main

import (
	"github.com/uala-challenge/tweet-service/cmd/api/post_tweet"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/get_item"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/save_item"
	"github.com/uala-challenge/tweet-service/internal/platform/sns/publish_tweet_event_sns"
	"github.com/uala-challenge/tweet-service/internal/store_tweet"
)

type repositories struct {
	PublishTweet publish_tweet_event_sns.Service
	SaveTweet    save_item.Service
	GetTweet     get_item.Service
	GetTweets    get_item.Service
}

type useCases struct {
	StoreTweet store_tweet.Service
}

type handlers struct {
	Test post_tweet.Service
}
