package retrieve_tweet

import (
	"context"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/get_item"
	"github.com/uala-challenge/tweet-service/kit"
)

type Service interface {
	Apply(ctx context.Context, pk map[string]interface{}) (kit.Tweet, error)
}

type Dependencies struct {
	DBRepository get_item.Service
	Log          log.Service
}
