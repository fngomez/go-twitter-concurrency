package service

import (
	"fmt"
	"github.com/fngomez/go-twitter-concurrency/src/domain"
)

type TweetManager struct {
	tweets       []domain.Tweet
	tweetsByUser map[string][]domain.Tweet
	channel ChannelTweetWriter
	tweetChannel chan domain.Tweet

}

func NewTweetManager(channel ChannelTweetWriter) *TweetManager {

	tweetManager := new(TweetManager)

	tweetManager.tweets = make([]domain.Tweet, 0)
	tweetManager.tweetsByUser = make(map[string][]domain.Tweet)

	tweetManager.channel = channel
	tweetManager.tweetChannel = make(chan domain.Tweet)


	return tweetManager
}

func (manager *TweetManager) PublishTweet(tweetToPublish domain.Tweet, quit chan bool) (int, error) {

	if tweetToPublish.GetUser() == "" {
		return 0, fmt.Errorf("user is required")
	}

	if tweetToPublish.GetText() == "" {
		return 0, fmt.Errorf("text is required")
	}

	if len(tweetToPublish.GetText()) > 140 {
		return 0, fmt.Errorf("text exceeds 140 characters")
	}

	manager.tweets = append(manager.tweets, tweetToPublish)

	tweetToPublish.SetId(len(manager.tweets))

	userTweets := manager.tweetsByUser[tweetToPublish.GetUser()]
	manager.tweetsByUser[tweetToPublish.GetUser()] = append(userTweets, tweetToPublish)

	manager.tweetChannel <- tweetToPublish

	manager.channel.WriteTweet(manager.tweetChannel, quit)

	return tweetToPublish.GetId(), nil
}

// GetTweet returns the last published tweet
func (manager *TweetManager) GetTweet() domain.Tweet {
	lastTweetIndex := len(manager.tweets) - 1

	return manager.tweets[lastTweetIndex]
}

func (manager *TweetManager) GetTweets() []domain.Tweet {
	return manager.tweets
}

func (manager *TweetManager) GetTweetById(id int) domain.Tweet {

	var tweet domain.Tweet

	tweetIndex := 0

	for tweetIndex < len(manager.tweets) && tweet == nil {

		actualTweet := manager.tweets[tweetIndex]

		if actualTweet.GetId() == id {
			tweet = actualTweet
		}

		tweetIndex++
	}

	return tweet
}

func (manager *TweetManager) CountTweetsByUser(user string) int {

	var count int

	for _, tweet := range manager.tweets {
		if tweet.GetUser() == user {
			count++
		}
	}

	return count
}

func (manager *TweetManager) GetTweetsByUser(user string) []domain.Tweet {

	return manager.tweetsByUser[user]
}
