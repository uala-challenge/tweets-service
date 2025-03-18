package main

import "github.com/uala-challenge/simple-toolkit/pkg/simplify/app_builder"

// @title tweets-service
// @version 1.0
// @description Servicio que maneja la publicación y consulta de tweets.
// @host localhost:8084
// @BasePath /
// @schemes http
//
// @contact.name Fredy Hernandez
// @contact.email hfredy717@outlook.com
//
// @swagger:meta
//
// tweets-service es una API diseñada para permitir la publicación y consulta de tweets.
//
// ## Funcionalidades principales:
// - **Publicar tweets**: Permite a los usuarios enviar tweets, los cuales se almacenan y son procesados asincrónicamente.
// - **Consultar tweets**: Proporciona acceso a tweets previamente publicados, basados en el ID del usuario y el ID del tweet.
//
// ## Endpoints disponibles:
// - `POST /tweets`: Publica un nuevo tweet y lo envía a la cola de mensajes para su procesamiento asincrónico.
// - `GET /tweets/{user_id}/{tweet_id}`: Obtiene un tweet específico por usuario y ID del tweet.
func main() {
	builder := NewAppBuilder()
	application := app_builder.Apply(builder)
	err := application.Run()
	if err != nil {
		panic(err)
	}
}
