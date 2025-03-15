package post_tweet

import (
	"fmt"
	"io"
	"net/http"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/error_handler"
	"github.com/uala-challenge/tweet-service/internal/store_tweet"
	"github.com/uala-challenge/tweet-service/kit"
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

// Init godoc
// @Summary Send a message to the integration hub
// @Description Send a message to the integration hub
// @Tags Send, Notification
// @Produce json
// @Param string
// @Success 200 {object} interface{}
// @Failure 400  {object}  string
// @Failure 424  {object}  error_wrapper.CommonApiError "The type or communication channel requested to send the notification is not enabled"
// @Failure 403  {object}  error_wrapper.CommonApiError "Preference is not allowed within the capability: The integration flow is not configured in the capability for this preference."
// @Failure 428  {object}  error_wrapper.CommonApiError "The preference is not enabled: the user has not enabled the delivery of this type of notification."
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
	_, _ = w.Write([]byte(fmt.Sprintf("DynamoItem creado: %s", tweetId)))

}
