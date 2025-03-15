package get_items

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge"
	"github.com/uala-challenge/tweet-service/kit"
)

type Service interface {
	Apply(ctx context.Context, items []map[string]types.AttributeValue) ([]*kit.DynamoItem, error)
}

type Dependencies struct {
	Client *dynamodb.Client
	Config uala_challenge.Config
	Log    log.Service
}
