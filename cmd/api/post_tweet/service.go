package post_tweet

import (
	"fmt"
	"io"
	"net/http"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/error_handler"
	"github.com/uala-challenge/tweets-service/internal/store_tweet"
	"github.com/uala-challenge/tweets-service/kit"
)

type service struct {
	useCase store_tweet.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		useCase: d.UseCaseStoreTweet,
	}
}

func (s service) Init(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = error_handler.HandleApiErrorResponse(error_handler.NewCommonApiError("bad request", err.Error(), err, http.StatusBadRequest), w)
		return
	}
	_ = r.Body.Close()

	rqt, _ := kit.BytesToModel[kit.TweetRequest](body)

	if err := rqt.Validate(); err != nil {
		_ = error_handler.HandleApiErrorResponse(error_handler.NewCommonApiError("validate error", err.Error(), err, http.StatusBadRequest), w)
		return
	}

	tweetId, err := s.useCase.Apply(r.Context(), rqt)
	if err != nil {
		_ = error_handler.HandleApiErrorResponse(error_handler.NewCommonApiError("error processing tweet", err.Error(), err, http.StatusInternalServerError), w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(fmt.Sprintf("tweet creado: %s", tweetId)))

}
