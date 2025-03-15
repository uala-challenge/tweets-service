package store_tweet

import (
	"context"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/save_item"
	"github.com/uala-challenge/tweet-service/internal/platform/sns/publish_tweet_event_sns"
	"github.com/uala-challenge/tweet-service/kit"
)

type Service interface {
	Apply(ctx context.Context, rqt kit.TweetRequest) (string, error)
}

type Dependencies struct {
	SNS    publish_tweet_event_sns.Service
	DB     save_item.Service
	Log    log.Service
	Config Config
}

type Config struct {
	Topic   string `json:"topic"`
	Retries int    `json:"retries"`
}
