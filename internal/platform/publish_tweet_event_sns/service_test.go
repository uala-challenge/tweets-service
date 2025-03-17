package publish_tweet_event_sns

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	sns2 "github.com/uala-challenge/simple-toolkit/pkg/client/sns"
	sns_mock "github.com/uala-challenge/simple-toolkit/pkg/client/sns/mock"
	log_mock "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
)

func TestPublishMessage_Success(t *testing.T) {
	cli := sns_mock.NewService(t)
	l := log_mock.NewService(t)

	pubInput := &sns.PublishInput{
		Message:  new(string),
		TopicArn: new(string),
	}

	cli.On("Publish", mock.Anything, pubInput).Return(&sns.PublishOutput{}, nil)

	l.On("Info", mock.Anything, "Mensaje publicado en SNSRepository", mock.Anything).Return()

	service := NewService(Dependencies{
		Client: &sns2.Sns{Cliente: cli},
		Log:    l,
	})

	err := service.Accept(context.TODO(), pubInput, 3)

	assert.NoError(t, err)
	cli.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestPublishMessage_FailAllRetries(t *testing.T) {
	cli := sns_mock.NewService(t)
	l := log_mock.NewService(t)

	pubInput := &sns.PublishInput{
		Message:  new(string),
		TopicArn: new(string),
	}

	expectedErr := errors.New("SNS Service Unavailable")
	cli.On("Publish", mock.Anything, pubInput).Return(nil, expectedErr).Times(3)

	for i := 1; i <= 3; i++ {
		l.On("Warn", mock.Anything, "Reintentando publicación de mensaje en SNSRepository...", mock.MatchedBy(func(m map[string]interface{}) bool {
			return m["attempt"] == i && m["error"] == "SNS Service Unavailable"
		})).Once()
	}

	l.On("Error", mock.Anything, expectedErr, "No se pudo publicar el mensaje después de los intentos máximos", mock.Anything).Return()

	service := NewService(Dependencies{
		Client: &sns2.Sns{Cliente: cli},
		Log:    l,
	})

	err := service.Accept(context.TODO(), pubInput, 3)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	cli.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestPublishMessage_SuccessAfterRetry(t *testing.T) {
	cli := sns_mock.NewService(t)
	l := log_mock.NewService(t)

	pubInput := &sns.PublishInput{
		Message:  new(string),
		TopicArn: new(string),
	}

	expectedErr := errors.New("SNS Service Unavailable")

	cli.On("Publish", mock.Anything, pubInput).Return(nil, expectedErr).Once()
	cli.On("Publish", mock.Anything, pubInput).Return(&sns.PublishOutput{}, nil).Once()

	l.On("Warn", mock.Anything, "Reintentando publicación de mensaje en SNSRepository...", mock.Anything).Once()

	l.On("Info", mock.Anything, "Mensaje publicado en SNSRepository", mock.Anything).Return()

	service := NewService(Dependencies{
		Client: &sns2.Sns{Cliente: cli},
		Log:    l,
	})

	err := service.Accept(context.TODO(), pubInput, 3)

	assert.NoError(t, err)
	cli.AssertExpectations(t)
	l.AssertExpectations(t)
}
