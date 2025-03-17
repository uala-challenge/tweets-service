package retrieve_tweet

import (
	"context"

	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/get_item"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweets-service/kit"
)

type Service interface {
	Apply(ctx context.Context, pk map[string]interface{}) (kit.Tweet, error)
}

type Dependencies struct {
	DBRepository get_item.Service
	Log          log.Service
	Config       Config
}

type Config struct {
	Table string `json:"table"`
}
