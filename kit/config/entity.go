package config

import (
	"github.com/uala-challenge/tweet-service/internal/platform/db/uala_challenge"
	"github.com/uala-challenge/tweet-service/internal/store_tweet"
)

type RepositoryConfig struct {
	DataBase uala_challenge.Config `json:"database"`
}

type UsesCasesConfig struct {
	Store store_tweet.Config `json:"store"`
}
