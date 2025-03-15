package batch_get_tweets

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweet-service/internal/batch_get_tweets/mappers"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/get_items"
	"github.com/uala-challenge/tweet-service/kit"
)

type service struct {
	db  get_items.Service
	log log.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		db:  d.DBRepository,
		log: d.Log,
	}
}

func (s service) Apply(ctx context.Context, tweets []kit.TweetPK) ([]kit.Tweet, error) {
	keys := GenerateTweetKeys(tweets)

	items, err := s.db.Apply(ctx, keys)
	if err != nil {
		return nil, err
	}

	return mappers.DynamoItemsToTweets(items), nil

}

func GenerateTweetKeys(tweets []kit.TweetPK) []map[string]types.AttributeValue {
	var keys []map[string]types.AttributeValue

	for _, tweet := range tweets {
		key := map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: tweet.TweetID},
			"SK": &types.AttributeValueMemberS{Value: tweet.UserID},
		}
		keys = append(keys, key)
	}

	return keys
}
