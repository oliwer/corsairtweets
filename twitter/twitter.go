package twitter

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/StalkR/goircbot/bot"
)

var (
	api *anaconda.TwitterApi
)

// Verify if the Twitter credentials are correct
func creds(e *bot.Event) {
	ok, err := api.VerifyCredentials()
	if ok {
		e.Bot.Privmsg(e.Target, "Twitter credentials verified.")
	} else {
		e.Bot.Privmsg(e.Target, fmt.Sprintf(
			"Failed to verify Twitter credentials: %v", err))
	}
}

// Post a new tweet
func tweet(e *bot.Event) {
	_, err := api.PostTweet(e.Args, nil)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("Failed to send tweet: %v", err))
	} else {
		e.Bot.Privmsg(e.Target, "Tweet sent.")
	}
}

// Register the twitter plugin with a bot
func Register(b bot.Bot, appkey, appsecret, token, toksecret string) {
	anaconda.SetConsumerKey(appkey)
	anaconda.SetConsumerSecret(appsecret)

	api = anaconda.NewTwitterApi(token, toksecret)

	b.Commands().Add("tweet", bot.Command{
		Help:    "Send a message on the #CorsairIRC Tweeter feed",
		Handler: tweet,
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	b.Commands().Add("twcreds", bot.Command{
		Help:    "Test Twitter credentials",
		Handler: creds,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
