package save_item

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge"
)

type service struct {
	client *dynamodb.Client
	log    log.Service
	config uala_challenge.Config
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) *service {
	return &service{
		client: d.Client,
		log:    d.Log,
		config: d.Config,
	}
}

func (s *service) Accept(ctx context.Context, itm map[string]interface{}) error {
	item, err := attributevalue.MarshalMap(itm)
	if err != nil {
		return s.log.WrapError(err, "Error serializando item")
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &s.config.Table,
		Item:      item,
	})
	if err != nil {
		return s.log.WrapError(err, "Error al guardar el item")
	}
	return nil
}
