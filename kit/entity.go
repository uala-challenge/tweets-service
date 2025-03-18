package kit


// Tweet representa la respuesta de un tweet consultado.
// @swagger:model Tweet
type Tweet struct {
	// ID del usuario propietario del tweet.
	// example: "user:12345"
	UserID string `json:"user_id"`

	// ID único del tweet.
	// example: "tweet:67890"
	TweetID string `json:"tweet_id"`

	// Fecha de creación del tweet en formato UNIX timestamp.
	// example: 1710597823
	Created int64 `json:"created"`

	// Contenido del tweet.
	// example: "Este es un tweet de prueba."
	Content string `json:"content"`
}

// CommonApiError representa la estructura de respuesta en caso de error.
// @swagger:model CommonApiError
// @example {"code": "ERR-400", "msg": "Solicitud incorrecta"}
type CommonApiError struct {
	// Código del error.
	// example: "ERR-400"
	Code string `json:"code"`

	// Mensaje descriptivo del error.
	// example: "Solicitud incorrecta"
	Msg string `json:"msg"`

	// Error interno (no expuesto en JSON).
	Err error `json:"-"`

	// Código HTTP asociado al error (no expuesto en JSON).
	HttpCode int `json:"-"`
}

// TweetRequest representa la estructura de solicitud para crear un tweet.
// @swagger:model TweetRequest
// @example {"user_id": "user:12345", "tweet": "Este es un tweet de ejemplo con una longitud adecuada."}
type TweetRequest struct {
	// ID del usuario que publica el tweet.
	// Required: true
	// example: "user:12345"
	UserID string `json:"user_id" validate:"required"`

	// Contenido del tweet.
	// Required: true
	// Max length: 280 caracteres
	// example: "Este es un tweet de prueba con una longitud adecuada."
	Tweet string `json:"tweet" validate:"required,max=280"`
}

type DynamoItem struct {
	PK      string `dynamodbav:"PK"`
	SK      string `dynamodbav:"SK"`
	GSI1PK  string `dynamodbav:"GSI1PK"`
	GSI1SK  string `dynamodbav:"GSI1SK"`
	Content string `dynamodbav:"content"`
	Created int64  `dynamodbav:"created"`
}

type TweetPK struct {
	UserID  string `json:"user_id"`
	TweetID string `json:"tweet_id"`
}
