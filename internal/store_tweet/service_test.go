package store_tweet

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	db_mock "github.com/uala-challenge/simple-toolkit/pkg/platform/db/save_item/mock"
	log_mock "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
	sns_mock "github.com/uala-challenge/tweets-service/internal/platform/publish_tweet_event_sns/mock"
	"github.com/uala-challenge/tweets-service/kit"
)

func TestStoreTweet_Success(t *testing.T) {
	mockSNS := sns_mock.NewService(t)
	mockDB := db_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	mockTweetRequest := kit.TweetRequest{
		UserID: "user-456",
		Tweet:  "Este es un tweet de prueba",
	}

	expectedSNSInput := mock.MatchedBy(func(input *sns.PublishInput) bool {
		return input != nil &&
			input.TopicArn != nil && *input.TopicArn == "sns-topic" &&
			input.Message != nil &&
			len(input.MessageAttributes) > 0
	})

	mockDB.On("Accept", mock.Anything, mock.Anything, "tweets_table").Return(nil)
	mockSNS.On("Accept", mock.Anything, expectedSNSInput, 3).Return(nil)

	service := NewService(Dependencies{
		SNSRepository: mockSNS,
		DBRepository:  mockDB,
		Log:           mockLog,
		Config:        Config{Table: "tweets_table", Topic: "sns-topic", Retries: 3},
	})

	tweetID, err := service.Apply(context.TODO(), mockTweetRequest)

	assert.NoError(t, err)
	assert.NotEmpty(t, tweetID)
	mockDB.AssertExpectations(t)
	mockSNS.AssertExpectations(t)
}

func TestStoreTweet_FailToSaveInDynamoDB(t *testing.T) {
	mockSNS := sns_mock.NewService(t)
	mockDB := db_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	mockTweetRequest := kit.TweetRequest{
		UserID: "user-456",
		Tweet:  "Este es un tweet de prueba",
	}

	expectedErr := errors.New("error guardando en DynamoDB")

	mockDB.On("Accept", mock.Anything, mock.Anything, "tweets_table").Return(expectedErr)
	mockLog.On("WrapError", expectedErr, "Error al guardar tweet").Return(expectedErr)

	service := NewService(Dependencies{
		SNSRepository: mockSNS,
		DBRepository:  mockDB,
		Log:           mockLog,
		Config:        Config{Table: "tweets_table", Topic: "sns-topic", Retries: 3},
	})

	tweetID, err := service.Apply(context.TODO(), mockTweetRequest)

	assert.Error(t, err)
	assert.Empty(t, tweetID)
	mockDB.AssertExpectations(t)
	mockSNS.AssertNotCalled(t, "Accept")
}

func TestStoreTweet_FailToPublishToSNS(t *testing.T) {
	mockSNS := sns_mock.NewService(t)
	mockDB := db_mock.NewService(t)
	mockLog := log_mock.NewService(t)

	mockTweetRequest := kit.TweetRequest{
		UserID: "user-456",
		Tweet:  "Este es un tweet de prueba",
	}

	expectedSNSInput := mock.MatchedBy(func(input *sns.PublishInput) bool {
		return input != nil &&
			input.TopicArn != nil && *input.TopicArn == "sns-topic" &&
			input.Message != nil &&
			len(input.MessageAttributes) > 0
	})

	expectedErr := errors.New("error publicando en SNS")

	mockDB.On("Accept", mock.Anything, mock.Anything, "tweets_table").Return(nil)
	mockSNS.On("Accept", mock.Anything, expectedSNSInput, 3).Return(expectedErr)
	mockLog.On("WrapError", expectedErr, "Error al publicar tweet").Return(expectedErr)

	service := NewService(Dependencies{
		SNSRepository: mockSNS,
		DBRepository:  mockDB,
		Log:           mockLog,
		Config:        Config{Table: "tweets_table", Topic: "sns-topic", Retries: 3},
	})

	tweetID, err := service.Apply(context.TODO(), mockTweetRequest)

	assert.Error(t, err)
	assert.Empty(t, tweetID)
	mockDB.AssertExpectations(t)
	mockSNS.AssertExpectations(t)
}
