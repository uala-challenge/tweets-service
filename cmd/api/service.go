package main

import (
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_builder"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_engine"
	"github.com/uala-challenge/tweet-service/cmd/api/post_tweet"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/save_item"
	publish_tweet_event_sns2 "github.com/uala-challenge/tweet-service/internal/platform/sns/publish_tweet_event_sns"
	"github.com/uala-challenge/tweet-service/internal/store_tweet"
	"github.com/uala-challenge/tweet-service/kit/config"
)

type engine struct {
	simplify         app_engine.Engine
	repositories     repositories
	useCases         useCases
	handlers         handlers
	repositoryConfig config.RepositoryConfig
	useCasesConfig   config.UsesCasesConfig
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
	a.engine.repositoryConfig = app_engine.GetConfig[config.RepositoryConfig](a.engine.simplify.RepositoriesConfig)
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
		Config: a.engine.repositoryConfig.DataBase,
	})
	return a
}

func (a AppBuilder) InitUseCases() app_builder.Builder {
	a.engine.useCases.StoreTweet = store_tweet.NewService(store_tweet.Dependencies{
		SNS:    a.engine.repositories.PublishTweet,
		DB:     a.engine.repositories.SaveTweet,
		Log:    a.engine.simplify.Log,
		Config: a.engine.useCasesConfig.Store,
	})
	return a
}

func (a AppBuilder) InitHandlers() app_builder.Builder {
	a.engine.handlers.Test = post_tweet.NewService(post_tweet.Dependencies{UseCaseStoreTweet: a.engine.useCases.StoreTweet})
	return a
}

func (a AppBuilder) InitRoutes() app_builder.Builder {
	a.engine.simplify.App.Router.Post("/tweet", a.engine.handlers.Test.Init)
	return a
}

func (a AppBuilder) Build() app_builder.App {
	return a.engine
}
