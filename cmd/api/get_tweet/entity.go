package get_tweet

import (
	"net/http"

	"github.com/uala-challenge/tweet-service/internal/retrieve_tweet"
)

type Service interface {
	Init(w http.ResponseWriter, r *http.Request)
}

type Dependencies struct {
	UseCaseRetrieveTweet retrieve_tweet.Service
}
