package mappers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/uala-challenge/tweets-service/kit"
)

var prefixes = map[string]string{
	"PK":     "tweet:",
	"GSI1SK": "tweet:",
}

func TweetDynamoMap(rqt kit.TweetRequest, tweetID string) kit.DynamoItem {
	return kit.DynamoItem{
		PK:      prefixes["PK"] + tweetID,
		SK:      rqt.UserID,
		GSI1PK:  rqt.UserID,
		GSI1SK:  prefixes["GSI1SK"] + tweetID,
		Content: rqt.Tweet,
		Created: time.Now().Unix(),
	}
}

func TweetSNSMap(rqt kit.DynamoItem) kit.Tweet {
	return kit.Tweet{
		UserID:  rqt.SK,
		TweetID: rqt.PK,
		Created: rqt.Created,
	}
}

func TweetToPublishInput(item kit.Tweet, topic string) *sns.PublishInput {
	jsonBytes, _ := json.Marshal(item)
	messageStr := string(jsonBytes)
	return &sns.PublishInput{
		TopicArn: aws.String(topic),
		Message:  aws.String(messageStr),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"user_id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(item.UserID),
			},
			"tweet_id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(item.TweetID),
			},
			"created": {
				DataType:    aws.String("Number"),
				StringValue: aws.String(fmt.Sprintf("%d", item.Created)),
			},
		},
	}
}
