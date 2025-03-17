package retrieve_tweet

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	gt "github.com/uala-challenge/simple-toolkit/pkg/platform/db/get_item/mock"
	log_mock "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
	"github.com/uala-challenge/tweets-service/kit"
)

func getMockDynamoItem() map[string]types.AttributeValue {
	item := kit.DynamoItem{
		PK:      "tweet-123",
		SK:      "user-456",
		GSI1PK:  "user-456",
		GSI1SK:  "tweet-123",
		Content: "Este es un tweet de prueba",
		Created: time.Now().Unix(),
	}
	marshaled, _ := attributevalue.MarshalMap(item)
	return marshaled
}

func TestRetrieveTweet_Success(t *testing.T) {
	mockDB := gt.NewService(t)
	mockLog := log_mock.NewService(t)

	mockPK := map[string]interface{}{
		"PK": "tweet-123",
	}

	mockDynamoItem := getMockDynamoItem()

	mockDB.On("Apply", mock.Anything, mockPK, "tweets_table").Return(mockDynamoItem, nil)

	service := NewService(Dependencies{
		DBRepository: mockDB,
		Log:          mockLog,
		Config:       Config{Table: "tweets_table"},
	})

	tweet, err := service.Apply(context.TODO(), mockPK)

	assert.NoError(t, err)
	assert.Equal(t, "tweet-123", tweet.TweetID)
	assert.Equal(t, "user-456", tweet.UserID)
	assert.Equal(t, "Este es un tweet de prueba", tweet.Content)
}

func TestRetrieveTweet_DBError(t *testing.T) {
	mockDB := gt.NewService(t)
	mockLog := log_mock.NewService(t)

	mockPK := map[string]interface{}{
		"PK": "tweet-123",
	}

	expectedErr := errors.New("error en DynamoDB")

	mockDB.On("Apply", mock.Anything, mockPK, "tweets_table").Return(nil, expectedErr)
	mockLog.On("WrapError", expectedErr, "error al obtener tweet").Return(expectedErr)

	service := NewService(Dependencies{
		DBRepository: mockDB,
		Log:          mockLog,
		Config:       Config{Table: "tweets_table"},
	})

	tweet, err := service.Apply(context.TODO(), mockPK)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, kit.Tweet{}, tweet)
	mockDB.AssertExpectations(t)
	mockLog.AssertExpectations(t)
}
