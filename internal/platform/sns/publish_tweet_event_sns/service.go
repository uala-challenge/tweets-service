package publish_tweet_event_sns

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type service struct {
	client *sns.Client
	log    log.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) *service {
	return &service{
		client: d.Client,
		log:    d.Log,
	}
}

func (s *service) Accept(ctx context.Context, pubInput *sns.PublishInput, retries int) error {
	retryDelays := make([]time.Duration, retries)
	for i := 0; i < retries; i++ {
		retryDelays[i] = time.Duration(1<<i) * time.Second
	}
	var lastErr error
	for attempt := 1; attempt <= retries; attempt++ {
		_, err := s.client.Publish(ctx, pubInput)
		if err == nil {
			s.log.Info(ctx, "Mensaje publicado en SNSRepository", nil)
			return nil
		}

		lastErr = err
		s.log.Warn(ctx, "Reintentando publicación de mensaje en SNSRepository...", map[string]interface{}{
			"attempt": attempt,
			"error":   err.Error(),
		})

		if attempt < retries {
			time.Sleep(retryDelays[attempt-1])
		}
	}

	s.log.Error(ctx, lastErr, "No se pudo publicar el mensaje después de los intentos máximos", nil)
	return lastErr
}
