package main

import (
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_builder"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/app_engine"
	"github.com/uala-challenge/tweet-service/cmd/api/get_tweet"
	"github.com/uala-challenge/tweet-service/cmd/api/list_tweets"
	"github.com/uala-challenge/tweet-service/cmd/api/post_tweet"
	"github.com/uala-challenge/tweet-service/internal/batch_get_tweets"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/get_item"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/get_items"
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge/save_item"
	publish_tweet_event_sns2 "github.com/uala-challenge/tweet-service/internal/platform/sns/publish_tweet_event_sns"
	"github.com/uala-challenge/tweet-service/internal/retrieve_tweet"
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
	a.engine.repositories.GetTweets = get_items.NewService(get_items.Dependencies{
		Client: a.engine.simplify.DynamoDBClient,
		Config: a.engine.repositoryConfig.DataBase,
		Log:    a.engine.simplify.Log,
	})
	a.engine.repositories.GetTweet = get_item.NewService(get_item.Dependencies{
		Client: a.engine.simplify.DynamoDBClient,
		Config: a.engine.repositoryConfig.DataBase,
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
	a.engine.useCases.BatchGetTweets = batch_get_tweets.NewService(batch_get_tweets.Dependencies{
		DBRepository: a.engine.repositories.GetTweets,
		Log:          a.engine.simplify.Log,
	})
	return a
}

func (a AppBuilder) InitHandlers() app_builder.Builder {
	a.engine.handlers.PostTweet = post_tweet.NewService(post_tweet.Dependencies{UseCaseStoreTweet: a.engine.useCases.StoreTweet})
	a.engine.handlers.GetTweet = get_tweet.NewService(get_tweet.Dependencies{UseCaseRetrieveTweet: a.engine.useCases.RetrieveTweet})
	a.engine.handlers.ListTweet = list_tweets.NewService(list_tweets.Dependencies{UseCaseListTweets: a.engine.useCases.BatchGetTweets})
	return a
}

func (a AppBuilder) InitRoutes() app_builder.Builder {
	a.engine.simplify.App.Router.Post("/tweets", a.engine.handlers.PostTweet.Init)
	a.engine.simplify.App.Router.Get("/tweets/{tweet_id}/users/{user_id}", a.engine.handlers.GetTweet.Init)
	a.engine.simplify.App.Router.Get("/tweets", a.engine.handlers.ListTweet.Init)
	return a
}

func (a AppBuilder) Build() app_builder.App {
	return a.engine
}
