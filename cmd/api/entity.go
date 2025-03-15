package main

import (
	"github.com/uala-challenge/tweet-service/cmd/api/get_tweet"
	"github.com/uala-challenge/tweet-service/cmd/api/list_tweets"
	"github.com/uala-challenge/tweet-service/cmd/api/post_tweet"
	"github.com/uala-challenge/tweet-service/internal/batch_get_tweets"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/get_item"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/get_items"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/save_item"
	"github.com/uala-challenge/tweet-service/internal/platform/sns/publish_tweet_event_sns"
	"github.com/uala-challenge/tweet-service/internal/retrieve_tweet"
	"github.com/uala-challenge/tweet-service/internal/store_tweet"
)

type repositories struct {
	PublishTweet publish_tweet_event_sns.Service
	SaveTweet    save_item.Service
	GetTweet     get_item.Service
	GetTweets    get_items.Service
}

type useCases struct {
	StoreTweet     store_tweet.Service
	BatchGetTweets batch_get_tweets.Service
	RetrieveTweet  retrieve_tweet.Service
}

type handlers struct {
	PostTweet post_tweet.Service
	GetTweet  get_tweet.Service
	ListTweet list_tweets.Service
}
