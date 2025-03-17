package retrieve_tweet

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/uala-challenge/tweets-service/internal/retrieve_tweet/mappers"

	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/get_item"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweets-service/kit"
)

type service struct {
	db     get_item.Service
	log    log.Service
	config Config
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		db:     d.DBRepository,
		log:    d.Log,
		config: d.Config,
	}
}

func (s service) Apply(ctx context.Context, pk map[string]interface{}) (kit.Tweet, error) {
	tw, err := s.db.Apply(ctx, pk, s.config.Table)
	if err != nil {
		return kit.Tweet{}, s.log.WrapError(err, "error al obtener tweet")
	}
	var itemDB kit.DynamoItem
	err = attributevalue.UnmarshalMap(tw, &itemDB)
	if err != nil {
		return kit.Tweet{}, s.log.WrapError(nil, "error deserializando item")
	}
	return mappers.DynamoItemToTweet(itemDB), nil
}
