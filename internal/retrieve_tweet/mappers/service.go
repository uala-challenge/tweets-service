package mappers

import "github.com/uala-challenge/tweet-service/kit"

func DynamoItemToTweet(item kit.DynamoItem) kit.Tweet {
	return kit.Tweet{
		UserID:  item.SK,
		TweetID: item.PK,
		Created: item.Created,
		Content: item.Content,
	}
}
