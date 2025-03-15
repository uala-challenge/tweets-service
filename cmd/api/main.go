package main

import "github.com/uala-challenge/simple-toolkit/pkg/simplify/app_builder"

// @title                     			tweets-service
// @description              			Servicio que maneja la publicación y consulta de tweets.
// @contact.name 						Fredy Hernandez
// @contact.email						hfredy717@outlook.com
// swag only supports one host. If the one selected for you is not correct, you should manually change it.
func main() {
	builder := NewAppBuilder()
	application := app_builder.Apply(builder)
	err := application.Run()
	if err != nil {
		panic(err)
	}
}
