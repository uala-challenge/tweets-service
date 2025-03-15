package get_item

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge"
	"github.com/uala-challenge/tweet-service/kit"
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

func (s service) Apply(ctx context.Context, item map[string]interface{}) (*kit.DynamoItem, error) {
	key, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, s.log.WrapError(err, "error serializando clave de búsqueda")
	}

	result, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &s.config.Table,
		Key:       key,
	})
	if err != nil {
		return nil, s.log.WrapError(err, "error al obtener el item")
	}

	if result.Item == nil {
		return nil, s.log.WrapError(nil, "item no encontrado")
	}

	var itemDB kit.DynamoItem
	err = attributevalue.UnmarshalMap(result.Item, &itemDB)
	if err != nil {
		return nil, s.log.WrapError(nil, "error deserializando item")
	}

	return &itemDB, nil
}
