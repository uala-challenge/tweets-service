package post_tweet

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	store_mock "github.com/uala-challenge/tweets-service/internal/store_tweet/mock"
	"github.com/uala-challenge/tweets-service/kit"
)

func TestPostTweet_Success(t *testing.T) {
	mockStoreTweet := store_mock.NewService(t)

	tweetReq := kit.TweetRequest{
		UserID: "user-123",
		Tweet:  "Este es un tweet de prueba",
	}
	tweetJSON, _ := json.Marshal(tweetReq)

	mockStoreTweet.On("Apply", mock.Anything, tweetReq).Return("tweet-456", nil)

	service := NewService(Dependencies{
		UseCaseStoreTweet: mockStoreTweet,
	})

	req, err := http.NewRequest("POST", "/tweet", bytes.NewBuffer(tweetJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "tweet creado: tweet-456")
	mockStoreTweet.AssertExpectations(t)
}

func TestPostTweet_EmptyBody(t *testing.T) {
	mockStoreTweet := store_mock.NewService(t)

	service := NewService(Dependencies{
		UseCaseStoreTweet: mockStoreTweet,
	})

	req, err := http.NewRequest("POST", "/tweet", bytes.NewBuffer([]byte{}))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "validate error")
	assert.Contains(t, rr.Body.String(), "UserID es requerido")

	mockStoreTweet.AssertNotCalled(t, "Apply")
}

func TestPostTweet_ValidationError(t *testing.T) {
	mockStoreTweet := store_mock.NewService(t)

	tweetReq := kit.TweetRequest{
		UserID: "user-123",
		Tweet:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent consequat tortor ac neque vestibulum, nec sagittis orci vestibulum. Aliquam erat volutpat. Sed at volutpat purus. Fusce feugiat libero et erat varius, eget tincidunt nulla volutpat. Aenean maximus justo non eros posuere, id lobortis sapien interdum.", // Más de 280 caracteres
	}
	tweetJSON, _ := json.Marshal(tweetReq)

	service := NewService(Dependencies{
		UseCaseStoreTweet: mockStoreTweet,
	})

	req, err := http.NewRequest("POST", "/tweet", bytes.NewBuffer(tweetJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "validate error")
	mockStoreTweet.AssertNotCalled(t, "Apply")
}

func TestPostTweet_StoreError(t *testing.T) {
	mockStoreTweet := store_mock.NewService(t)

	tweetReq := kit.TweetRequest{
		UserID: "user-123",
		Tweet:  "Este es un tweet de prueba",
	}
	tweetJSON, _ := json.Marshal(tweetReq)

	expectedErr := errors.New("error al procesar tweet")

	mockStoreTweet.On("Apply", mock.Anything, tweetReq).Return("", expectedErr)

	service := NewService(Dependencies{
		UseCaseStoreTweet: mockStoreTweet,
	})

	req, err := http.NewRequest("POST", "/tweet", bytes.NewBuffer(tweetJSON))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	service.Init(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "error processing tweet")
	mockStoreTweet.AssertExpectations(t)
}
