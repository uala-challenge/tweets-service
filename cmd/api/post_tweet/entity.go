package post_tweet

import (
	"net/http"

	"github.com/uala-challenge/tweet-service/internal/store_tweet"
)

type Service interface {
	Init(w http.ResponseWriter, r *http.Request)
}

type Dependencies struct {
	UseCaseStoreTweet store_tweet.Service
}
