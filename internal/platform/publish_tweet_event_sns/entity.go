package publish_tweet_event_sns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	sns2 "github.com/uala-challenge/simple-toolkit/pkg/client/sns"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	Accept(ctx context.Context, pubInput *sns.PublishInput, retries int) error
}

type Dependencies struct {
	Client *sns2.Sns
	Log    log.Service
}
