package post_tweet

import (
	"fmt"
	"io"
	"net/http"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/error_handler"
	"github.com/uala-challenge/tweets-service/internal/store_tweet"
	"github.com/uala-challenge/tweets-service/kit"
)

// service representa el servicio que maneja la publicación de tweets.
type service struct {
	useCase store_tweet.Service
}

var _ Service = (*service)(nil)

// NewService crea una nueva instancia del servicio de publicación de tweets.
func NewService(d Dependencies) Service {
	return &service{
		useCase: d.UseCaseStoreTweet,
	}
}

// Init maneja la solicitud POST para publicar un tweet.
//
// @Summary Publicar un tweet
// @Description Permite a los usuarios publicar un nuevo tweet, almacenarlo en la base de datos y enviarlo a la cola de mensajes para su procesamiento asincrónico.
// @Tags tweets
// @Accept  json
// @Produce  json
// @Param body body kit.TweetRequest true "Datos del tweet a publicar"
// @Success 201 {string} string "Tweet creado correctamente"
// @Failure 400 {object} kit.CommonApiError "Solicitud incorrecta"
// @Failure 500 {object} kit.CommonApiError "Error interno del servidor"
// @Router /tweets [post]
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(fmt.Sprintf("tweet creado: %s", tweetId)))
}