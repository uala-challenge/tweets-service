package batch_get_tweets

import (
	"context"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/get_items"
	"github.com/uala-challenge/tweet-service/kit"
)

type Service interface {
	Apply(ctx context.Context, tweets []kit.TweetPK) ([]kit.Tweet, error)
}

type Dependencies struct {
	DBRepository get_items.Service
	Log          log.Service
}
