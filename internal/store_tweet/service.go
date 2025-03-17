package store_tweet

import (
	"context"
	"github.com/uala-challenge/tweets-service/internal/platform/publish_tweet_event_sns"

	"github.com/google/uuid"
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/save_item"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweets-service/internal/store_tweet/mappers"
	"github.com/uala-challenge/tweets-service/kit"
)

type service struct {
	sns    publish_tweet_event_sns.Service
	db     save_item.Service
	log    log.Service
	config Config
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		sns:    d.SNSRepository,
		db:     d.DBRepository,
		log:    d.Log,
		config: d.Config,
	}
}

func (s service) Apply(ctx context.Context, rqt kit.TweetRequest) (string, error) {
	tweetID := uuid.New().String()
	tweet := mappers.TweetDynamoMap(rqt, tweetID)
	tweetDynamo, err := kit.StructToMap[kit.DynamoItem](tweet)
	if err != nil {
		return "", s.log.WrapError(err, "Error al mapear tweet")
	}
	err = s.db.Accept(ctx, tweetDynamo, s.config.Table)
	if err != nil {
		return "", s.log.WrapError(err, "Error al guardar tweet")
	}
	tweetSNS := mappers.TweetSNSMap(tweet)
	err = s.sns.Accept(ctx, mappers.TweetToPublishInput(tweetSNS, s.config.Topic), s.config.Retries)
	if err != nil {
		return "", s.log.WrapError(err, "Error al publicar tweet")
	}
	return tweetID, nil
}
