package get_tweet

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	retrieve_mock "github.com/uala-challenge/tweets-service/internal/retrieve_tweet/mock"
	"github.com/uala-challenge/tweets-service/kit"
)

func TestGetTweet_Success(t *testing.T) {
	mockRetrieveTweet := retrieve_mock.NewService(t)

	mockTweet := kit.Tweet{
		TweetID: "tweet-123",
		UserID:  "user-456",
		Content: "Este es un tweet de prueba",
	}

	mockRetrieveTweet.On("Apply", mock.Anything, mock.MatchedBy(func(pk map[string]interface{}) bool {
		expectedPK := map[string]interface{}{"SK": "user-456", "PK": "tweet-123"}
		return assert.ObjectsAreEqual(expectedPK, pk)
	})).Return(mockTweet, nil)

	service := NewService(Dependencies{
		UseCaseRetrieveTweet: mockRetrieveTweet,
	})

	req := httptest.NewRequest("GET", "/tweet/tweet-123/user-456", nil)
	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("user_id", "user-456")
	reqCtx.URLParams.Add("tweet_id", "tweet-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqCtx)) // Asegurar que los parámetros sean accesibles

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var responseTweet kit.Tweet
	err := json.Unmarshal(rr.Body.Bytes(), &responseTweet)
	assert.NoError(t, err)
	assert.Equal(t, mockTweet, responseTweet)

	mockRetrieveTweet.AssertExpectations(t)
}

func TestGetTweet_TweetNotFound(t *testing.T) {
	mockRetrieveTweet := retrieve_mock.NewService(t)

	mockRetrieveTweet.On("Apply", mock.Anything, mock.Anything).Return(kit.Tweet{}, errors.New("tweet not found"))

	service := NewService(Dependencies{
		UseCaseRetrieveTweet: mockRetrieveTweet,
	})

	req := httptest.NewRequest("GET", "/tweet/tweet-123/user-456", nil)
	reqCtx := chi.NewRouteContext()
	reqCtx.URLParams.Add("user_id", "user-456")
	reqCtx.URLParams.Add("tweet_id", "tweet-123")
	req = req.WithContext(context.Background())

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error to get tweet")

	mockRetrieveTweet.AssertExpectations(t)
}
