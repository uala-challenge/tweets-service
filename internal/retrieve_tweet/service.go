package retrieve_tweet

import (
	"context"

	"github.com/uala-challenge/tweet-service/internal/retrieve_tweet/mappers"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/get_item"
	"github.com/uala-challenge/tweet-service/kit"
)

type service struct {
	db  get_item.Service
	log log.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		db:  d.DBRepository,
		log: d.Log,
	}
}

func (s service) Apply(ctx context.Context, pk map[string]interface{}) (kit.Tweet, error) {
	tw, err := s.db.Apply(ctx, pk)
	if err != nil {
		return kit.Tweet{}, err
	}
	return mappers.DynamoItemToTweet(*tw), nil
}
