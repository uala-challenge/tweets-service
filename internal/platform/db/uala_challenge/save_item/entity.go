package save_item

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge"
)

type Service interface {
	Accept(ctx context.Context, itm map[string]interface{}) error
}

type Dependencies struct {
	Client *dynamodb.Client
	Config uala_challenge.Config
	Log    log.Service
}
