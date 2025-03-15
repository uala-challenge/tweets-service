package list_tweets

import (
	"io"
	"net/http"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/error_handler"
	"github.com/uala-challenge/tweet-service/internal/batch_get_tweets"
	"github.com/uala-challenge/tweet-service/kit"
)

type service struct {
	useCase batch_get_tweets.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) Service {
	return &service{
		useCase: d.UseCaseListTweets,
	}
}

func (s service) Init(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		_ = error_handler.HandleApiErrorResponse(error_handler.NewCommonApiError("bad request", err.Error(), err, http.StatusBadRequest), w)
		return
	}
	_ = r.Body.Close()

	rqt, _ := kit.BytesToSlice[kit.TweetPK](body)

	tweets, err := s.useCase.Apply(r.Context(), rqt)
	if err != nil {
		_ = error_handler.HandleApiErrorResponse(error_handler.NewCommonApiError("error processing tweet", err.Error(), err, http.StatusInternalServerError), w)
		return
	}

	rps, _ := kit.SliceToBytes[kit.Tweet](tweets)

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(rps)

}
