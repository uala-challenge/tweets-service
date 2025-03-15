package list_tweets

import (
	"net/http"

	"github.com/uala-challenge/tweet-service/internal/batch_get_tweets"
)

type Service interface {
	Init(w http.ResponseWriter, r *http.Request)
}

type Dependencies struct {
	UseCaseListTweets batch_get_tweets.Service
}
