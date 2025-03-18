package get_tweet

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/error_handler"

	"github.com/uala-challenge/tweets-service/internal/retrieve_tweet"
)

type service struct {
	useCase retrieve_tweet.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		useCase: d.UseCaseRetrieveTweet,
	}
}

func (s service) Init(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	tweetID := chi.URLParam(r, "tweet_id")

	tweet, err := s.useCase.Apply(r.Context(), map[string]interface{}{"SK": userID, "PK": tweetID})
	if err != nil {
		_ = error_handler.HandleApiErrorResponse(error_handler.NewCommonApiError("error to get tweet", err.Error(), err, http.StatusInternalServerError), w)
		return
	}

	rsp, err := json.Marshal(tweet)
	if err != nil {
		_ = error_handler.HandleApiErrorResponse(error_handler.NewCommonApiError("error to marshal tweet", err.Error(), err, http.StatusInternalServerError), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rsp)
}
