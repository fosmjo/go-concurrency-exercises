//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer szenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream) (tweets []*Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return tweets
		}

		tweets = append(tweets, tweet)
	}
}

func producer2(done <-chan bool, stream Stream) <-chan *Tweet {
	tweets := make(chan *Tweet)

	go func() {
		// defer fmt.Println("\nend \"producer2\"") // check goroutine leak
		defer close(tweets)

		for {
			select {
			case <-done:
				return
			default:
				tweet, err := stream.Next()
				if err == ErrEOF {
					return
				}
				tweets <- tweet
			}
		}
	}()

	return tweets
}

func consumer(tweets []*Tweet) {
	for _, t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func consumer2(tweets <-chan *Tweet) {
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func consumer3(tweets <-chan *Tweet) {
	var wg sync.WaitGroup

	for t := range tweets {
		wg.Add(1)

		go func(t *Tweet) {
			defer wg.Done()

			if t.IsTalkingAboutGo() {
				fmt.Println(t.Username, "\ttweets about golang")
			} else {
				fmt.Println(t.Username, "\tdoes not tweet about golang")
			}
		}(t)
	}

	wg.Wait()
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	done := make(chan bool)
	defer close(done)

	// Producer
	tweets := producer2(done, stream)

	// Consumer
	consumer3(tweets)

	fmt.Printf("Process took %s\n", time.Since(start))
}
