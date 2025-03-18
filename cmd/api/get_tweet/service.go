package get_tweet

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/error_handler"

	"github.com/uala-challenge/tweets-service/internal/retrieve_tweet"
	_ "github.com/uala-challenge/tweets-service/kit"
)

// service representa el servicio que maneja la obtención de tweets.
type service struct {
	useCase retrieve_tweet.Service
}

var _ Service = (*service)(nil)

// NewService crea una nueva instancia del servicio de obtención de tweets.
func NewService(d Dependencies) Service {
	return &service{
		useCase: d.UseCaseRetrieveTweet,
	}
}

// Init maneja la solicitud GET para obtener un tweet específico.
//
// @Summary Obtener un tweet específico
// @Description Retorna un tweet basado en el `user_id` y `tweet_id` proporcionados.
// @Tags tweets
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID del usuario propietario del tweet"
// @Param tweet_id path string true "ID del tweet a consultar"
// @Success 200 {object} kit.Tweet "Tweet obtenido correctamente"
// @Failure 400 {object} kit.CommonApiError "Solicitud incorrecta"
// @Failure 500 {object} kit.CommonApiError "Error interno del servidor"
// @Router /tweets/{user_id}/{tweet_id} [get]
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
