package twitter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/StalkR/goircbot/bot"
	"github.com/garyburd/go-oauth/oauth"
)

const (
	tokenFilename = "access-token.dat"
	sep           = " "
)

var (
	api      *anaconda.TwitterApi
	tmpCreds *oauth.Credentials
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

// Request a login URL
func login(e *bot.Event) {
	url, creds, err := api.AuthorizationURL("oob")
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf(
			"failed to get Twitter authorization URL: %v", err))
		return
	}

	e.Bot.Privmsg(e.Target, "Please open this URL: "+url)
	e.Bot.Privmsg(e.Target, "After allowing the application to access your "+
		"account, Twitter will give you a PIN code. You must pass this "+
		"PIN to the bot with the following command:")
	e.Bot.Privmsg(e.Target, "> twpin [pincode]")

	tmpCreds = creds
}

// Validate an authentication PIN code from Twitter
func pin(e *bot.Event) {
	pin := strings.Trim(e.Args, " ")

	cred, _, err := api.GetCredentials(tmpCreds, pin)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("invalid pin code: %s", err))
		return
	}

	e.Bot.Privmsg(e.Target, "Pin code is valid. Login complete!")

	api = anaconda.NewTwitterApi(cred.Token, cred.Secret)
	tmpCreds = nil
	saveAccessToken(cred.Token, cred.Secret)
}

// Load the token from filesystem
func loadAccessToken() (token, toksecret string) {
	tokens, err := ioutil.ReadFile(tokenFilename)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Println("failed to load Twitter Access Token:", err)
		}
		return "", ""
	}

	tokens = bytes.TrimSpace(tokens)

	list := strings.Split(string(tokens), sep)
	return list[0], list[1]
}

// Save the token to a file
func saveAccessToken(token, toksecret string) {
	content := strings.Join([]string{token, toksecret}, sep)

	err := ioutil.WriteFile(tokenFilename, []byte(content), 0600)
	if err != nil {
		log.Println("failed to save access token:", err)
	}
}

// Register the twitter plugin with a bot
func Register(b bot.Bot, appkey, appsecret string) {
	anaconda.SetConsumerKey(appkey)
	anaconda.SetConsumerSecret(appsecret)

	token, toksecret := loadAccessToken()
	api = anaconda.NewTwitterApi(token, toksecret)

	b.Commands().Add("tweet", bot.Command{
		Help:    "Post a tweet on @CorsairIRC (for real)",
		Handler: tweet,
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	b.Commands().Add("twcreds", bot.Command{
		Help:    "Test Twitter credentials",
		Handler: creds,
		Pub:     false,
		Priv:    true,
		Hidden:  true})

	b.Commands().Add("twlogin", bot.Command{
		Help:    "Login to Twitter",
		Handler: login,
		Pub:     false,
		Priv:    true,
		Hidden:  true})

	b.Commands().Add("twpin", bot.Command{
		Help:    "Enter the PIN code to finalize login",
		Handler: pin,
		Pub:     false,
		Priv:    true,
		Hidden:  true})
}
