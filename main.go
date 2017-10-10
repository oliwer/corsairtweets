// Twitter IRC Bot for #Corsair
// Written by Olivier Duclos (odyssey)

package main

import (
	"flag"
	"os"
	"github.com/StalkR/goircbot/bot"
	"github.com/oliwer/corsairtweets/twitter"
)

var (
	host      = flag.String("host", "irc.oftc.net", "Server host[:port]")
	ssl       = flag.Bool("ssl", true, "Connect with SSL")
	nick      = flag.String("nick", "twittard", "Bot nick")
	ident     = flag.String("ident", "corsairtwitterebot", "Bot ident")
	channel   = flag.String("channel", "#Corsair", "Channel to join")
	appkey    = os.Getenv("TWITTER_APP_KEY")
	appsecret = os.Getenv("TWITTER_APP_SECRET")
)

func main() {
	flag.Parse()

	b := bot.NewBot(*host, *ssl, *nick, *ident, []string{*channel})
	twitter.Register(b, appkey, appsecret)
	b.Run()
}
