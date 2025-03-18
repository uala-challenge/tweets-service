package main

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/get_item"
	"github.com/uala-challenge/simple-toolkit/pkg/platform/db/save_item"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_builder"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_engine"
	"github.com/uala-challenge/tweets-service/cmd/api/get_tweet"
	"github.com/uala-challenge/tweets-service/cmd/api/post_tweet"
	_ "github.com/uala-challenge/tweets-service/docs"
	publish_tweet_event_sns2 "github.com/uala-challenge/tweets-service/internal/platform/publish_tweet_event_sns"
	"github.com/uala-challenge/tweets-service/internal/retrieve_tweet"
	"github.com/uala-challenge/tweets-service/internal/store_tweet"
	"github.com/uala-challenge/tweets-service/kit/config"
)

type engine struct {
	simplify       app_engine.Engine
	repositories   repositories
	useCases       useCases
	handlers       handlers
	useCasesConfig config.UsesCasesConfig
}

type AppBuilder struct {
	engine *engine
}

var _ app_builder.Builder = (*AppBuilder)(nil)

func NewAppBuilder() *AppBuilder {
	a := *app_engine.NewApp()
	return &AppBuilder{
		engine: &engine{
			simplify: a,
		},
	}
}

func (a engine) Run() error {
	return a.simplify.App.Run()
}

func (a AppBuilder) LoadConfig() app_builder.Builder {
	a.engine.useCasesConfig = app_engine.GetConfig[config.UsesCasesConfig](a.engine.simplify.UsesCasesConfig)
	return a
}

func (a AppBuilder) InitRepositories() app_builder.Builder {
	a.engine.repositories.PublishTweet = publish_tweet_event_sns2.NewService(publish_tweet_event_sns2.Dependencies{
		Client: a.engine.simplify.SNSClient,
		Log:    a.engine.simplify.Log,
	})
	a.engine.repositories.SaveTweet = save_item.NewService(save_item.Dependencies{
		Client: a.engine.simplify.DynamoDBClient,
		Log:    a.engine.simplify.Log,
	})
	a.engine.repositories.GetTweet = get_item.NewService(get_item.Dependencies{
		Client: a.engine.simplify.DynamoDBClient,
		Log:    a.engine.simplify.Log,
	})
	return a
}

func (a AppBuilder) InitUseCases() app_builder.Builder {
	a.engine.useCases.StoreTweet = store_tweet.NewService(store_tweet.Dependencies{
		SNSRepository: a.engine.repositories.PublishTweet,
		DBRepository:  a.engine.repositories.SaveTweet,
		Log:           a.engine.simplify.Log,
		Config:        a.engine.useCasesConfig.Store,
	})
	a.engine.useCases.RetrieveTweet = retrieve_tweet.NewService(retrieve_tweet.Dependencies{
		DBRepository: a.engine.repositories.GetTweet,
		Log:          a.engine.simplify.Log,
	})
	return a
}

func (a AppBuilder) InitHandlers() app_builder.Builder {
	a.engine.handlers.PostTweet = post_tweet.NewService(post_tweet.Dependencies{UseCaseStoreTweet: a.engine.useCases.StoreTweet})
	a.engine.handlers.GetTweet = get_tweet.NewService(get_tweet.Dependencies{UseCaseRetrieveTweet: a.engine.useCases.RetrieveTweet})
	return a
}

func (a AppBuilder) InitRoutes() app_builder.Builder {
	a.engine.simplify.App.Router.Post("/tweet", a.engine.handlers.PostTweet.Init)
	a.engine.simplify.App.Router.Get("/tweet/{tweet_id}/user/{user_id}", a.engine.handlers.GetTweet.Init)
	a.engine.simplify.App.Router.Get("/swagger/*", httpSwagger.WrapHandler)
	return a
}

func (a AppBuilder) Build() app_builder.App {
	return a.engine
}
