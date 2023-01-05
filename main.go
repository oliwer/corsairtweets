// Twitter IRC Bot for #Corsair
// Written by Olivier Duclos (odyssey)

package main

import (
	"flag"
	"os"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/plugins/idle"
	"github.com/StalkR/goircbot/plugins/imdb"
	"github.com/StalkR/goircbot/plugins/sed"
	"github.com/StalkR/goircbot/plugins/up"
	"github.com/StalkR/goircbot/plugins/urban"
	"github.com/StalkR/goircbot/plugins/weather"
	"github.com/oliwer/corsairtweets/lastseen"
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
	ignored   = strings.Split(os.Getenv("IDLE_IGNORE"), ",")
	owmkey    = os.Getenv("OPENWEATHERMAP_KEY")
)

func main() {
	flag.Parse()

	b := bot.NewBot(*host, *ssl, *nick, *ident, []string{*channel})
	idle.Register(b, ignored)
	imdb.Register(b)
	lastseen.Register(b, ignored)
	sed.Register(b)
	twitter.Register(b, appkey, appsecret)
	up.Register(b)
	urban.Register(b)
	weather.Register(b, owmkey)
	b.Run()
}
