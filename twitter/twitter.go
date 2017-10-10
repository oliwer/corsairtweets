package twitter

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	// rm -rf anaconda/vendor/github.com/garyburd/oauth to fix compilation
	"github.com/ChimeraCoder/anaconda"
	"github.com/StalkR/goircbot/bot"
	"github.com/garyburd/go-oauth/oauth"
)

const (
	TokenFilename = "access-token.dat"
	Sep = " "
)

var (
	api *anaconda.TwitterApi
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
	url, creds, err := anaconda.AuthorizationURL("oob")
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf(
			"failed to get Twitter authorization URL: %v", err))
		return
	}

	e.Bot.Privmsg(e.Target, "Please open this URL: " + url)
	e.Bot.Privmsg(e.Target, "After allowing the application to access your " +
		"account, Twitter will give you a PIN code. You must pass this " +
		"PIN to the bot with the following command:")
	e.Bot.Privmsg(e.Target, "> twpin [pincode]")

	tmpCreds = creds
}

// Validate an authentication PIN code from Twitter
func pin(e *bot.Event) {
	pin := strings.Trim(e.Args, " ")

	cred, _, err := anaconda.GetCredentials(tmpCreds, pin)
	if err != nil {
		e.Bot.Privmsg(e.Target, fmt.Sprintf("invalid pin code: %s", err))
		return
	}

	api = anaconda.NewTwitterApi(cred.Token, cred.Secret)
	tmpCreds = nil
	saveAccessToken(cred.Token, cred.Secret)
}

// Load the token from filesystem
func loadAccessToken() (token, toksecret string) {
	bytes, err := ioutil.ReadFile(TokenFilename)
	if err != nil {
		if ! os.IsNotExist(err) {
			log.Println("failed to load Twitter Access Token:", err)
		}
		return "", ""
	}

	list := strings.Split(string(bytes), Sep)
	return list[0], list[1]
}

// Save the token to a file
func saveAccessToken(token, toksecret string) {
	content := strings.Join([]string{token,toksecret}, Sep)

	err := ioutil.WriteFile(TokenFilename, []byte(content), 0600)
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

	b.Commands().Add("twlogin", bot.Command{
		Help:    "Login to Twitter",
		Handler: login,
		Pub:     true,
		Priv:    true,
		Hidden:  false})

	b.Commands().Add("twpin", bot.Command{
		Help:    "Enter the PIN code to finalize login",
		Handler: pin,
		Pub:     true,
		Priv:    true,
		Hidden:  false})
}
