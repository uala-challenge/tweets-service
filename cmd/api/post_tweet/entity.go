package post_tweet

import (
	"github.com/uala-challenge/tweet-service/internal/store_tweet"
	"net/http"
)

type Service interface {
	Init(w http.ResponseWriter, r *http.Request)
}

type Dependencies struct {
	UseCaseStoreTweet store_tweet.Service
}
