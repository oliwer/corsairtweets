// Twitter IRC Bot for #Corsair
// Written by Olivier Duclos (odyssey)

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/plugins/imdb"
	"github.com/StalkR/goircbot/plugins/sed"
	"github.com/StalkR/goircbot/plugins/up"
	"github.com/StalkR/goircbot/plugins/urban"
	"github.com/StalkR/goircbot/plugins/weather"
	"github.com/oliwer/corsairtweets/lastseen"
	"github.com/oliwer/corsairtweets/timein"
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
	ignored   = strings.Split(os.Getenv("IGNORE_NICKS"), ",")
	owmkey    = os.Getenv("OPENWEATHERMAP_KEY")
)

func main() {
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Done() // it not allowed to call wg.Wait() before it has been used.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-c
		log.Printf("Received signal %s. Exiting...\n", sig.String())
		wg.Wait()
		os.Exit(1)
	}()

	b := bot.NewBot(*host, *ssl, *nick, *ident, []string{*channel})
	imdb.Register(b)
	lastseen.Register(b, ignored, &wg)
	sed.Register(b)
	timein.Register(b)
	twitter.Register(b, appkey, appsecret)
	up.Register(b)
	urban.Register(b)
	weather.Register(b, owmkey)
	b.Run()
}
