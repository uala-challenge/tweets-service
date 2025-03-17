package store_tweet

import (
	"context"
	"github.com/uala-challenge/tweets-service/internal/platform/publish_tweet_event_sns"

	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/save_item"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweets-service/kit"
)

type Service interface {
	Apply(ctx context.Context, rqt kit.TweetRequest) (string, error)
}

type Dependencies struct {
	SNSRepository publish_tweet_event_sns.Service
	DBRepository  save_item.Service
	Log           log.Service
	Config        Config
}

type Config struct {
	Topic   string `json:"topic"`
	Retries int    `json:"retries"`
	Table   string `json:"table"`
}
