package service

import (
	"github.com/fngomez/go-twitter-concurrency/src/domain"
	"os"
)

type TweetWriter interface {
	WriteTweet(domain.Tweet)
}

type MemoryTweetWriter struct {
	Tweets []domain.Tweet
}

type ChannelTweetWriter struct {
	tweetWriter TweetWriter
}

func NewMemoryTweetWriter() *MemoryTweetWriter {
	return &MemoryTweetWriter{}
}

func NewChannelTweetWriter(tweetWriter TweetWriter) ChannelTweetWriter {
	return ChannelTweetWriter{tweetWriter }
}

func (memoryTweetWriter *MemoryTweetWriter) WriteTweet(tweet domain.Tweet){
	memoryTweetWriter.Tweets = append(memoryTweetWriter.Tweets, tweet)
}

func (channel *ChannelTweetWriter) WriteTweet(tweetsToWrite chan domain.Tweet, quit chan bool){
	tweet, open := <-tweetsToWrite

	for open {
		channel.tweetWriter.WriteTweet(tweet)
		tweet, open = <-tweetsToWrite
	}

	quit <- true
}

type FileTweetWriter struct {
	file *os.File
}

func NewFileTweetWriter() *FileTweetWriter {

	file, _ := os.OpenFile(
		"tweets.txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)

	writer := new(FileTweetWriter)
	writer.file = file

	return writer
}

func (writer *FileTweetWriter) WriteTweet(tweet domain.Tweet) {

	if writer.file != nil {
		byteSlice := []byte(tweet.PrintableTweet() + "\n")
		writer.file.Write(byteSlice)
	}
}